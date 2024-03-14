package internal

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
)

func InitConfig() {
	serverAddress := flag.String("a", "", "listen address")
	baseURL := flag.String("b", "", "base url")
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

	if Config.ServerAddress == "" {
		Config.ServerAddress = ":8080"
	}
	if Config.BaseURL == "" {
		Config.BaseURL = "http://localhost:8080"
	}
}
