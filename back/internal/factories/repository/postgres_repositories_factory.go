package repository

import (
	"keden-service/back/internal/repositories/database/postgres/company"
	"keden-service/back/internal/repositories/database/postgres/document"
	"keden-service/back/internal/repositories/database/postgres/role"
	"keden-service/back/internal/repositories/database/postgres/subscription"
	"keden-service/back/internal/repositories/database/postgres/token"
	"keden-service/back/internal/repositories/database/postgres/user"

	"gorm.io/gorm"
)

func NewPostgresRepositoryFactory(db *gorm.DB) *PostgresRepositoryFactory {
	return &PostgresRepositoryFactory{
		UserRepository:         user.NewUserRepository(db),
		RoleRepository:         role.NewRoleRepository(db),
		CompanyRepository:      company.NewCompanyRepository(db),
		SubscriptionRepository: subscription.NewSubscriptionRepository(db),
		DocumentRepository:     document.NewDocumentRepository(db),
		RefreshTokenRepository: token.NewRefreshTokenRepository(db),
	}
}

type PostgresRepositoryFactory struct {
	UserRepository         *user.UserRepository
	RoleRepository         *role.RoleRepository
	CompanyRepository      *company.CompanyRepository
	SubscriptionRepository *subscription.SubscriptionRepository
	DocumentRepository     *document.DocumentRepository
	RefreshTokenRepository *token.RefreshTokenRepository
}
