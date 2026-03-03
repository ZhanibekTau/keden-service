package config

type BaseConfig struct {
	AppName        string `mapstructure:"APP_NAME"`
	Hostname       string `mapstructure:"HOSTNAME"`
	AppVersion     string `mapstructure:"APP_VERSION"`
	ServerAddress  string `mapstructure:"SERVER_ADDRESS"`
	HandlerTimeout int    `mapstructure:"HANDLER_TIMEOUT"`
}
