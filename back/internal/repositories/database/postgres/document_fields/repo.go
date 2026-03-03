package document_fields

import (
	"context"
	"keden-service/back/internal/models"

	"gorm.io/gorm"
)

type IDocumentFieldsRepository interface {
	Create(ctx context.Context, f *models.DocumentFields) error
	GetByDocumentID(ctx context.Context, docID uint) (*models.DocumentFields, error)
	Update(ctx context.Context, f *models.DocumentFields) error
}

type DocumentFieldsRepository struct {
	db *gorm.DB
}

func NewDocumentFieldsRepository(db *gorm.DB) *DocumentFieldsRepository {
	return &DocumentFieldsRepository{db: db}
}

func (r *DocumentFieldsRepository) Create(ctx context.Context, f *models.DocumentFields) error {
	return r.db.WithContext(ctx).Create(f).Error
}

func (r *DocumentFieldsRepository) GetByDocumentID(ctx context.Context, docID uint) (*models.DocumentFields, error) {
	var f models.DocumentFields
	err := r.db.WithContext(ctx).Where("document_id = ?", docID).First(&f).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &f, nil
}

func (r *DocumentFieldsRepository) Update(ctx context.Context, f *models.DocumentFields) error {
	return r.db.WithContext(ctx).Save(f).Error
}
