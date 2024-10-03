// Package handlers Модуль http хендлеров
package handlers

import (
	"cmp"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sashaaro/url-shortener/internal"
	"github.com/sashaaro/url-shortener/internal/adapters"
	"github.com/sashaaro/url-shortener/internal/domain"
	"github.com/sashaaro/url-shortener/internal/utils"
	"go.uber.org/zap"
	"io"
	"net"
	"net/http"
	"net/url"
	"slices"
	"strings"
)

// HTTPHandlers основные хендлеры
type HTTPHandlers struct {
	service          *domain.ShortenerService
	genShortURLToken domain.GenShortURLToken
	logger           zap.SugaredLogger
	pool             *pgxpool.Pool
}

// NewHTTPHandlers конструктор
func NewHTTPHandlers(
	service *domain.ShortenerService,
	genShortURLToken domain.GenShortURLToken,
	logger zap.SugaredLogger,
	pool *pgxpool.Pool,
) *HTTPHandlers {
	return &HTTPHandlers{
		service:          service,
		genShortURLToken: genShortURLToken,
		logger:           logger,
		pool:             pool,
	}
}

func (r *HTTPHandlers) createShortHandler(writer http.ResponseWriter, request *http.Request) {
	b, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	originURL, err := url.Parse(string(b))
	if err != nil || !strings.HasPrefix(originURL.Scheme, "http") {
		http.Error(writer, "invalid url", http.StatusBadRequest)
		return
	}
	key := r.genShortURLToken()
	err = r.service.CreateShort(request.Context(), key, *originURL, adapters.MustUserIDFromReq(request))
	var dupErr *domain.ErrURLAlreadyExists
	if errors.As(err, &dupErr) {
		writer.WriteHeader(http.StatusConflict)
		_, _ = writer.Write([]byte(adapters.CreatePublicURL(dupErr.HashKey)))
		return
	}
	if err != nil {
		r.logger.Debug("cannot batch add urls", zap.Error(err))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	_, _ = writer.Write([]byte(adapters.CreatePublicURL(key)))
}

func (r *HTTPHandlers) getOriginLinkHandler(writer http.ResponseWriter, request *http.Request) {
	hashkey := chi.URLParam(request, "hash")
	originURL, err := r.service.GetOriginLink(request.Context(), hashkey)
	if errors.Is(err, domain.ErrURLDeleted) {
		writer.WriteHeader(http.StatusGone)
		return
	}
	if err != nil {
		r.logger.Debug("cannot get url by hash", zap.Error(err))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if originURL == nil {
		http.Error(writer, "short url not found", http.StatusNotFound)
		return
	}
	http.Redirect(writer, request, originURL.String(), http.StatusTemporaryRedirect)
}

// ShortenRequest - запрос на укорочение ссылоки
type ShortenRequest struct {
	URL string `json:"url"`
}

// ShortenResponse -ответ на укорочение ссылоки
type ShortenResponse struct {
	Result string `json:"result"`
}

func (r *HTTPHandlers) shorten(w http.ResponseWriter, request *http.Request) {
	var req ShortenRequest
	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		r.logger.Debug("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	originURL, err := url.Parse(req.URL)
	if err != nil {
		http.Error(w, "Invalid url", http.StatusBadRequest)
		return
	}

	key := r.genShortURLToken()
	err = r.service.CreateShort(request.Context(), key, *originURL, adapters.MustUserIDFromReq(request))
	var dupErr *domain.ErrURLAlreadyExists
	if errors.As(err, &dupErr) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		err = json.NewEncoder(w).Encode(ShortenResponse{Result: adapters.CreatePublicURL(dupErr.HashKey)})
		if err != nil {
			r.logger.Debug("cannot encode response JSON", zap.Error(err))
		}
		return
	}

	if err != nil {
		r.logger.Debug("cannot add url", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(ShortenResponse{Result: adapters.CreatePublicURL(key)})
	if err != nil {
		r.logger.Debug("cannot encode response JSON", zap.Error(err))
	}
}

// ShortenBatchItem - запрос на укорочение нескольких ссылок
type ShortenBatchItem struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// ShortenItemRes - ответ на укорочение нескольких ссылок
type ShortenItemRes struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func (r *HTTPHandlers) batchShorten(w http.ResponseWriter, request *http.Request) {
	var req []ShortenBatchItem
	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		r.logger.Debug("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(req) == 0 {
		w.WriteHeader(http.StatusCreated)
		return
	}

	slices.SortFunc(req, func(a, b ShortenBatchItem) int {
		return cmp.Compare(a.CorrelationID, b.CorrelationID)
	})

	originURLs := make([]domain.BatchItem, 0, len(req))
	for _, item := range req {
		var u *url.URL
		u, err = url.Parse(item.OriginalURL)
		if err != nil {
			http.Error(w, "Invalid url", http.StatusBadRequest)
			return
		}
		originURLs = append(originURLs, domain.BatchItem{
			HashKey: r.genShortURLToken(),
			URL:     *u,
		})
	}

	err = r.service.BatchAdd(request.Context(), originURLs, adapters.MustUserIDFromReq(request))
	var dupErr *domain.ErrURLAlreadyExists
	if errors.As(err, &dupErr) {
		w.WriteHeader(http.StatusConflict)
		return
	}
	if err != nil {
		r.logger.Debug("cannot batch add urls", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := make([]ShortenItemRes, 0, len(originURLs))
	for i, item := range originURLs {
		resp = append(resp, ShortenItemRes{
			CorrelationID: req[i].CorrelationID,
			ShortURL:      adapters.CreatePublicURL(item.HashKey),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		r.logger.Debug("cannot encode response JSON", zap.Error(err))
	}
}

func (r *HTTPHandlers) ping(w http.ResponseWriter, request *http.Request) {
	err := r.pool.Ping(request.Context())
	if err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (r *HTTPHandlers) getMyUrls(w http.ResponseWriter, request *http.Request) {
	list, err := r.service.GetByUser(request.Context(), adapters.MustUserIDFromReq(request))
	if err != nil {
		r.logger.Debug("cannot get urls", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(list)
	if err != nil {
		r.logger.Debug("cannot encode response JSON", zap.Error(err))
	}
}

func (r *HTTPHandlers) deleteUrls(w http.ResponseWriter, request *http.Request) {
	keys := []string{}
	err := json.NewDecoder(request.Body).Decode(&keys)
	if err != nil {
		r.logger.Debug("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	keys = utils.Filter(keys, func(key string) bool {
		return len([]rune(key)) > 0
	})
	if len(keys) != 0 {
		_, err = r.service.DeleteByUser(request.Context(), keys, adapters.MustUserIDFromReq(request))
		if err != nil {
			r.logger.Error("cannot delete urls", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Cannot delete urls"))
			return
		}
	}

	w.WriteHeader(http.StatusAccepted)
}

func (r *HTTPHandlers) stats(w http.ResponseWriter, request *http.Request) {
	stats, err := r.service.Stats(request.Context())

	if err != nil {
		r.logger.Error("cannot get stats", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Cannot get stats"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(stats)
}

// CreateServeMux - создание основных хендлеров
func CreateServeMux(service *domain.ShortenerService, logger zap.SugaredLogger, pool *pgxpool.Pool) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	handlers := NewHTTPHandlers(service, adapters.GenBase64ShortURLToken, logger, pool)

	statsHandler := WithAuth(false, gzipHandle(WithLogging(logger, handlers.stats)))
	if internal.Config.TrustedSubnet != "" {
		_, subnet, err := net.ParseCIDR(internal.Config.TrustedSubnet)
		if err != nil {
			panic(err)
		}

		statsHandler = TrustedClientMiddleware(logger, subnet)(statsHandler)
	}

	r.Post("/", WithAuth(false, gzipHandle(WithLogging(logger, handlers.createShortHandler))))
	r.Get("/{hash}", WithAuth(false, gzipHandle(WithLogging(logger, handlers.getOriginLinkHandler))))
	r.Post("/api/shorten", WithAuth(false, gzipHandle(WithLogging(logger, handlers.shorten))))
	r.Post("/api/shorten/batch", WithAuth(false, gzipHandle(WithLogging(logger, handlers.batchShorten))))
	r.Get("/api/user/urls", WithAuth(true, gzipHandle(WithLogging(logger, handlers.getMyUrls))))
	r.Delete("/api/user/urls", WithAuth(false, gzipHandle(WithLogging(logger, handlers.deleteUrls))))
	r.Get("/api/internal/stats", statsHandler)

	r.Get("/ping", handlers.ping)
	r.Mount("/debug", middleware.Profiler())

	return r
}
