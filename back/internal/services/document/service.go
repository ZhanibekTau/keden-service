package document

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"keden-service/back/internal/models"
	docRepo "keden-service/back/internal/repositories/database/postgres/document"
	fieldsRepo "keden-service/back/internal/repositories/database/postgres/document_fields"
	itemsRepo "keden-service/back/internal/repositories/database/postgres/document_items"
	"keden-service/back/internal/services/ai"
	"keden-service/back/internal/services/excel"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	ErrDocumentNotFound = errors.New("document not found")
	ErrNotPDF           = errors.New("only PDF files are allowed")
	ErrFileTooLarge     = errors.New("file size exceeds 50MB limit")
	ErrExcelNotReady    = errors.New("excel file is not ready yet")
	ErrAccessDenied     = errors.New("access denied")
	ErrAIDataNotReady   = errors.New("AI data is not ready yet")
)

const MaxFileSize = 50 * 1024 * 1024 // 50MB

type DocumentService struct {
	docRepo      docRepo.IDocumentRepository
	fieldsRepo   fieldsRepo.IDocumentFieldsRepository
	itemsRepo    itemsRepo.IDocumentItemRepository
	aiService    *ai.AIService
	excelService *excel.ExcelService
}

func NewDocumentService(
	dr docRepo.IDocumentRepository,
	fr fieldsRepo.IDocumentFieldsRepository,
	ir itemsRepo.IDocumentItemRepository,
	as *ai.AIService,
	es *excel.ExcelService,
) *DocumentService {
	return &DocumentService{
		docRepo:      dr,
		fieldsRepo:   fr,
		itemsRepo:    ir,
		aiService:    as,
		excelService: es,
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

	pdfData, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	aiResp, err := s.aiService.ProcessDocument(ctx, pdfData)
	if err != nil {
		return nil, fmt.Errorf("AI processing failed: %w", err)
	}

	now := time.Now()
	doc := &models.Document{
		UserID:       userID,
		OriginalName: header.Filename,
		Status:       models.DocumentStatusCompleted,
		FileSize:     header.Size,
		ProcessedAt:  &now,
	}
	if err := s.docRepo.Create(ctx, doc); err != nil {
		return nil, err
	}

	docFields := aiRespToFields(doc.ID, aiResp)
	if err := s.fieldsRepo.Create(ctx, docFields); err != nil {
		return nil, fmt.Errorf("failed to save document fields: %w", err)
	}

	docItems := aiRespToItems(doc.ID, aiResp.Items)
	if err := s.itemsRepo.CreateBatch(ctx, docItems); err != nil {
		return nil, fmt.Errorf("failed to save document items: %w", err)
	}

	return doc, nil
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
	if doc.Status != models.DocumentStatusCompleted {
		return nil, "", ErrExcelNotReady
	}

	fields, items, err := s.loadFieldsAndItems(ctx, docID)
	if err != nil {
		return nil, "", err
	}

	excelData, err := s.excelService.GenerateFromData(fields, items)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate excel: %w", err)
	}

	name := strings.TrimSuffix(doc.OriginalName, filepath.Ext(doc.OriginalName)) + ".xlsx"
	return excelData, name, nil
}

func (s *DocumentService) GetAllDocuments(ctx context.Context) ([]models.Document, error) {
	return s.docRepo.GetAll(ctx)
}

func (s *DocumentService) GetStats(ctx context.Context) (total int64, completed int64, err error) {
	total, err = s.docRepo.GetTotalCount(ctx)
	if err != nil {
		return
	}
	completed, err = s.docRepo.GetCountByStatus(ctx, models.DocumentStatusCompleted)
	return
}

func (s *DocumentService) GetAIData(ctx context.Context, docID uint, userID uint) (*ai.AIResponse, error) {
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

	fields, items, err := s.loadFieldsAndItems(ctx, docID)
	if err != nil {
		return nil, err
	}

	return fieldsAndItemsToAIResp(fields, items), nil
}

func (s *DocumentService) UpdateAIData(ctx context.Context, docID uint, userID uint, aiResp *ai.AIResponse) error {
	doc, err := s.docRepo.GetByID(ctx, docID)
	if err != nil {
		return err
	}
	if doc == nil {
		return ErrDocumentNotFound
	}
	if doc.UserID != userID {
		return ErrAccessDenied
	}

	// Update fields
	existing, err := s.fieldsRepo.GetByDocumentID(ctx, docID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrAIDataNotReady
	}

	updated := aiRespToFields(docID, aiResp)
	updated.ID = existing.ID
	updated.CreatedAt = existing.CreatedAt
	if err := s.fieldsRepo.Update(ctx, updated); err != nil {
		return err
	}

	// Replace items: delete all then recreate
	if err := s.itemsRepo.DeleteByDocumentID(ctx, docID); err != nil {
		return err
	}
	newItems := aiRespToItems(docID, aiResp.Items)
	return s.itemsRepo.CreateBatch(ctx, newItems)
}

func (s *DocumentService) DownloadXML(ctx context.Context, docID uint, userID uint) ([]byte, string, error) {
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

	fields, items, err := s.loadFieldsAndItems(ctx, docID)
	if err != nil {
		return nil, "", err
	}

	xmlData, err := generateXML(fieldsAndItemsToAIResp(fields, items))
	if err != nil {
		return nil, "", err
	}

	xmlName := strings.TrimSuffix(doc.OriginalName, filepath.Ext(doc.OriginalName)) + ".xml"
	return xmlData, xmlName, nil
}

// loadFieldsAndItems fetches both tables for a document.
func (s *DocumentService) loadFieldsAndItems(ctx context.Context, docID uint) (*models.DocumentFields, []models.DocumentItem, error) {
	fields, err := s.fieldsRepo.GetByDocumentID(ctx, docID)
	if err != nil {
		return nil, nil, err
	}
	if fields == nil {
		return nil, nil, ErrAIDataNotReady
	}

	items, err := s.itemsRepo.GetByDocumentID(ctx, docID)
	if err != nil {
		return nil, nil, err
	}

	return fields, items, nil
}

// aiRespToFields maps ai.AIResponse.Fields → models.DocumentFields.
func aiRespToFields(docID uint, resp *ai.AIResponse) *models.DocumentFields {
	f := resp.Fields
	return &models.DocumentFields{
		DocumentID:        docID,
		DocumentType:      resp.DocumentType,
		DeclarationNumber: strVal(f, "declaration_number"),
		Date:              strVal(f, "date"),
		Sender:            strVal(f, "sender"),
		Receiver:          strVal(f, "receiver"),
		CountryOrigin:     strVal(f, "country_origin"),
		CountryDest:       strVal(f, "country_dest"),
		Currency:          strVal(f, "currency"),
		TotalValue:        floatVal(f, "total_value"),
		CustomsValue:      floatVal(f, "customs_value"),
	}
}

// aiRespToItems maps []map → []models.DocumentItem.
func aiRespToItems(docID uint, raw []map[string]interface{}) []models.DocumentItem {
	items := make([]models.DocumentItem, 0, len(raw))
	for _, m := range raw {
		item := models.DocumentItem{
			DocumentID:  docID,
			Number:      intVal(m, "number"),
			HSCode:      strVal(m, "hs_code"),
			Description: strVal(m, "description"),
			Quantity:    floatVal(m, "quantity"),
			Unit:        strVal(m, "unit"),
			WeightNet:   floatVal(m, "weight_net"),
			WeightGross: floatVal(m, "weight_gross"),
			Value:       floatVal(m, "value"),
			DutyRate:    strVal(m, "duty_rate"),
			VATRate:     strVal(m, "vat_rate"),
		}
		items = append(items, item)
	}
	return items
}

// fieldsAndItemsToAIResp converts DB models back to the API format used by the frontend.
func fieldsAndItemsToAIResp(f *models.DocumentFields, items []models.DocumentItem) *ai.AIResponse {
	fields := map[string]interface{}{
		"declaration_number": f.DeclarationNumber,
		"date":               f.Date,
		"sender":             f.Sender,
		"receiver":           f.Receiver,
		"country_origin":     f.CountryOrigin,
		"country_dest":       f.CountryDest,
		"currency":           f.Currency,
		"total_value":        f.TotalValue,
		"customs_value":      f.CustomsValue,
	}

	raw := make([]map[string]interface{}, 0, len(items))
	for _, it := range items {
		raw = append(raw, map[string]interface{}{
			"number":       it.Number,
			"hs_code":      it.HSCode,
			"description":  it.Description,
			"quantity":     it.Quantity,
			"unit":         it.Unit,
			"weight_net":   it.WeightNet,
			"weight_gross": it.WeightGross,
			"value":        it.Value,
			"duty_rate":    it.DutyRate,
			"vat_rate":     it.VATRate,
		})
	}

	return &ai.AIResponse{
		DocumentType: f.DocumentType,
		Fields:       fields,
		Items:        raw,
	}
}

// ---- helper converters ----

func strVal(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok && v != nil {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func floatVal(m map[string]interface{}, key string) float64 {
	v, ok := m[key]
	if !ok || v == nil {
		return 0
	}
	switch n := v.(type) {
	case float64:
		return n
	case float32:
		return float64(n)
	case int:
		return float64(n)
	case int64:
		return float64(n)
	case string:
		f, _ := strconv.ParseFloat(n, 64)
		return f
	}
	return 0
}

func intVal(m map[string]interface{}, key string) int {
	v, ok := m[key]
	if !ok || v == nil {
		return 0
	}
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	case int64:
		return int(n)
	case string:
		i, _ := strconv.Atoi(n)
		return i
	}
	return 0
}

// ---- XML generation ----

func generateXML(resp *ai.AIResponse) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	buf.WriteString("<Document>\n")
	buf.WriteString(fmt.Sprintf("  <Type>%s</Type>\n", xmlEsc(resp.DocumentType)))
	buf.WriteString("  <Fields>\n")
	for k, v := range resp.Fields {
		buf.WriteString(fmt.Sprintf("    <Field name=%q>%s</Field>\n", k, xmlEsc(fmt.Sprintf("%v", v))))
	}
	buf.WriteString("  </Fields>\n")
	if len(resp.Items) > 0 {
		buf.WriteString("  <Items>\n")
		for _, item := range resp.Items {
			buf.WriteString("    <Item>\n")
			for k, v := range item {
				buf.WriteString(fmt.Sprintf("      <Field name=%q>%s</Field>\n", k, xmlEsc(fmt.Sprintf("%v", v))))
			}
			buf.WriteString("    </Item>\n")
		}
		buf.WriteString("  </Items>\n")
	}
	buf.WriteString("</Document>")
	return buf.Bytes(), nil
}

func xmlEsc(s string) string {
	var b bytes.Buffer
	_ = xml.EscapeText(&b, []byte(s))
	return b.String()
}
