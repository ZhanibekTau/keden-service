package company

import (
	"context"
	"keden-service/back/internal/models"

	"gorm.io/gorm"
)

type ICompanyRepository interface {
	Create(ctx context.Context, company *models.Company) error
	GetByID(ctx context.Context, id uint) (*models.Company, error)
	GetByUserID(ctx context.Context, userID uint) (*models.Company, error)
	Update(ctx context.Context, company *models.Company) error
}

type CompanyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) Create(ctx context.Context, company *models.Company) error {
	// Omit("User") prevents GORM from nullifying user_id when User pointer is nil (BelongsTo behavior)
	return r.db.WithContext(ctx).Omit("User").Create(company).Error
}

func (r *CompanyRepository) GetByID(ctx context.Context, id uint) (*models.Company, error) {
	var company models.Company
	result := r.db.WithContext(ctx).First(&company, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &company, nil
}

func (r *CompanyRepository) GetByUserID(ctx context.Context, userID uint) (*models.Company, error) {
	var company models.Company
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&company)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &company, nil
}

func (r *CompanyRepository) Update(ctx context.Context, company *models.Company) error {
	return r.db.WithContext(ctx).Save(company).Error
}
