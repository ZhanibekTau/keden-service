package config

import (
	"os"

	"github.com/spf13/viper"
)

func ReadEnv() error {
	p := "../back/.env"

	if _, err := os.Stat(p); err == nil {
		viper.SetConfigFile(p)
		viper.AutomaticEnv()
		return viper.ReadInConfig()
	}

	// Fallback: just use env variables
	viper.AutomaticEnv()
	return nil
}

func InitConfig(cfg interface{}) error {
	return viper.Unmarshal(cfg)
}
