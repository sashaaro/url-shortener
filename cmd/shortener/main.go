package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/sashaaro/url-shortener/internal"
	"github.com/sashaaro/url-shortener/internal/adapters"
	"github.com/sashaaro/url-shortener/internal/handlers"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func main() {
	internal.InitConfig()

	logger := adapters.CreateLogger()

	var conn *pgx.Conn
	var err error
	conn, err = pgx.Connect(context.Background(), internal.Config.DatabaseDSN)
	if err != nil {
		logger.Warn("can't connect to database", zap.Error(err))
	}
	defer conn.Close(context.Background())

	urlRepo := adapters.NewFileURLRepository(internal.Config.FileStoragePath, adapters.NewMemURLRepository(), logger)
	log.Fatal(http.ListenAndServe(internal.Config.ServerAddress, handlers.CreateServeMux(urlRepo, logger, conn)))
}
