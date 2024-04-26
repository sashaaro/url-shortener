package handlers

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"github.com/sashaaro/url-shortener/internal"
	"github.com/sashaaro/url-shortener/internal/adapters"
	"github.com/sashaaro/url-shortener/internal/infra"
	"github.com/sashaaro/url-shortener/internal/utils"
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

	internal.InitConfig()

	logger := adapters.CreateLogger()

	urlRepo := adapters.NewMemURLRepository()

	if internal.Config.DatabaseDSN != "" {
		conn := infra.CreatePgxConn()
		//nolint:errcheck
		defer conn.Close(context.Background())
		_, err := conn.Exec(context.Background(), "TRUNCATE TABLE urls")
		require.NoError(t, err)
		urlRepo = adapters.NewPgURLRepository(conn)
	}

	testServer := httptest.NewServer(CreateServeMux(urlRepo, logger, nil))
	defer testServer.Close()
	internal.Config.BaseURL = testServer.URL

	t.Run("create short url, pass through short url", func(t *testing.T) {
		resp, err := httpClient.Post(testServer.URL, "text/plain", strings.NewReader(`https://github.com`))
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		require.NotEmpty(t, resp.Header.Get("Authorization"))
		authorization := resp.Header.Get("Authorization")

		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		u, err := url.Parse(string(b))
		require.NoError(t, err)

		resp, err = httpClient.Get(u.String())
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
		require.Equal(t, "https://github.com", resp.Header.Get("Location"))

		if internal.Config.DatabaseDSN != "" {
			req := utils.Must(http.NewRequest("DELETE", testServer.URL+"/api/user/urls", strings.NewReader(`["`+strings.TrimPrefix(u.Path, "/")+`"]`)))
			req.Header.Set("Authorization", authorization)
			resp, err = httpClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			require.Equal(t, http.StatusAccepted, resp.StatusCode)

			resp, err = httpClient.Get(u.String())
			require.NoError(t, err)
			defer resp.Body.Close()
			require.Equal(t, http.StatusGone, resp.StatusCode)
		}

		resp, err = httpClient.Get(testServer.URL + "/NoExistShortUrl")
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusNotFound, resp.StatusCode)

		if internal.Config.DatabaseDSN != "" { // check unique key for database
			resp, err = httpClient.Post(testServer.URL, "text/plain", strings.NewReader(`https://github.com`))
			require.NoError(t, err)
			defer resp.Body.Close()
			require.Equal(t, http.StatusConflict, resp.StatusCode)
		}
	})

	t.Run("create short url use POST /shorten, pass through short url", func(t *testing.T) {

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
		require.NotEmpty(t, resp.Header.Get("Authorization"))
		authorization := resp.Header.Get("Authorization")

		var buf bytes.Buffer
		g := gzip.NewWriter(&buf)
		_, err = g.Write([]byte(`{"url": "https://rambler.ru"}`))
		require.NoError(t, err)
		g.Close()

		req, err := http.NewRequest("POST", testServer.URL+"/api/shorten", &buf)
		require.NoError(t, err)
		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authorization)

		resp, err = httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		req, err = http.NewRequest("GET", testServer.URL+"/api/user/urls", nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", authorization)
		resp, err = httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
