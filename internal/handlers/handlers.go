package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sashaaro/url-shortener/internal"
	"github.com/sashaaro/url-shortener/internal/domain"
	"io"
	"net/http"
	"net/url"
)

func GenShortURLToken() string {
	length := 8
	bufSize := length*6/8 + 1
	buf := make([]byte, bufSize)
	n, err := rand.Read(buf)
	if err != nil || n != bufSize {
		panic(fmt.Errorf("error while retriving random data: %d %v", n, err.Error()))
	}
	return base64.URLEncoding.EncodeToString(buf)[:length]
}

type HttpHandlers struct {
	urlRepo domain.URLRepository
}

func NewHttpHandlers(urlRepo domain.URLRepository) *HttpHandlers {
	return &HttpHandlers{urlRepo: urlRepo}
}

func (r *HttpHandlers) createShortHandler(writer http.ResponseWriter, request *http.Request) {
	b, err := io.ReadAll(request.Body)
	if err != nil {
		return
	}

	originURL, err := url.Parse(string(b))
	if err != nil {
		http.Error(writer, "Invalid url", http.StatusBadRequest)
		return
	}
	key := GenShortURLToken()
	r.urlRepo.Add(key, *originURL)

	writer.WriteHeader(http.StatusCreated)

	writer.Write([]byte(internal.Config.BaseURL + "/" + key))
}

func (r *HttpHandlers) getShortHandler(writer http.ResponseWriter, request *http.Request) {
	hashkey := chi.URLParam(request, "hash")
	originURL, ok := r.urlRepo.GetByHash(hashkey)
	if !ok {
		http.Error(writer, "Short url not found", http.StatusNotFound)
		return
	}
	http.Redirect(writer, request, originURL.String(), http.StatusTemporaryRedirect)
}

func CreateServeMux(urlRepo domain.URLRepository) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	handlers := NewHttpHandlers(urlRepo)
	r.Post("/", handlers.createShortHandler)
	r.Get("/{hash}", handlers.getShortHandler)

	return r
}
