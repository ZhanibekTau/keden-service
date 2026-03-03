package subscription

import (
	"context"
	"keden-service/back/internal/models"
	"time"

	"gorm.io/gorm"
)

type ISubscriptionRepository interface {
	Create(ctx context.Context, sub *models.Subscription) error
	GetByID(ctx context.Context, id uint) (*models.Subscription, error)
	GetByUserID(ctx context.Context, userID uint) ([]models.Subscription, error)
	GetActiveByUserID(ctx context.Context, userID uint) (*models.Subscription, error)
	GetActiveRequests(ctx context.Context) ([]models.Subscription, error)
	Update(ctx context.Context, sub *models.Subscription) error
	GetActiveCount(ctx context.Context) (int64, error)
	GetPendingCount(ctx context.Context) (int64, error)
	ExpireOverdue(ctx context.Context) error
}

type SubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(ctx context.Context, sub *models.Subscription) error {
	return r.db.WithContext(ctx).Omit("User", "ApprovedBy").Create(sub).Error
}

func (r *SubscriptionRepository) GetByID(ctx context.Context, id uint) (*models.Subscription, error) {
	var sub models.Subscription
	result := r.db.WithContext(ctx).Preload("User").Preload("User.Role").First(&sub, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &sub, nil
}

func (r *SubscriptionRepository) GetByUserID(ctx context.Context, userID uint) ([]models.Subscription, error) {
	var subs []models.Subscription
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&subs)
	return subs, result.Error
}

func (r *SubscriptionRepository) GetActiveByUserID(ctx context.Context, userID uint) (*models.Subscription, error) {
	var sub models.Subscription
	now := time.Now()
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ? AND end_date > ?", userID, models.SubscriptionStatusActive, now).
		Order("end_date DESC").
		First(&sub)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &sub, nil
}

// GetActiveRequests returns all subscriptions not yet finalized (pending, in_progress, invoice_sent).
func (r *SubscriptionRepository) GetActiveRequests(ctx context.Context) ([]models.Subscription, error) {
	var subs []models.Subscription
	result := r.db.WithContext(ctx).
		Preload("User").
		Preload("User.Role").
		Where("status IN ?", []string{
			models.SubscriptionStatusPending,
			models.SubscriptionStatusInProgress,
			models.SubscriptionStatusInvoiceSent,
		}).
		Order("requested_at ASC").
		Find(&subs)
	return subs, result.Error
}

func (r *SubscriptionRepository) Update(ctx context.Context, sub *models.Subscription) error {
	return r.db.WithContext(ctx).Omit("User", "ApprovedBy").Save(sub).Error
}

func (r *SubscriptionRepository) GetActiveCount(ctx context.Context) (int64, error) {
	var count int64
	now := time.Now()
	err := r.db.WithContext(ctx).Model(&models.Subscription{}).
		Where("status = ? AND end_date > ?", models.SubscriptionStatusActive, now).
		Count(&count).Error
	return count, err
}

func (r *SubscriptionRepository) GetPendingCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Subscription{}).
		Where("status IN ?", []string{
			models.SubscriptionStatusPending,
			models.SubscriptionStatusInProgress,
			models.SubscriptionStatusInvoiceSent,
		}).
		Count(&count).Error
	return count, err
}

func (r *SubscriptionRepository) ExpireOverdue(ctx context.Context) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.Subscription{}).
		Where("status = ? AND end_date <= ?", models.SubscriptionStatusActive, now).
		Update("status", models.SubscriptionStatusExpired).Error
}
