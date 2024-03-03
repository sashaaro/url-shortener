package infra

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/sashaaro/url-shortener/internal/domain"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func GenShortURLToken() string {
	length := 8
	bufSize := length*6/8 + 1
	buf := make([]byte, bufSize)
	n, err := rand.Read(buf)
	if err != nil || n != bufSize {
		panic(fmt.Errorf("error while retriving random data: %d %v", n, err.Error()))
	}
	return base64.URLEncoding.EncodeToString(buf)[:length]
}

func CreateServeMux(urlRepo domain.URLRepository) *http.ServeMux {
	mux := http.NewServeMux()

	createShortHandler := func(writer http.ResponseWriter, request *http.Request) {
		b, err := io.ReadAll(request.Body)
		if err != nil {
			return
		}

		originURL, err := url.Parse(string(b))
		if err != nil {
			http.Error(writer, "Invalid url", http.StatusBadRequest)
			return
		}
		key := GenShortURLToken()
		urlRepo.Add(key, *originURL)

		writer.WriteHeader(http.StatusCreated)

		writer.Write([]byte("http://" + request.Host + "/" + key))
	}

	getShortHandler := func(writer http.ResponseWriter, request *http.Request) {
		key := strings.Trim(request.URL.Path, "/")

		originURL, ok := urlRepo.GetByHash(key)
		if !ok {
			http.Error(writer, "Short url not found", http.StatusNotFound)
			return
		}
		http.Redirect(writer, request, originURL.String(), http.StatusTemporaryRedirect)
	}

	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodPost {
			createShortHandler(writer, request)
		} else if request.Method == http.MethodGet {
			getShortHandler(writer, request)
		} else {
			http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}
