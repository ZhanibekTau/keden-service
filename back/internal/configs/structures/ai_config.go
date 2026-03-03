package structures

type AIConfig struct {
	ServiceURL string `mapstructure:"AI_SERVICE_URL"`
	Timeout    string `mapstructure:"AI_SERVICE_TIMEOUT"`
}
