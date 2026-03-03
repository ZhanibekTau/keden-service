package document

import (
	"context"
	"keden-service/back/internal/models"

	"gorm.io/gorm"
)

type IDocumentRepository interface {
	Create(ctx context.Context, doc *models.Document) error
	GetByID(ctx context.Context, id uint) (*models.Document, error)
	GetByUserID(ctx context.Context, userID uint) ([]models.Document, error)
	GetAll(ctx context.Context) ([]models.Document, error)
	Update(ctx context.Context, doc *models.Document) error
	UpdateStatus(ctx context.Context, id uint, status string, errorMsg string) error
	GetTotalCount(ctx context.Context) (int64, error)
	GetCountByStatus(ctx context.Context, status string) (int64, error)
}

type DocumentRepository struct {
	db *gorm.DB
}

func NewDocumentRepository(db *gorm.DB) *DocumentRepository {
	return &DocumentRepository{db: db}
}

func (r *DocumentRepository) Create(ctx context.Context, doc *models.Document) error {
	return r.db.WithContext(ctx).Omit("User").Create(doc).Error
}

func (r *DocumentRepository) GetByID(ctx context.Context, id uint) (*models.Document, error) {
	var doc models.Document
	result := r.db.WithContext(ctx).Preload("User").First(&doc, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &doc, nil
}

func (r *DocumentRepository) GetByUserID(ctx context.Context, userID uint) ([]models.Document, error) {
	var docs []models.Document
	result := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&docs)
	return docs, result.Error
}

func (r *DocumentRepository) GetAll(ctx context.Context) ([]models.Document, error) {
	var docs []models.Document
	result := r.db.WithContext(ctx).Preload("User").Order("created_at DESC").Find(&docs)
	return docs, result.Error
}

func (r *DocumentRepository) Update(ctx context.Context, doc *models.Document) error {
	return r.db.WithContext(ctx).Save(doc).Error
}

func (r *DocumentRepository) UpdateStatus(ctx context.Context, id uint, status string, errorMsg string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if errorMsg != "" {
		updates["error_message"] = errorMsg
	}
	return r.db.WithContext(ctx).Model(&models.Document{}).Where("id = ?", id).Updates(updates).Error
}

func (r *DocumentRepository) GetTotalCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Document{}).Count(&count).Error
	return count, err
}

func (r *DocumentRepository) GetCountByStatus(ctx context.Context, status string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Document{}).Where("status = ?", status).Count(&count).Error
	return count, err
}
