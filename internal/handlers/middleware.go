package handlers

import (
	"compress/gzip"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sashaaro/url-shortener/internal"
	"github.com/sashaaro/url-shortener/internal/adapters"
	"github.com/sashaaro/url-shortener/internal/utils"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Claims - payload jwt токена
type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID `json:"user_id,omitempty"`
}

// WithLogging - добавление лога запроса
func WithLogging(logger zap.SugaredLogger, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		logger.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.status,
			"duration", duration,
			"size", responseData.size,
		)
	}
}

type responseData struct {
	status int
	size   int
}

// добавляем реализацию http.ResponseWriter
type loggingResponseWriter struct {
	http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
	responseData        *responseData
}

// декоратор для ResponseWriter.Write
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

// декоратор для ResponseWriter.WriteHeader
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

func gzipHandle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			uncompressed, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = uncompressed
			//nolint:errcheck
			defer uncompressed.Close()
		}

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			//nolint:errcheck
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	}
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

// декоратор для Write
func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

// декоратор для gzip.Reader().Read
func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// декоратор для Reader.Close()
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

// JwtTTL - время жизни токена по умолчанию
const JwtTTL = 15 * time.Minute

// BuildJWTString - создание токена
func BuildJWTString(secretKey string, userID uuid.UUID) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(JwtTTL)),
		},
		// собственное утверждение
		UserID: userID,
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}

func fetchUserIDFromToken(secretKey string, tokenStr string) (uuid.UUID, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
	if err != nil {
		return uuid.Nil, err
	}

	return claims.UserID, nil
}

// WithAuth провекра jwt и добавление пользователя в context
func WithAuth(authRequired bool, h http.HandlerFunc) http.HandlerFunc {
	hostname := utils.Must(url.Parse(internal.Config.BaseURL)).Hostname()
	return func(w http.ResponseWriter, r *http.Request) {
		authCookie, _ := r.Cookie("access_token")

		var userID uuid.UUID
		var accessToken string

		if authCookie != nil && authCookie.Value != "" {
			accessToken = authCookie.Value
		}

		authHeader := r.Header.Get("Authorization")
		if accessToken == "" && authHeader != "" {
			accessToken = strings.TrimPrefix(authHeader, "Bearer ")
		}

		if accessToken != "" {
			var err error
			userID, err = fetchUserIDFromToken(internal.Config.JwtSecret, accessToken)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte("invalid access token"))
				return
			}
		}
		if userID == uuid.Nil {
			if authRequired {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte("Unauthorized"))
				return
			}
			userID = uuid.New()
		}

		if accessToken == "" {
			var err error
			accessToken, err = BuildJWTString(internal.Config.JwtSecret, userID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    accessToken,
			Path:     "/",
			Domain:   hostname,
			Expires:  time.Now().Add(JwtTTL),
			Secure:   true,
			HttpOnly: true,
		})

		h.ServeHTTP(w, r.WithContext(adapters.UserIDToCxt(r.Context(), userID)))
	}
}
