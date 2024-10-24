// Package internal - кишки
package internal

// Config основной экземляр конфига приложения
var Config = config{}

// config конфиг приложения
type config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	GrpcPort        int    `env:"GRPC_PORT"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	JwtSecret       string `env:"JWT_SECRET"`
	EnableHTTPS     bool   `env:"ENABLE_HTTPS"`
	TrustedSubnet   string `env:"TRUSTED_SUBNET"`
}
