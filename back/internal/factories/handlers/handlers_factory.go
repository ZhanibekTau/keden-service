package handlers

import (
	"keden-service/back/internal/factories/services"
	adminHandler "keden-service/back/internal/handlers/http/admin"
	authHandler "keden-service/back/internal/handlers/http/auth"
	companyHandler "keden-service/back/internal/handlers/http/company"
	documentHandler "keden-service/back/internal/handlers/http/document"
	subscriptionHandler "keden-service/back/internal/handlers/http/subscription"
)

func NewHandlerFactory(svc *services.ServiceFactory) *HandlerFactory {
	h := &HandlerFactory{}

	h.AuthHandler = authHandler.NewAuthHandler(svc.AuthService)
	h.CompanyHandler = companyHandler.NewCompanyHandler(svc.CompanyService)
	h.SubscriptionHandler = subscriptionHandler.NewSubscriptionHandler(svc.SubscriptionService)
	h.DocumentHandler = documentHandler.NewDocumentHandler(svc.DocumentService)
	h.AdminHandler = adminHandler.NewAdminHandler(
		svc.CompanyService,
		svc.SubscriptionService,
		svc.DocumentService,
	)

	return h
}

type HandlerFactory struct {
	AuthHandler         *authHandler.AuthHandler
	CompanyHandler      *companyHandler.CompanyHandler
	SubscriptionHandler *subscriptionHandler.SubscriptionHandler
	DocumentHandler     *documentHandler.DocumentHandler
	AdminHandler        *adminHandler.AdminHandler
}
