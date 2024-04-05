package handlers

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/sashaaro/url-shortener/internal"
	"github.com/sashaaro/url-shortener/internal/adapters"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestIteration2(t *testing.T) {
	httpClient := http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}}

	t.Run("create short url, pass through short url", func(t *testing.T) {
		testServer := httptest.NewServer(CreateServeMux(adapters.NewMemURLRepository()))
		defer testServer.Close()
		internal.Config.BaseURL = testServer.URL
		resp, err := httpClient.Post(testServer.URL, "text/plain", strings.NewReader(`https://github.com`))
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		u, err := url.Parse(string(b))
		require.NoError(t, err)

		resp, err = httpClient.Get(u.String())
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
		require.Equal(t, "https://github.com", resp.Header.Get("Location"))
	})

	t.Run("create short url use POST /shorten, pass through short url", func(t *testing.T) {
		urlRepo := adapters.NewMemURLRepository()
		//urlRepo := adapters.NewFileURLRepository("/tmp/short-url-db.json", urlRepo)
		testServer := httptest.NewServer(CreateServeMux(
			urlRepo,
		))
		//defer urlRepo.Close()
		defer testServer.Close()
		internal.Config.BaseURL = testServer.URL

		resp, err := httpClient.Post(testServer.URL+"/api/shorten", "application/json", strings.NewReader(`{"url": "https://yandex.ru"}`))
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		require.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		var shortenRes ShortenResponse
		err = json.NewDecoder(resp.Body).Decode(&shortenRes)
		require.NoError(t, err)
		u, err := url.Parse(shortenRes.Result)
		require.NoError(t, err)

		resp, err = httpClient.Get(u.String())
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
		require.Equal(t, "https://yandex.ru", resp.Header.Get("Location"))

		var buf bytes.Buffer
		g := gzip.NewWriter(&buf)
		_, err = g.Write([]byte(`{"url": "https://yandex.ru"}`))
		require.NoError(t, err)
		g.Close()

		req, err := http.NewRequest("POST", testServer.URL+"/api/shorten", &buf)
		require.NoError(t, err)
		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Set("Content-Type", "application/json")

		resp, err = httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusCreated, resp.StatusCode)
	})
}
