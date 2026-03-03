package configs

import (
	"keden-service/back/internal/configs/structures"
	"keden-service/back/internal/pkg/config"

	"github.com/sirupsen/logrus"
)

func InitDbConfig() (*structures.DbConfig, error) {
	dbConfig := &structures.DbConfig{}
	if err := config.InitConfig(dbConfig); err != nil {
		return nil, err
	}
	logrus.Info("DB config initialized")
	return dbConfig, nil
}

func InitRabbitConfig() (*structures.RabbitConfig, error) {
	rabbitConfig := &structures.RabbitConfig{}
	if err := config.InitConfig(rabbitConfig); err != nil {
		return nil, err
	}
	logrus.Info("RabbitMQ config initialized")
	return rabbitConfig, nil
}

func InitJWTConfig() (*structures.JWTConfig, error) {
	jwtConfig := &structures.JWTConfig{}
	if err := config.InitConfig(jwtConfig); err != nil {
		return nil, err
	}
	logrus.Info("JWT config initialized")
	return jwtConfig, nil
}

func InitS3Config() (*structures.S3Config, error) {
	s3Config := &structures.S3Config{}
	if err := config.InitConfig(s3Config); err != nil {
		return nil, err
	}
	logrus.Info("S3 config initialized")
	return s3Config, nil
}

func InitAIConfig() (*structures.AIConfig, error) {
	aiConfig := &structures.AIConfig{}
	if err := config.InitConfig(aiConfig); err != nil {
		return nil, err
	}
	logrus.Info("AI config initialized")
	return aiConfig, nil
}

func InitAdminConfig() (*structures.AdminConfig, error) {
	adminConfig := &structures.AdminConfig{}
	if err := config.InitConfig(adminConfig); err != nil {
		return nil, err
	}
	logrus.Info("Admin config initialized")
	return adminConfig, nil
}

func InitBaseConfig() (*config.BaseConfig, error) {
	baseConfig := &config.BaseConfig{HandlerTimeout: 30}
	if err := config.InitConfig(baseConfig); err != nil {
		return nil, err
	}
	logrus.Info("Base config initialized")
	return baseConfig, nil
}
