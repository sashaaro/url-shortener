package main

import (
	"github.com/sashaaro/url-shortener/internal"
	"github.com/sashaaro/url-shortener/internal/infra"
	"log"
	"net/http"
)

func main() {
	urlRepo := infra.NewMemURLRepository()
	log.Fatal(http.ListenAndServe(*internal.HTTPAddr, infra.CreateServeMux(urlRepo, *internal.BaseURL)))
}
