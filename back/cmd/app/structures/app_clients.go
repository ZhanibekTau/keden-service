package structures

import (
	"gorm.io/gorm"
)

type AppClients struct {
	DbClient *gorm.DB
}
