package document_items

import (
	"context"
	"keden-service/back/internal/models"

	"gorm.io/gorm"
)

type IDocumentItemRepository interface {
	CreateBatch(ctx context.Context, items []models.DocumentItem) error
	GetByDocumentID(ctx context.Context, docID uint) ([]models.DocumentItem, error)
	DeleteByDocumentID(ctx context.Context, docID uint) error
}

type DocumentItemRepository struct {
	db *gorm.DB
}

func NewDocumentItemRepository(db *gorm.DB) *DocumentItemRepository {
	return &DocumentItemRepository{db: db}
}

func (r *DocumentItemRepository) CreateBatch(ctx context.Context, items []models.DocumentItem) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&items).Error
}

func (r *DocumentItemRepository) GetByDocumentID(ctx context.Context, docID uint) ([]models.DocumentItem, error) {
	var items []models.DocumentItem
	err := r.db.WithContext(ctx).
		Where("document_id = ?", docID).
		Order("number ASC").
		Find(&items).Error
	return items, err
}

func (r *DocumentItemRepository) DeleteByDocumentID(ctx context.Context, docID uint) error {
	return r.db.WithContext(ctx).Where("document_id = ?", docID).Delete(&models.DocumentItem{}).Error
}
