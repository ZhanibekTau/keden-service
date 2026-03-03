package structures

type S3Config struct {
	Endpoint  string `mapstructure:"S3_ENDPOINT"`
	AccessKey string `mapstructure:"S3_ACCESS_KEY"`
	SecretKey string `mapstructure:"S3_SECRET_KEY"`
	Bucket    string `mapstructure:"S3_BUCKET"`
	Region    string `mapstructure:"S3_REGION"`
	UseSSL    bool   `mapstructure:"S3_USE_SSL"`
}
