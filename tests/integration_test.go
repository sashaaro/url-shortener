package tests

import (
	"github.com/sashaaro/url-shortener/internal/infra"
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
		testServer := httptest.NewServer(infra.CreateServeMux(infra.NewMemURLRepository()))
		defer testServer.Close()
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
}
