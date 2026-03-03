package auth

import (
	"context"
	"errors"
	"keden-service/back/internal/configs/structures"
	"keden-service/back/internal/models"
	companyRepo "keden-service/back/internal/repositories/database/postgres/company"
	roleRepo "keden-service/back/internal/repositories/database/postgres/role"
	tokenRepo "keden-service/back/internal/repositories/database/postgres/token"
	userRepo "keden-service/back/internal/repositories/database/postgres/user"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailExists        = errors.New("email already registered")
	ErrBINExists          = errors.New("BIN already registered")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrAccountInactive    = errors.New("account is inactive")
	ErrRoleNotFound       = errors.New("role not found")
	ErrCompanyFieldsRequired = errors.New("company fields are required for company account type")
)

type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	Phone       string `json:"phone" binding:"required"`
	AccountType string `json:"account_type" binding:"required,oneof=individual company"`
	// Company fields (required if account_type == "company")
	CompanyName   string `json:"company_name"`
	LegalName     string `json:"legal_name"`
	BIN           string `json:"bin"`
	ContactPerson string `json:"contact_person"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AuthResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresAt    time.Time   `json:"expires_at"`
	User         models.User `json:"user"`
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type AuthService struct {
	userRepo    userRepo.IUserRepository
	roleRepo    roleRepo.IRoleRepository
	companyRepo companyRepo.ICompanyRepository
	tokenRepo   tokenRepo.IRefreshTokenRepository
	jwtConfig   *structures.JWTConfig
}

func NewAuthService(
	ur userRepo.IUserRepository,
	rr roleRepo.IRoleRepository,
	cr companyRepo.ICompanyRepository,
	tr tokenRepo.IRefreshTokenRepository,
	jwtCfg *structures.JWTConfig,
) *AuthService {
	return &AuthService{
		userRepo:    ur,
		roleRepo:    rr,
		companyRepo: cr,
		tokenRepo:   tr,
		jwtConfig:   jwtCfg,
	}
}

func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	existing, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailExists
	}

	if req.AccountType == "company" {
		if req.CompanyName == "" || req.LegalName == "" || req.BIN == "" {
			return nil, ErrCompanyFieldsRequired
		}
	}

	clientRole, err := s.roleRepo.GetByName(ctx, "client")
	if err != nil {
		return nil, err
	}
	if clientRole == nil {
		return nil, ErrRoleNotFound
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Phone:        req.Phone,
		RoleID:       clientRole.ID,
		AccountType:  req.AccountType,
		IsActive:     true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	if req.AccountType == "company" {
		company := &models.Company{
			UserID:        user.ID,
			CompanyName:   req.CompanyName,
			LegalName:     req.LegalName,
			BIN:           req.BIN,
			ContactPerson: req.ContactPerson,
		}
		if err := s.companyRepo.Create(ctx, company); err != nil {
			return nil, ErrBINExists
		}
	}

	user.Role = clientRole

	return s.generateTokenPair(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, ErrAccountInactive
	}

	return s.generateTokenPair(ctx, user)
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshTokenStr string) (*AuthResponse, error) {
	token, err := s.tokenRepo.GetByToken(ctx, refreshTokenStr)
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, ErrInvalidToken
	}

	user, err := s.userRepo.GetByID(ctx, token.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidToken
	}

	_ = s.tokenRepo.DeleteByToken(ctx, refreshTokenStr)

	return s.generateTokenPair(ctx, user)
}

func (s *AuthService) Logout(ctx context.Context, refreshTokenStr string) error {
	return s.tokenRepo.DeleteByToken(ctx, refreshTokenStr)
}

func (s *AuthService) generateTokenPair(ctx context.Context, user *models.User) (*AuthResponse, error) {
	accessExpiry, err := time.ParseDuration(s.jwtConfig.AccessTokenExpiry)
	if err != nil {
		accessExpiry = 15 * time.Minute
	}

	refreshExpiry, err := time.ParseDuration(s.jwtConfig.RefreshTokenExpiry)
	if err != nil {
		refreshExpiry = 7 * 24 * time.Hour
	}

	roleName := ""
	if user.Role != nil {
		roleName = user.Role.Name
	}

	now := time.Now()
	accessExpiresAt := now.Add(accessExpiry)

	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   roleName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   user.Email,
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenStr, err := accessToken.SignedString([]byte(s.jwtConfig.Secret))
	if err != nil {
		return nil, err
	}

	refreshTokenStr := uuid.New().String()
	refreshToken := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenStr,
		ExpiresAt: now.Add(refreshExpiry),
	}

	if err := s.tokenRepo.Create(ctx, refreshToken); err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
		ExpiresAt:    accessExpiresAt,
		User:         *user,
	}, nil
}

func (s *AuthService) ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtConfig.Secret), nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
