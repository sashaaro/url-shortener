package handlers

import (
	"github.com/sashaaro/url-shortener/internal/adapters"
	"github.com/sashaaro/url-shortener/internal/domain"
	"log"
	"net/http"
)

func Example() {
	urlRepo := adapters.NewMemURLRepository()
	logger := adapters.CreateLogger()

	mux := CreateServeMux(domain.NewShortenerService(urlRepo), logger, nil)

	log.Fatal(http.ListenAndServe(":8080", mux))
	// use with server
}
