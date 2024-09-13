package main

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sashaaro/url-shortener/internal"
	"github.com/sashaaro/url-shortener/internal/adapters"
	"github.com/sashaaro/url-shortener/internal/domain"
	"github.com/sashaaro/url-shortener/internal/handlers"
	"github.com/sashaaro/url-shortener/internal/infra"
	"github.com/sashaaro/url-shortener/internal/version"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
	"os"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	version.Build{
		BuildVersion: buildVersion,
		BuildDate:    buildDate,
		BuildCommit:  buildCommit,
	}.Print(os.Stdout)

	internal.InitConfig()

	logger := adapters.CreateLogger()

	var urlRepo domain.URLRepository

	var pool *pgxpool.Pool
	if internal.Config.DatabaseDSN != "" {
		pool = infra.CreatePgxPool()
		//nolint:errcheck
		defer pool.Close()
		urlRepo = adapters.NewPgURLRepository(pool)
	} else {
		urlRepo = adapters.NewMemURLRepository()
		if internal.Config.FileStoragePath != "" {
			urlRepo = adapters.NewFileURLRepository(internal.Config.FileStoragePath, urlRepo, logger) // wrap with file storage
		}
	}

	if internal.Config.EnableHTTPS {
		log.Fatal(
			http.Serve(
				autocert.NewListener(internal.Config.ServerAddress, "url-shortener.ru", "www.url-shortener.ru"),
				handlers.CreateServeMux(urlRepo, logger, pool),
			),
		)
	} else {
		log.Fatal(http.ListenAndServe(internal.Config.ServerAddress, handlers.CreateServeMux(urlRepo, logger, pool)))
	}
}
