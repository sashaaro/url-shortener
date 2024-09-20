package internal

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"log"
	"net/url"
	"os"
	"strings"
)

// InitConfig инициализация конфигурации
func InitConfig() {
	serverAddress := flag.String("a", "", "listen address")
	baseURL := flag.String("b", "", "base url")
	databaseDSN := flag.String("d", "", "database dsn")
	enableHttps := flag.String("s", "", "Enable https")
	configFile := flag.String("c", "", "Config file")

	fileStoragePath := flag.String("f", "/tmp/short-url-db.json", "file path")

	flag.Parse()

	if err := env.Parse(&Config); err != nil {
		fmt.Printf("%+v\n", err)
	}
	if Config.ServerAddress == "" {
		Config.ServerAddress = *serverAddress
	}
	if Config.BaseURL == "" {
		Config.BaseURL = *baseURL
	}
	if Config.FileStoragePath == "" {
		Config.FileStoragePath = *fileStoragePath
	}

	if Config.ServerAddress == "" {
		Config.ServerAddress = ":8080"
	}
	if Config.BaseURL == "" {
		Config.BaseURL = "http://localhost:8080"
	}

	if enableHttps != nil && len(*enableHttps) > 0 && *enableHttps != "0" {
		Config.EnableHTTPS = true
	}

	_, err := url.Parse(Config.BaseURL)
	if err != nil {
		log.Fatal("invalid base url: ", err)
	}

	if Config.DatabaseDSN == "" {
		Config.DatabaseDSN = *databaseDSN
	}

	Config.DatabaseDSN = strings.TrimSpace(Config.DatabaseDSN)

	if Config.JwtSecret == "" {
		Config.JwtSecret = "secret"
	}

	parseFromConfigFile(configFile)
}

func parseFromConfigFile(configFile *string) {
	if configFile == nil || *configFile == "" {
		return
	}
	jsonFile, err := os.Open(*configFile)
	if err != nil {
		log.Fatal("no config file: ", err)
	}
	defer jsonFile.Close()

	c := &jsonConfig{}
	err = json.NewDecoder(jsonFile).Decode(c)
	if err != nil {
		log.Fatal("invalid config file: ", err)
	}

	if c.EnableHTTPS {
		Config.EnableHTTPS = true
	}
	if c.ServerAddress != "" {
		Config.ServerAddress = c.ServerAddress
	}
	if c.BaseURL != "" {
		Config.BaseURL = c.BaseURL
	}
	if c.DatabaseDSN != "" {
		Config.DatabaseDSN = c.DatabaseDSN
	}
}

type jsonConfig struct {
	ServerAddress   string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDSN     string `json:"database_dsn"`
	EnableHTTPS     bool   `json:"enable_https"`
}
