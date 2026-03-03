package structures

import (
	"keden-service/back/internal/pkg/rabbitmq"

	"gorm.io/gorm"
)

type AppClients struct {
	DbClient *gorm.DB
	RabbitMq *rabbitmq.AmqpPubSub
}
