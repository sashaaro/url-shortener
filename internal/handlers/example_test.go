package handlers

import (
	"github.com/sashaaro/url-shortener/internal/adapters"
	"log"
	"net/http"
)

func Example() {
	urlRepo := adapters.NewMemURLRepository()
	logger := adapters.CreateLogger()

	mux := CreateServeMux(urlRepo, logger, nil)

	log.Fatal(http.ListenAndServe(":8080", mux))
	// use with server
}
