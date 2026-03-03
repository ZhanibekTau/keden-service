package user

import (
	"context"
	"keden-service/back/internal/models"

	"gorm.io/gorm"
)

type IUserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id uint) (*models.User, error)
	GetAll(ctx context.Context) ([]models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
	GetStats(ctx context.Context) (int64, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).Preload("Role").Where("email = ? AND deleted_at IS NULL", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).Preload("Role").Where("id = ? AND deleted_at IS NULL", id).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	result := r.db.WithContext(ctx).
		Preload("Role").
		Joins("JOIN roles ON roles.id = users.role_id").
		Where("roles.name = ? AND users.deleted_at IS NULL", "client").
		Order("users.created_at DESC").
		Find(&users)
	return users, result.Error
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Unscoped().Delete(&models.User{}, id).Error
}

func (r *UserRepository) GetStats(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.User{}).
		Joins("JOIN roles ON roles.id = users.role_id").
		Where("roles.name = ? AND users.deleted_at IS NULL", "client").
		Count(&count).Error
	return count, err
}
