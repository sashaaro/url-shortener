package internal

// основной экземляр конфига приложения
var Config = config{}

// конфиг приложения
type config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	JwtSecret       string `env:"JWT_SECRET"`
}
