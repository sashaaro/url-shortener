package adapters

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sashaaro/url-shortener/internal/domain"
	"net/url"
)

var _ domain.URLRepository = &PgURLRepository{}

// PgURLRepository - хранение ссылок в postgres
type PgURLRepository struct {
	pool *pgxpool.Pool
}

// CountUrls количество ссылок
func (r *PgURLRepository) CountUrls(ctx context.Context) (int64, error) {
	row := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM urls")
	var count int64
	err := row.Scan(&count)
	return count, err
}

// CountUsers количество пользователей
func (r *PgURLRepository) CountUsers(ctx context.Context) (int64, error) {
	row := r.pool.QueryRow(ctx, "SELECT COUNT(DISTINCT user_id) FROM urls")
	var count int64
	err := row.Scan(&count)
	return count, err
}

// DeleteByUser -удаление
func (r *PgURLRepository) DeleteByUser(ctx context.Context, keys []domain.HashKey, userID uuid.UUID) (bool, error) {
	res, err := r.pool.Exec(ctx, "UPDATE urls SET is_deleted = true WHERE key = ANY($1)", keys)

	return res.RowsAffected() == int64(len(keys)), err
}

// GetByUser получение
func (r *PgURLRepository) GetByUser(ctx context.Context, userID uuid.UUID) ([]domain.URLEntry, error) {
	rows, err := r.pool.Query(ctx, "SELECT key, url FROM urls WHERE user_id = $1", userID.String())
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
		if err = rows.Scan(&key, &u); err != nil {
			return nil, err
		}
		urls = append(urls, domain.URLEntry{
			ShortURL:    CreatePublicURL(key),
			OriginalURL: u,
		})
	}
	return urls, nil
}

// BatchAdd - добавление нескольких ссылок
func (r *PgURLRepository) BatchAdd(ctx context.Context, batch []domain.BatchItem, userID uuid.UUID) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	// nolint:errcheck
	defer tx.Rollback(ctx)

	for _, item := range batch {
		_, err = tx.Exec(ctx, "INSERT INTO urls (key, url, user_id) VALUES ($1, $2, $3)", item.HashKey, item.URL.String(), userID.String())
		if err != nil {
			pgErr := &pgconn.PgError{}
			ok := errors.As(err, &pgErr)
			if ok && pgErr.Code == pgerrcode.UniqueViolation {
				var existKey string
				err = r.pool.QueryRow(ctx, "SELECT key FROM urls WHERE url = $1 LIMIT 1", item.URL.String()).Scan(&existKey)
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

// Add добавление ссылки
func (r *PgURLRepository) Add(ctx context.Context, key domain.HashKey, u url.URL, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, "INSERT INTO urls (key, url, user_id) VALUES ($1, $2, $3)", key, u.String(), userID.String())
	if err != nil {
		pgErr := &pgconn.PgError{}
		ok := errors.As(err, &pgErr)
		if ok && pgErr.Code == pgerrcode.UniqueViolation {
			var existKey string
			err = r.pool.QueryRow(ctx, "SELECT key FROM urls WHERE url = $1 LIMIT 1", u.String()).Scan(&existKey)
			if err != nil {
				return err
			}
			return &domain.ErrURLAlreadyExists{HashKey: existKey}
		}
	}

	return err
}

// GetByHash - получение ссылки по ключу
func (r *PgURLRepository) GetByHash(ctx context.Context, key domain.HashKey) (*url.URL, error) {
	var res string
	var isDeleted bool
	err := r.pool.QueryRow(ctx, "SELECT url, is_deleted FROM urls WHERE key = $1", key).Scan(&res, &isDeleted)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if isDeleted {
		return nil, domain.ErrURLDeleted
	}
	return url.Parse(res)
}

// NewPgURLRepository - конструктор
func NewPgURLRepository(pool *pgxpool.Pool) *PgURLRepository {
	repo := &PgURLRepository{
		pool: pool,
	}
	return repo
}
