package subscription

import (
	"context"
	"errors"
	"keden-service/back/internal/models"
	companyRepo "keden-service/back/internal/repositories/database/postgres/company"
	subRepo "keden-service/back/internal/repositories/database/postgres/subscription"
	"time"
)

var (
	ErrSubscriptionNotFound = errors.New("subscription not found")
	ErrAlreadyPending       = errors.New("you already have a pending subscription request")
	ErrAlreadyActive        = errors.New("you already have an active subscription")
	ErrCannotApprove        = errors.New("subscription is not in pending status")
	ErrCannotTransition     = errors.New("invalid status transition")
)

// SubscriptionDetail enriches Subscription with company fields for the admin view.
type SubscriptionDetail struct {
	models.Subscription
	CompanyName string `json:"company_name,omitempty"`
	LegalName   string `json:"legal_name,omitempty"`
	BIN         string `json:"bin,omitempty"`
}

// allowed status transitions: current status → allowed next statuses
var allowedTransitions = map[string][]string{
	models.SubscriptionStatusPending:     {models.SubscriptionStatusInProgress, models.SubscriptionStatusRejected},
	models.SubscriptionStatusInProgress:  {models.SubscriptionStatusInvoiceSent, models.SubscriptionStatusRejected},
	models.SubscriptionStatusInvoiceSent: {models.SubscriptionStatusActive, models.SubscriptionStatusRejected},
}

type SubscriptionService struct {
	subRepo     subRepo.ISubscriptionRepository
	companyRepo companyRepo.ICompanyRepository
}

func NewSubscriptionService(sr subRepo.ISubscriptionRepository, cr companyRepo.ICompanyRepository) *SubscriptionService {
	return &SubscriptionService{subRepo: sr, companyRepo: cr}
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
		if sub.Status == models.SubscriptionStatusPending ||
			sub.Status == models.SubscriptionStatusInProgress ||
			sub.Status == models.SubscriptionStatusInvoiceSent {
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

// GetActiveRequests returns all open subscription requests enriched with company data.
func (s *SubscriptionService) GetActiveRequests(ctx context.Context) ([]SubscriptionDetail, error) {
	subs, err := s.subRepo.GetActiveRequests(ctx)
	if err != nil {
		return nil, err
	}

	details := make([]SubscriptionDetail, 0, len(subs))
	for _, sub := range subs {
		d := SubscriptionDetail{Subscription: sub}
		if sub.User != nil && sub.User.AccountType == "company" {
			company, _ := s.companyRepo.GetByUserID(ctx, sub.UserID)
			if company != nil {
				d.CompanyName = company.CompanyName
				d.LegalName = company.LegalName
				d.BIN = company.BIN
			}
		}
		details = append(details, d)
	}
	return details, nil
}

// UpdateStatus moves a subscription to the next status following the allowed transition map.
func (s *SubscriptionService) UpdateStatus(ctx context.Context, subscriptionID uint, adminID uint, newStatus string, comment string) error {
	sub, err := s.subRepo.GetByID(ctx, subscriptionID)
	if err != nil {
		return err
	}
	if sub == nil {
		return ErrSubscriptionNotFound
	}

	allowed, ok := allowedTransitions[sub.Status]
	if !ok {
		return ErrCannotTransition
	}
	valid := false
	for _, s := range allowed {
		if s == newStatus {
			valid = true
			break
		}
	}
	if !valid {
		return ErrCannotTransition
	}

	now := time.Now()
	sub.Status = newStatus
	sub.AdminComment = comment
	sub.ApprovedByID = &adminID
	sub.ApprovedAt = &now

	if newStatus == models.SubscriptionStatusActive {
		endDate := now.AddDate(0, 1, 0)
		sub.StartDate = &now
		sub.EndDate = &endDate
	}

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
