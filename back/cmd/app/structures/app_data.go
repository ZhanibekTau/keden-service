package structures

import (
	"keden-service/back/internal/configs/structures"
	"keden-service/back/internal/pkg/config"
)

type AppData struct {
	BaseConfig  *config.BaseConfig
	DbConfig    *structures.DbConfig
	JWTConfig   *structures.JWTConfig
	AIConfig    *structures.AIConfig
	AdminConfig *structures.AdminConfig
}
