package handlers

import (
	"cmp"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/sashaaro/url-shortener/internal/adapters"
	"github.com/sashaaro/url-shortener/internal/domain"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"slices"
)

type HTTPHandlers struct {
	urlRepo          domain.URLRepository
	genShortURLToken domain.GenShortURLToken
	logger           zap.SugaredLogger
	conn             *pgx.Conn
}

func NewHTTPHandlers(urlRepo domain.URLRepository, genShortURLToken domain.GenShortURLToken, logger zap.SugaredLogger, conn *pgx.Conn) *HTTPHandlers {
	return &HTTPHandlers{
		urlRepo:          urlRepo,
		genShortURLToken: genShortURLToken,
		logger:           logger,
		conn:             conn,
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
	err = r.urlRepo.Add(request.Context(), key, *originURL, adapters.MustUserIDFromReq(request))
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

func (r *HTTPHandlers) getShortHandler(writer http.ResponseWriter, request *http.Request) {
	hashkey := chi.URLParam(request, "hash")
	originURL, err := r.urlRepo.GetByHash(request.Context(), hashkey)
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
	err = r.urlRepo.Add(request.Context(), key, *originURL, adapters.MustUserIDFromReq(request))
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

type ShortenBatchItem struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

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
		u, err := url.Parse(item.OriginalURL)
		if err != nil {
			http.Error(w, "Invalid url", http.StatusBadRequest)
			return
		}
		originURLs = append(originURLs, domain.BatchItem{
			HashKey: r.genShortURLToken(),
			URL:     *u,
		})
	}

	err = r.urlRepo.BatchAdd(request.Context(), originURLs, adapters.MustUserIDFromReq(request))
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
	err := r.conn.Ping(request.Context())
	if err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (r *HTTPHandlers) getMyUrls(w http.ResponseWriter, request *http.Request) {
	list, err := r.urlRepo.GetByUser(request.Context(), adapters.MustUserIDFromReq(request))
	if err != nil {
		r.logger.Debug("cannot get urls", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(list)
	if err != nil {
		r.logger.Debug("cannot encode response JSON", zap.Error(err))
	}
}

func CreateServeMux(urlRepo domain.URLRepository, logger zap.SugaredLogger, conn *pgx.Conn) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	handlers := NewHTTPHandlers(urlRepo, adapters.GenBase64ShortURLToken, logger, conn)

	r.Post("/", WithAuth(gzipHandle(WithLogging(logger, handlers.createShortHandler))))
	r.Get("/{hash}", WithAuth(gzipHandle(WithLogging(logger, handlers.getShortHandler))))
	r.Post("/api/shorten", WithAuth(gzipHandle(WithLogging(logger, handlers.shorten))))
	r.Post("/api/shorten/batch", WithAuth(gzipHandle(WithLogging(logger, handlers.batchShorten))))
	r.Get("/api/user/urls", WithAuth(gzipHandle(WithLogging(logger, handlers.getMyUrls))))
	r.Get("/ping", handlers.ping)

	return r
}
