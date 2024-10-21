package handlers

import (
	"context"
	"github.com/sashaaro/url-shortener/internal"
	"github.com/sashaaro/url-shortener/internal/adapters"
	"github.com/sashaaro/url-shortener/internal/domain"
	"github.com/sashaaro/url-shortener/internal/infra"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func BenchmarkApplication_APIHandlerGetURL(b *testing.B) {
	httpClient := http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}}

	logger := adapters.CreateLogger()

	urlRepo := adapters.NewMemURLRepository()

	if internal.Config.DatabaseDSN != "" {
		pool := infra.CreatePgxPool()
		//nolint:errcheck
		defer pool.Close()
		_, err := pool.Exec(context.Background(), "TRUNCATE TABLE urls")
		require.NoError(b, err)
		urlRepo = adapters.NewPgURLRepository(pool)
	}

	testServer := httptest.NewServer(CreateServeMux(domain.NewShortenerService(urlRepo, adapters.GenBase64ShortURLToken), logger, nil))
	defer testServer.Close()
	internal.Config.BaseURL = testServer.URL

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resp, err := httpClient.Post(testServer.URL, "text/plain", strings.NewReader(`https://github.com`))
		require.NoError(b, err)
		defer resp.Body.Close()
		require.Equal(b, http.StatusCreated, resp.StatusCode)

	}
}
