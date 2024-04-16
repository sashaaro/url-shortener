package adapters

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/sashaaro/url-shortener/internal/domain"
	"net/url"
)

var _ domain.URLRepository = &PgURLRepository{}

type PgURLRepository struct {
	conn *pgx.Conn
}

func (r *PgURLRepository) BatchAdd(batch []domain.BatchItem) {
	tx, err := r.conn.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		panic(err)
	}
	// nolint:errcheck
	defer tx.Rollback(context.Background())

	for _, item := range batch {
		_, err := tx.Exec(context.Background(), "INSERT INTO urls (key, url) VALUES ($1, $2)", item.HashKey, item.URL.String())
		if err != nil {
			panic(err)
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		panic(err)
	}
}

func (r *PgURLRepository) Add(key domain.HashKey, u url.URL) {
	_, err := r.conn.Exec(context.Background(), "INSERT INTO urls (key, url) VALUES ($1, $2)", key, u.String())
	if err != nil {
		panic(err)
	}
}

func (r *PgURLRepository) GetByHash(key domain.HashKey) (url.URL, bool) {
	var res string
	err := r.conn.QueryRow(context.Background(), "SELECT url FROM urls WHERE key = $1", key).Scan(&res)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return url.URL{}, false
		}
		panic(err)
	}
	u, err := url.Parse(res)
	if err != nil {
		panic(err)
	}
	return *u, true
}

func NewPgURLRepository(conn *pgx.Conn) *PgURLRepository {
	repo := &PgURLRepository{
		conn: conn,
	}
	_, err := conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS urls (key text PRIMARY KEY, url text)")
	if err != nil {
		panic(err)
	}
	return repo
}
