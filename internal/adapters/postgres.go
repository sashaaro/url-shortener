package adapters

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sashaaro/url-shortener/internal/domain"
	"net/url"
)

var _ domain.URLRepository = &PgURLRepository{}

type PgURLRepository struct {
	conn *pgx.Conn
}

func (r *PgURLRepository) GetByUser(ctx context.Context, userID uuid.UUID) ([]domain.URLEntry, error) {
	rows, err := r.conn.Query(ctx, "SELECT key, url FROM urls WHERE user_id = $1", userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	urls := []domain.URLEntry{}
	var key, u string
	for rows.Next() {
		if err := rows.Scan(&key, &u); err != nil {
			return nil, err
		}
		shortURL, err := url.Parse(CreatePublicURL(key))
		if err != nil {
			return nil, err
		}
		originalURL, err := url.Parse(u)
		if err != nil {
			return nil, err
		}
		urls = append(urls, domain.URLEntry{
			ShortURL:    *shortURL,
			OriginalURL: *originalURL,
		})
	}
	return urls, nil
}

func (r *PgURLRepository) BatchAdd(ctx context.Context, batch []domain.BatchItem, userID uuid.UUID) error {
	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	// nolint:errcheck
	defer tx.Rollback(ctx)

	for _, item := range batch {
		_, err := tx.Exec(ctx, "INSERT INTO urls (key, url, user_id) VALUES ($1, $2, $3)", item.HashKey, item.URL.String(), userID)
		if err != nil {
			pgErr := &pgconn.PgError{}
			ok := errors.As(err, &pgErr)
			if ok && pgErr.Code == pgerrcode.UniqueViolation {
				var existKey string
				err := r.conn.QueryRow(ctx, "SELECT key FROM urls WHERE url = $1 LIMIT 1", item.URL.String()).Scan(&existKey)
				if err != nil {
					return err
				}
				return &domain.ErrURLAlreadyExists{HashKey: existKey}
			}
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *PgURLRepository) Add(ctx context.Context, key domain.HashKey, u url.URL, userID uuid.UUID) error {
	_, err := r.conn.Exec(ctx, "INSERT INTO urls (key, url, user_id) VALUES ($1, $2, $3)", key, u.String(), userID)
	if err != nil {
		pgErr := &pgconn.PgError{}
		ok := errors.As(err, &pgErr)
		if ok && pgErr.Code == pgerrcode.UniqueViolation {
			var existKey string
			err := r.conn.QueryRow(ctx, "SELECT key FROM urls WHERE url = $1 LIMIT 1", u.String()).Scan(&existKey)
			if err != nil {
				return err
			}
			return &domain.ErrURLAlreadyExists{HashKey: existKey}
		}
	}

	return err
}

func (r *PgURLRepository) GetByHash(ctx context.Context, key domain.HashKey) (*url.URL, error) {
	var res string
	err := r.conn.QueryRow(ctx, "SELECT url FROM urls WHERE key = $1", key).Scan(&res)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return url.Parse(res)
}

func NewPgURLRepository(conn *pgx.Conn) *PgURLRepository {
	repo := &PgURLRepository{
		conn: conn,
	}
	return repo
}
