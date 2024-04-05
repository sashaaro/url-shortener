package main

import (
	"github.com/sashaaro/url-shortener/internal"
	"github.com/sashaaro/url-shortener/internal/adapters"
	"github.com/sashaaro/url-shortener/internal/handlers"
	"log"
	"net/http"
)

func main() {
	internal.InitConfig()

	urlRepo := adapters.NewFileURLRepository(internal.Config.FileStoragePath, adapters.NewMemURLRepository())
	log.Fatal(http.ListenAndServe(internal.Config.ServerAddress, handlers.CreateServeMux(urlRepo)))
}
