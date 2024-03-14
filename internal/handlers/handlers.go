package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sashaaro/url-shortener/internal"
	"github.com/sashaaro/url-shortener/internal/adapters"
	"github.com/sashaaro/url-shortener/internal/domain"
	"io"
	"net/http"
	"net/url"
)

type HTTPHandlers struct {
	urlRepo          domain.URLRepository
	genShortURLToken domain.GenShortURLToken
}

func NewHTTPHandlers(urlRepo domain.URLRepository, genShortURLToken domain.GenShortURLToken) *HTTPHandlers {
	return &HTTPHandlers{
		urlRepo:          urlRepo,
		genShortURLToken: genShortURLToken,
	}
}

func (r *HTTPHandlers) createShortHandler(writer http.ResponseWriter, request *http.Request) {
	b, err := io.ReadAll(request.Body)
	if err != nil {
		return
	}

	originURL, err := url.Parse(string(b))
	if err != nil {
		http.Error(writer, "Invalid url", http.StatusBadRequest)
		return
	}
	key := r.genShortURLToken()
	r.urlRepo.Add(key, *originURL)

	writer.WriteHeader(http.StatusCreated)

	writer.Write([]byte(internal.Config.BaseURL + "/" + key))
}

func (r *HTTPHandlers) getShortHandler(writer http.ResponseWriter, request *http.Request) {
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
	handlers := NewHTTPHandlers(urlRepo, adapters.GenBase64ShortURLToken)
	r.Post("/", handlers.createShortHandler)
	r.Get("/{hash}", handlers.getShortHandler)

	return r
}
