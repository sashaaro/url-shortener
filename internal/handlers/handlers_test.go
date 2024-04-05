package handlers

import (
	"encoding/json"
	"github.com/sashaaro/url-shortener/internal"
	"github.com/sashaaro/url-shortener/internal/adapters"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestIteration2(t *testing.T) {
	t.Run("create short url, pass through short url", func(t *testing.T) {
		testServer := httptest.NewServer(CreateServeMux(adapters.NewMemURLRepository()))
		defer testServer.Close()
		internal.Config.BaseURL = testServer.URL
		resp, err := http.Post(testServer.URL, "text/plain", strings.NewReader(`https://github.com`))
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		b, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		u, err := url.Parse(string(b))
		assert.NoError(t, err)

		httpClient := http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}
		resp, err = httpClient.Get(u.String())
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
		assert.Equal(t, "https://github.com", resp.Header.Get("Location"))
	})

	t.Run("create short url use POST /shorten, pass through short url", func(t *testing.T) {
		testServer := httptest.NewServer(CreateServeMux(adapters.NewMemURLRepository()))
		defer testServer.Close()
		internal.Config.BaseURL = testServer.URL
		resp, err := http.Post(testServer.URL+"/api/shorten", "application/json", strings.NewReader(`{"url": "https://yandex.ru"}`))
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		var shortenRes ShortenResponse
		err = json.NewDecoder(resp.Body).Decode(&shortenRes)
		assert.NoError(t, err)
		u, err := url.Parse(shortenRes.Result)
		assert.NoError(t, err)

		httpClient := http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}
		resp, err = httpClient.Get(u.String())
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
		assert.Equal(t, "https://yandex.ru", resp.Header.Get("Location"))
	})
}
