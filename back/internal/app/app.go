package app

import (
	"context"
	"fmt"
	"keden-service/back/cmd/app/structures"
	"keden-service/back/internal/configs"
	"keden-service/back/internal/factories/handlers"
	"keden-service/back/internal/factories/repository"
	"keden-service/back/internal/factories/services"
	"keden-service/back/internal/middleware"
	"keden-service/back/internal/pkg/config"
	"keden-service/back/internal/pkg/database"
	localRabbit "keden-service/back/internal/pkg/rabbitmq"

	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewApp() (*MyApp, error) {
	newApp := &MyApp{}

	if err := newApp.PrepareConfigs(); err != nil {
		return nil, err
	}
	if err := newApp.PrepareComponents(); err != nil {
		return nil, err
	}

	return newApp, nil
}

type MyApp struct {
	appData         *structures.AppData
	clients         *structures.AppClients
	handlersFactory *handlers.HandlerFactory
	repoFactory     *repository.PostgresRepositoryFactory
	servicesFactory *services.ServiceFactory
	router          *gin.Engine
}

func (a *MyApp) PrepareConfigs() error {
	if err := config.ReadEnv(); err != nil {
		return err
	}

	a.appData = &structures.AppData{}

	baseConfig, err := configs.InitBaseConfig()
	if err != nil {
		return err
	}
	a.appData.BaseConfig = baseConfig

	dbConfig, err := configs.InitDbConfig()
	if err != nil {
		return err
	}
	a.appData.DbConfig = dbConfig

	rabbitConfig, err := configs.InitRabbitConfig()
	if err != nil {
		return err
	}
	a.appData.RabbitConfig = rabbitConfig

	jwtConfig, err := configs.InitJWTConfig()
	if err != nil {
		return err
	}
	a.appData.JWTConfig = jwtConfig

	aiConfig, err := configs.InitAIConfig()
	if err != nil {
		return err
	}
	a.appData.AIConfig = aiConfig

	adminConfig, err := configs.InitAdminConfig()
	if err != nil {
		return err
	}
	a.appData.AdminConfig = adminConfig

	return nil
}

func (a *MyApp) PrepareComponents() error {
	a.clients = &structures.AppClients{}

	if err := a.initDb(); err != nil {
		return err
	}

	if err := a.initRabbitClient(); err != nil {
		logrus.Warnf("RabbitMQ init failed (non-fatal): %v", err)
	}

	a.repoFactory = repository.NewPostgresRepositoryFactory(a.clients.DbClient)
	a.servicesFactory = services.NewServiceFactory(a.repoFactory, a.appData, a.clients.RabbitMq)
	a.handlersFactory = handlers.NewHandlerFactory(a.servicesFactory)

	a.setupRouter()

	return nil
}

func (a *MyApp) initDb() error {
	dbConfig := database.DbConfig{
		Driver:             database.Postgres,
		Host:               a.appData.DbConfig.Host,
		User:               a.appData.DbConfig.User,
		Password:           a.appData.DbConfig.Password,
		Db:                 a.appData.DbConfig.DbName,
		Port:               a.appData.DbConfig.Port,
		SslMode:            false,
		MaxOpenConnections: a.appData.DbConfig.MaxOpenConnections,
		MaxIdleConnections: a.appData.DbConfig.MaxIdleConnections,
		Logging:            a.appData.DbConfig.Logging,
	}

	dbClient, err := database.GetGormConnection(dbConfig)
	if err != nil {
		return err
	}

	a.clients.DbClient = dbClient

	if err := RunMigrations(dbClient); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	if err := SeedAdmin(dbClient, a.appData.AdminConfig); err != nil {
		logrus.Warnf("Admin seed failed: %v", err)
	}

	return nil
}

func (a *MyApp) initRabbitClient() error {
	amqpURL := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		a.appData.RabbitConfig.RabbitLogin,
		a.appData.RabbitConfig.RabbitPassword,
		a.appData.RabbitConfig.RabbitHost,
		a.appData.RabbitConfig.RabbitPort,
	)

	amqClient, err := localRabbit.NewAmqpPubSub(amqp.ConnectionConfig{
		AmqpURI: amqpURL,
	})
	if err != nil {
		return err
	}

	a.clients.RabbitMq = amqClient
	return nil
}

func (a *MyApp) setupRouter() {
	router := gin.Default()
	router.Use(middleware.CORS())

	api := router.Group("/api/v1")

	// Public routes
	authRoutes := api.Group("/auth")
	{
		authRoutes.POST("/register", a.handlersFactory.AuthHandler.Register)
		authRoutes.POST("/login", a.handlersFactory.AuthHandler.Login)
		authRoutes.POST("/refresh", a.handlersFactory.AuthHandler.Refresh)
	}

	// Authenticated routes
	protected := api.Group("")
	protected.Use(middleware.JWTAuth(a.servicesFactory.AuthService))
	{
		protected.POST("/auth/logout", a.handlersFactory.AuthHandler.Logout)

		// Company profile
		protected.GET("/company/profile", a.handlersFactory.CompanyHandler.GetProfile)
		protected.PUT("/company/profile", a.handlersFactory.CompanyHandler.UpdateProfile)
		protected.PUT("/company/password", a.handlersFactory.CompanyHandler.ChangePassword)

		// Subscription
		protected.GET("/subscription/current", a.handlersFactory.SubscriptionHandler.GetCurrentSubscription)
		protected.POST("/subscription/request", a.handlersFactory.SubscriptionHandler.RequestSubscription)
		protected.GET("/subscription/history", a.handlersFactory.SubscriptionHandler.GetSubscriptionHistory)

		// Documents (require active subscription)
		docs := protected.Group("/documents")
		docs.Use(middleware.CheckActiveSubscription(a.servicesFactory.SubscriptionService))
		{
			docs.POST("/upload", a.handlersFactory.DocumentHandler.Upload)
			docs.GET("", a.handlersFactory.DocumentHandler.GetDocuments)
			docs.GET("/:id", a.handlersFactory.DocumentHandler.GetDocumentByID)
			docs.GET("/:id/download", a.handlersFactory.DocumentHandler.DownloadExcel)
		}
	}

	// Admin routes
	admin := api.Group("/admin")
	admin.Use(middleware.JWTAuth(a.servicesFactory.AuthService))
	admin.Use(middleware.RequireRole("admin"))
	{
		admin.GET("/stats", a.handlersFactory.AdminHandler.GetStats)
		admin.GET("/companies", a.handlersFactory.CompanyHandler.GetAllClients)
		admin.PUT("/companies/:id/status", a.handlersFactory.CompanyHandler.UpdateUserStatus)
		admin.GET("/subscriptions/pending", a.handlersFactory.SubscriptionHandler.GetPendingRequests)
		admin.POST("/subscriptions/:id/approve", a.handlersFactory.SubscriptionHandler.ApproveSubscription)
		admin.POST("/subscriptions/:id/reject", a.handlersFactory.SubscriptionHandler.RejectSubscription)
		admin.GET("/documents", a.handlersFactory.DocumentHandler.GetAllDocuments)
	}

	a.router = router
}

func (a *MyApp) RunRestServer() error {
	addr := a.appData.BaseConfig.ServerAddress
	logrus.Infof("Starting REST server on %s", addr)
	return a.router.Run(addr)
}

func (a *MyApp) RunConsumer() error {
	if a.clients.RabbitMq == nil {
		logrus.Warn("RabbitMQ not available, consumer not started")
		return nil
	}

	queueConfig := localRabbit.NewConsumerTopicDurableConfig(
		"document.processing",
		"keden.documents.exchange",
		"keden.documents.processing.queue",
		1,
	)

	if err := a.clients.RabbitMq.RegisterHandler(
		a.handlersFactory.DocumentProcessorHandler.ConsumeMessage,
		queueConfig...,
	); err != nil {
		return err
	}

	logrus.Info("Document processing consumer started")
	if err := a.clients.RabbitMq.Consume(context.Background()); err != nil {
		return err
	}

	defer a.clients.RabbitMq.CloseConnection()
	return nil
}
