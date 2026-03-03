package subscription

import (
	"context"
	"errors"
	"keden-service/back/internal/models"
	subRepo "keden-service/back/internal/repositories/database/postgres/subscription"
	"time"
)

var (
	ErrSubscriptionNotFound = errors.New("subscription not found")
	ErrAlreadyPending       = errors.New("you already have a pending subscription request")
	ErrAlreadyActive        = errors.New("you already have an active subscription")
	ErrCannotApprove        = errors.New("subscription is not in pending status")
)

type SubscriptionService struct {
	subRepo subRepo.ISubscriptionRepository
}

func NewSubscriptionService(sr subRepo.ISubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{subRepo: sr}
}

func (s *SubscriptionService) RequestSubscription(ctx context.Context, userID uint) (*models.Subscription, error) {
	active, err := s.subRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if active != nil {
		return nil, ErrAlreadyActive
	}

	subs, err := s.subRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	for _, sub := range subs {
		if sub.Status == models.SubscriptionStatusPending {
			return nil, ErrAlreadyPending
		}
	}

	now := time.Now()
	sub := &models.Subscription{
		UserID:      userID,
		Status:      models.SubscriptionStatusPending,
		Amount:      12990,
		RequestedAt: now,
	}

	if err := s.subRepo.Create(ctx, sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *SubscriptionService) GetCurrentSubscription(ctx context.Context, userID uint) (*models.Subscription, error) {
	active, err := s.subRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if active != nil {
		return active, nil
	}

	subs, err := s.subRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(subs) > 0 {
		return &subs[0], nil
	}
	return nil, nil
}

func (s *SubscriptionService) GetSubscriptionHistory(ctx context.Context, userID uint) ([]models.Subscription, error) {
	return s.subRepo.GetByUserID(ctx, userID)
}

func (s *SubscriptionService) GetPendingRequests(ctx context.Context) ([]models.Subscription, error) {
	return s.subRepo.GetPendingRequests(ctx)
}

func (s *SubscriptionService) ApproveSubscription(ctx context.Context, subscriptionID uint, adminID uint, comment string) error {
	sub, err := s.subRepo.GetByID(ctx, subscriptionID)
	if err != nil {
		return err
	}
	if sub == nil {
		return ErrSubscriptionNotFound
	}
	if sub.Status != models.SubscriptionStatusPending {
		return ErrCannotApprove
	}

	now := time.Now()
	endDate := now.AddDate(0, 1, 0)

	sub.Status = models.SubscriptionStatusActive
	sub.StartDate = &now
	sub.EndDate = &endDate
	sub.ApprovedAt = &now
	sub.ApprovedByID = &adminID
	sub.AdminComment = comment

	return s.subRepo.Update(ctx, sub)
}

func (s *SubscriptionService) RejectSubscription(ctx context.Context, subscriptionID uint, adminID uint, comment string) error {
	sub, err := s.subRepo.GetByID(ctx, subscriptionID)
	if err != nil {
		return err
	}
	if sub == nil {
		return ErrSubscriptionNotFound
	}
	if sub.Status != models.SubscriptionStatusPending {
		return ErrCannotApprove
	}

	now := time.Now()
	sub.Status = models.SubscriptionStatusRejected
	sub.ApprovedAt = &now
	sub.ApprovedByID = &adminID
	sub.AdminComment = comment

	return s.subRepo.Update(ctx, sub)
}

func (s *SubscriptionService) CheckActiveSubscription(ctx context.Context, userID uint) (bool, error) {
	active, err := s.subRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		return false, err
	}
	return active != nil, nil
}

func (s *SubscriptionService) GetActiveCount(ctx context.Context) (int64, error) {
	return s.subRepo.GetActiveCount(ctx)
}

func (s *SubscriptionService) GetPendingCount(ctx context.Context) (int64, error) {
	return s.subRepo.GetPendingCount(ctx)
}
