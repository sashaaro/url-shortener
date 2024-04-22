package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/sashaaro/url-shortener/internal"
	"github.com/sashaaro/url-shortener/internal/adapters"
	"github.com/sashaaro/url-shortener/internal/domain"
	"github.com/sashaaro/url-shortener/internal/handlers"
	"github.com/sashaaro/url-shortener/internal/infra"
	"log"
	"net/http"
)

func main() {
	internal.InitConfig()

	logger := adapters.CreateLogger()

	var urlRepo domain.URLRepository

	var conn *pgx.Conn
	if internal.Config.DatabaseDSN != "" {
		conn = infra.CreatePgxConn()
		//nolint:errcheck
		defer conn.Close(context.Background())
		urlRepo = adapters.NewPgURLRepository(conn)
	} else {
		urlRepo = adapters.NewMemURLRepository()
		if internal.Config.FileStoragePath != "" {
			urlRepo = adapters.NewFileURLRepository(internal.Config.FileStoragePath, urlRepo, logger) // wrap with file storage
		}
	}

	log.Fatal(http.ListenAndServe(internal.Config.ServerAddress, handlers.CreateServeMux(urlRepo, logger, conn)))
}
