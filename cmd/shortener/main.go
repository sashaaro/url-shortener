package main

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sashaaro/url-shortener/internal"
	"github.com/sashaaro/url-shortener/internal/adapters"
	"github.com/sashaaro/url-shortener/internal/domain"
	shortenerGrpc "github.com/sashaaro/url-shortener/internal/grpc"
	"github.com/sashaaro/url-shortener/internal/handlers"
	"github.com/sashaaro/url-shortener/internal/infra"
	"github.com/sashaaro/url-shortener/internal/version"
	"github.com/sashaaro/url-shortener/proto"
	"golang.org/x/crypto/acme/autocert"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
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

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	shrtenerService := domain.NewShortenerService(urlRepo)

	srv := http.Server{
		Addr:    internal.Config.ServerAddress,
		Handler: handlers.CreateServeMux(shrtenerService, logger, pool),
	}

	signalClosed := make(chan struct{})

	go func() {
		<-sigint
		log.Println("graceful shutdown...")

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}

		if pool != nil {
			log.Println("shutting down pool")
			pool.Close()
		}

		close(signalClosed)
	}()

	grpcTcpListen, err := net.Listen("tcp", ":"+strconv.Itoa(internal.Config.GrpcPort))
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	proto.RegisterURLShortenerServer(grpcServer, shortenerGrpc.NewGrpcService(shrtenerService, adapters.GenBase64ShortURLToken))

	go func() {
		log.Printf("Listen grpc")
		if err := grpcServer.Serve(grpcTcpListen); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		log.Printf("Listen http")
		if internal.Config.EnableHTTPS {
			err = srv.Serve(autocert.NewListener(internal.Config.ServerAddress, "url-shortener.ru", "www.url-shortener.ru"))
		} else {
			err = srv.ListenAndServe()
		}
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	<-signalClosed
}
