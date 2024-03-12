package infra

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

func CreateServeMux(urlRepo domain.URLRepository) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	createShortHandler := func(writer http.ResponseWriter, request *http.Request) {
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
		urlRepo.Add(key, *originURL)

		writer.WriteHeader(http.StatusCreated)

		writer.Write([]byte("http://" + request.Host + "/" + key))
	}

	getShortHandler := func(writer http.ResponseWriter, request *http.Request) {
		hashkey := chi.URLParam(request, "hash")
		originURL, ok := urlRepo.GetByHash(hashkey)
		if !ok {
			http.Error(writer, "Short url not found", http.StatusNotFound)
			return
		}
		http.Redirect(writer, request, originURL.String(), http.StatusTemporaryRedirect)
	}

	r.Post("/", createShortHandler)
	r.Get("/{hash}", getShortHandler)

	return r
}
