package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/sashaaro/url-shortener/internal"
	"github.com/sashaaro/url-shortener/internal/infra"
	"log"
	"net/http"
)

func main() {
	initConfig()

	urlRepo := infra.NewMemURLRepository()
	log.Fatal(http.ListenAndServe(internal.Config.ServerAddress, infra.CreateServeMux(urlRepo)))
}

func initConfig() {
	serverAddress := flag.String("a", "", "listen address")
	baseURL := flag.String("b", "", "base url")
	flag.Parse()

	if err := env.Parse(&internal.Config); err != nil {
		fmt.Printf("%+v\n", err)
	}
	if internal.Config.ServerAddress == "" {
		internal.Config.ServerAddress = *serverAddress
	}
	if internal.Config.BaseURL == "" {
		internal.Config.BaseURL = *baseURL
	}

	if internal.Config.ServerAddress == "" {
		internal.Config.ServerAddress = ":8080"
	}
	if internal.Config.BaseURL == "" {
		internal.Config.BaseURL = "http://localhost:8080"
	}
}
