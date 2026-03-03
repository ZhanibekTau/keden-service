package company

import (
	"context"
	"errors"
	"keden-service/back/internal/models"
	companyRepo "keden-service/back/internal/repositories/database/postgres/company"
	userRepo "keden-service/back/internal/repositories/database/postgres/user"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrCompanyNotFound = errors.New("company not found")
	ErrWrongPassword   = errors.New("current password is incorrect")
)

type UpdateProfileRequest struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Phone         string `json:"phone"`
	CompanyName   string `json:"company_name"`
	LegalName     string `json:"legal_name"`
	ContactPerson string `json:"contact_person"`
}

type ProfileResponse struct {
	User    *models.User    `json:"user"`
	Company *models.Company `json:"company,omitempty"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

type CompanyService struct {
	userRepo    userRepo.IUserRepository
	companyRepo companyRepo.ICompanyRepository
}

func NewCompanyService(ur userRepo.IUserRepository, cr companyRepo.ICompanyRepository) *CompanyService {
	return &CompanyService{
		userRepo:    ur,
		companyRepo: cr,
	}
}

func (s *CompanyService) GetProfile(ctx context.Context, userID uint) (*ProfileResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	resp := &ProfileResponse{User: user}

	if user.AccountType == "company" {
		company, err := s.companyRepo.GetByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}
		resp.Company = company
	}

	return resp, nil
}

func (s *CompanyService) UpdateProfile(ctx context.Context, userID uint, req UpdateProfileRequest) (*ProfileResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	resp := &ProfileResponse{User: user}

	if user.AccountType == "company" {
		company, err := s.companyRepo.GetByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}
		if company != nil {
			if req.CompanyName != "" {
				company.CompanyName = req.CompanyName
			}
			if req.LegalName != "" {
				company.LegalName = req.LegalName
			}
			if req.ContactPerson != "" {
				company.ContactPerson = req.ContactPerson
			}
			if err := s.companyRepo.Update(ctx, company); err != nil {
				return nil, err
			}
			resp.Company = company
		}
	}

	return resp, nil
}

func (s *CompanyService) ChangePassword(ctx context.Context, userID uint, req ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		return ErrWrongPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hash)
	return s.userRepo.Update(ctx, user)
}

func (s *CompanyService) GetAllClients(ctx context.Context) ([]models.User, error) {
	return s.userRepo.GetAll(ctx)
}

func (s *CompanyService) UpdateUserStatus(ctx context.Context, userID uint, isActive bool) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}
	user.IsActive = isActive
	return s.userRepo.Update(ctx, user)
}

func (s *CompanyService) GetStats(ctx context.Context) (int64, error) {
	return s.userRepo.GetStats(ctx)
}
