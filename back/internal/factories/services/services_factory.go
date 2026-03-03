package services

import (
	"keden-service/back/cmd/app/structures"
	"keden-service/back/internal/factories/repository"
	"keden-service/back/internal/pkg/rabbitmq"
	"keden-service/back/internal/services/ai"
	"keden-service/back/internal/services/auth"
	"keden-service/back/internal/services/company"
	"keden-service/back/internal/services/document"
	"keden-service/back/internal/services/excel"
	"keden-service/back/internal/services/subscription"
)

func NewServiceFactory(
	pgRepo *repository.PostgresRepositoryFactory,
	data *structures.AppData,
	rabbit *rabbitmq.AmqpPubSub,
) *ServiceFactory {
	s := &ServiceFactory{}

	s.AuthService = auth.NewAuthService(
		pgRepo.UserRepository,
		pgRepo.RoleRepository,
		pgRepo.CompanyRepository,
		pgRepo.RefreshTokenRepository,
		data.JWTConfig,
	)
	s.CompanyService = company.NewCompanyService(pgRepo.UserRepository, pgRepo.CompanyRepository)
	s.SubscriptionService = subscription.NewSubscriptionService(pgRepo.SubscriptionRepository)
	s.DocumentService = document.NewDocumentService(
		pgRepo.DocumentRepository,
		rabbit,
	)
	s.AIService = ai.NewAIService(data.AIConfig)
	s.ExcelService = excel.NewExcelService()

	return s
}

type ServiceFactory struct {
	AuthService         *auth.AuthService
	CompanyService      *company.CompanyService
	SubscriptionService *subscription.SubscriptionService
	DocumentService     *document.DocumentService
	AIService           *ai.AIService
	ExcelService        *excel.ExcelService
}
