package document

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"keden-service/back/internal/models"
	docRepo "keden-service/back/internal/repositories/database/postgres/document"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	localRabbit "keden-service/back/internal/pkg/rabbitmq"
)

var (
	ErrDocumentNotFound = errors.New("document not found")
	ErrNotPDF           = errors.New("only PDF files are allowed")
	ErrFileTooLarge     = errors.New("file size exceeds 50MB limit")
	ErrExcelNotReady    = errors.New("excel file is not ready yet")
	ErrAccessDenied     = errors.New("access denied")
)

const MaxFileSize = 50 * 1024 * 1024 // 50MB

type DocumentService struct {
	docRepo  docRepo.IDocumentRepository
	rabbitMq *localRabbit.AmqpPubSub
}

func NewDocumentService(
	dr docRepo.IDocumentRepository,
	rmq *localRabbit.AmqpPubSub,
) *DocumentService {
	return &DocumentService{
		docRepo:  dr,
		rabbitMq: rmq,
	}
}

func (s *DocumentService) UploadDocument(ctx context.Context, userID uint, file multipart.File, header *multipart.FileHeader) (*models.Document, error) {
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".pdf" {
		return nil, ErrNotPDF
	}

	if header.Size > MaxFileSize {
		return nil, ErrFileTooLarge
	}

	doc := &models.Document{
		UserID:       userID,
		OriginalName: header.Filename,
		Status:       models.DocumentStatusUploaded,
		FileSize:     header.Size,
	}

	if err := s.docRepo.Create(ctx, doc); err != nil {
		return nil, err
	}

	go func() {
		if err := s.queueForProcessing(doc.ID); err != nil {
			logrus.Errorf("Failed to queue document %d: %v", doc.ID, err)
		}
	}()

	return doc, nil
}

func (s *DocumentService) queueForProcessing(docID uint) error {
	ctx := context.Background()

	if err := s.docRepo.UpdateStatus(ctx, docID, models.DocumentStatusQueued, ""); err != nil {
		return err
	}

	payload, _ := json.Marshal(map[string]interface{}{
		"document_id": docID,
	})

	msg := message.NewMessage(uuid.New().String(), payload)
	return s.rabbitMq.Publish("document.processing", msg)
}

func (s *DocumentService) GetDocuments(ctx context.Context, userID uint) ([]models.Document, error) {
	return s.docRepo.GetByUserID(ctx, userID)
}

func (s *DocumentService) GetDocumentByID(ctx context.Context, docID uint, userID uint) (*models.Document, error) {
	doc, err := s.docRepo.GetByID(ctx, docID)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, ErrDocumentNotFound
	}
	if doc.UserID != userID {
		return nil, ErrAccessDenied
	}
	return doc, nil
}

func (s *DocumentService) DownloadExcel(ctx context.Context, docID uint, userID uint) ([]byte, string, error) {
	doc, err := s.docRepo.GetByID(ctx, docID)
	if err != nil {
		return nil, "", err
	}
	if doc == nil {
		return nil, "", ErrDocumentNotFound
	}
	if doc.UserID != userID {
		return nil, "", ErrAccessDenied
	}
	if doc.Status != models.DocumentStatusCompleted || doc.ExcelFilePath == "" {
		return nil, "", ErrExcelNotReady
	}

	data, err := os.ReadFile(doc.ExcelFilePath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read excel file: %w", err)
	}

	excelName := strings.TrimSuffix(doc.OriginalName, filepath.Ext(doc.OriginalName)) + ".xlsx"
	return data, excelName, nil
}

func (s *DocumentService) GetAllDocuments(ctx context.Context) ([]models.Document, error) {
	return s.docRepo.GetAll(ctx)
}

func (s *DocumentService) GetDocumentForProcessing(ctx context.Context, docID uint) (*models.Document, error) {
	doc, err := s.docRepo.GetByID(ctx, docID)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, ErrDocumentNotFound
	}
	return doc, nil
}

func (s *DocumentService) UpdateDocumentAfterAI(ctx context.Context, docID uint, aiJSON string, excelData []byte) error {
	doc, err := s.docRepo.GetByID(ctx, docID)
	if err != nil {
		return err
	}
	if doc == nil {
		return ErrDocumentNotFound
	}

	dir := fmt.Sprintf("./uploads/excel/%d", doc.UserID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	filePath := filepath.Join(dir, uuid.New().String()+".xlsx")
	if err := os.WriteFile(filePath, excelData, 0644); err != nil {
		return fmt.Errorf("failed to write excel file: %w", err)
	}

	now := time.Now()
	doc.AIResponseJSON = aiJSON
	doc.ExcelFilePath = filePath
	doc.Status = models.DocumentStatusCompleted
	doc.ProcessedAt = &now

	return s.docRepo.Update(ctx, doc)
}

func (s *DocumentService) MarkDocumentError(ctx context.Context, docID uint, errMsg string) error {
	return s.docRepo.UpdateStatus(ctx, docID, models.DocumentStatusError, errMsg)
}

func (s *DocumentService) MarkDocumentProcessing(ctx context.Context, docID uint) error {
	return s.docRepo.UpdateStatus(ctx, docID, models.DocumentStatusProcessing, "")
}

func (s *DocumentService) GetStats(ctx context.Context) (total int64, completed int64, err error) {
	total, err = s.docRepo.GetTotalCount(ctx)
	if err != nil {
		return
	}
	completed, err = s.docRepo.GetCountByStatus(ctx, models.DocumentStatusCompleted)
	return
}
