package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sashaaro/url-shortener/internal"
	"github.com/sashaaro/url-shortener/internal/adapters"
	"github.com/sashaaro/url-shortener/internal/domain"
	"go.uber.org/zap"
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
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	originURL, err := url.Parse(string(b))
	if err != nil {
		http.Error(writer, "invalid url", http.StatusBadRequest)
		return
	}
	key := r.genShortURLToken()
	r.urlRepo.Add(key, *originURL)

	writer.WriteHeader(http.StatusCreated)

	_, _ = writer.Write([]byte(createPublicURL(key)))
}

func createPublicURL(key domain.HashKey) string {
	return internal.Config.BaseURL + "/" + key
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

type ShortenRequest struct {
	URL string `json:"url"`
}
type ShortenResponse struct {
	Result string `json:"result"`
}

func (r *HTTPHandlers) shorten(w http.ResponseWriter, request *http.Request) {
	var req ShortenRequest
	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		adapters.Logger.Debug("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	originURL, err := url.Parse(req.URL)
	if err != nil {
		http.Error(w, "Invalid url", http.StatusBadRequest)
		return
	}

	key := r.genShortURLToken()
	r.urlRepo.Add(key, *originURL)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(ShortenResponse{Result: createPublicURL(key)})
	if err != nil {
		adapters.Logger.Debug("cannot encode response JSON", zap.Error(err))
	}
}

func CreateServeMux(urlRepo domain.URLRepository) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	handlers := NewHTTPHandlers(urlRepo, adapters.GenBase64ShortURLToken)
	r.Post("/", WithLogging(adapters.Logger, handlers.createShortHandler))
	r.Get("/{hash}", WithLogging(adapters.Logger, handlers.getShortHandler))
	r.Post("/api/shorten", gzipHandle(WithLogging(adapters.Logger, handlers.shorten)))

	return r
}
