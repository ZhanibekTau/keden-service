package structures

type JWTConfig struct {
	Secret             string `mapstructure:"JWT_SECRET"`
	AccessTokenExpiry  string `mapstructure:"JWT_ACCESS_TOKEN_EXPIRY"`
	RefreshTokenExpiry string `mapstructure:"JWT_REFRESH_TOKEN_EXPIRY"`
}
