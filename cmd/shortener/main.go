package main

import (
	"flag"
	"github.com/sashaaro/url-shortener/internal"
	"github.com/sashaaro/url-shortener/internal/infra"
	"log"
	"net/http"
)

func main() {
	flag.Parse()

	urlRepo := infra.NewMemURLRepository()
	log.Fatal(http.ListenAndServe(*internal.HTTPAddr, infra.CreateServeMux(urlRepo)))
}
