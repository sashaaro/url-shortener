package main

import (
	"github.com/sashaaro/url-shortener/internal/config"
	"github.com/sashaaro/url-shortener/internal/infra"
	"log"
	"net/http"
)

func main() {
	urlRepo := infra.NewMemURLRepository()
	log.Fatal(http.ListenAndServe(":"+config.HTTPPort, infra.CreateServeMux(urlRepo)))
}
