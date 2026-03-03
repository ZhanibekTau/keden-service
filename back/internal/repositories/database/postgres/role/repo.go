package role

import (
	"context"
	"keden-service/back/internal/models"

	"gorm.io/gorm"
)

type IRoleRepository interface {
	GetByName(ctx context.Context, name string) (*models.Role, error)
	GetAll(ctx context.Context) ([]models.Role, error)
}

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) GetByName(ctx context.Context, name string) (*models.Role, error) {
	var role models.Role
	result := r.db.WithContext(ctx).Where("name = ?", name).First(&role)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &role, nil
}

func (r *RoleRepository) GetAll(ctx context.Context) ([]models.Role, error) {
	var roles []models.Role
	result := r.db.WithContext(ctx).Find(&roles)
	return roles, result.Error
}
