package infra

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/sashaaro/url-shortener/internal"
	_ "github.com/sashaaro/url-shortener/migrations"
	"log"
)

func CreatePgxConn() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), internal.Config.DatabaseDSN)
	if err != nil {
		log.Fatal("can't connect to database", err)
	}

	db := stdlib.OpenDB(*conn.Config())

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal("can't set dialect: ", err)
	}

	if err := goose.Up(db, "./"); err != nil {
		log.Fatal("can't run migrations: ", err)
	}
	return conn
}
