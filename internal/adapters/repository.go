package adapters

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/url"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/sashaaro/url-shortener/internal/domain"
	"go.uber.org/zap"
)

var _ domain.URLRepository = &memURLRepository{
	urlStore: map[string]memEntry{},
}

type memEntry struct {
	url    url.URL
	hash   domain.HashKey
	userID uuid.UUID
}

type memURLRepository struct {
	mx       sync.Mutex
	urlStore map[domain.HashKey]memEntry
}

func (m *memURLRepository) GetByUser(ctx context.Context, userID uuid.UUID) ([]domain.UrlEntry, error) {
	l := make([]domain.UrlEntry, 0)
	m.mx.Lock()
	defer m.mx.Unlock()
	for _, v := range m.urlStore {
		if v.userID == userID {
			shortUrl, err := url.Parse(CreatePublicURL(v.hash))
			if err != nil {
				return nil, err
			}

			l = append(l, domain.UrlEntry{
				ShortUrl:    *shortUrl,
				OriginalUrl: v.url,
			})
		}
	}
	return l, nil
}

func (m *memURLRepository) BatchAdd(ctx context.Context, batch []domain.BatchItem, userID uuid.UUID) error {
	for _, item := range batch {
		err := m.Add(ctx, item.HashKey, item.URL, userID)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewMemURLRepository() domain.URLRepository {
	return &memURLRepository{
		urlStore: map[domain.HashKey]memEntry{},
	}
}

func (m *memURLRepository) Add(ctx context.Context, key domain.HashKey, u url.URL, userID uuid.UUID) error {
	m.mx.Lock()
	defer m.mx.Unlock()
	m.urlStore[key] = memEntry{
		url:    u,
		hash:   key,
		userID: userID,
	}
	return nil
}

func (m *memURLRepository) GetByHash(ctx context.Context, key domain.HashKey) (*url.URL, error) {
	m.mx.Lock()
	defer m.mx.Unlock()
	u, ok := m.urlStore[key]
	if ok {
		return &u.url, nil
	} else {
		return nil, nil
	}
}

var _ domain.URLRepository = &FileURLRepository{}

func NewFileURLRepository(
	filePath string,
	wrapped domain.URLRepository,
	logger zap.SugaredLogger,
) *FileURLRepository {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	r := &FileURLRepository{
		file:    file,
		wrapped: wrapped,
		encoder: json.NewEncoder(file),
		logger:  logger,
	}
	err = r.load()
	if err != nil {
		log.Fatal(err)
	}

	return r
}

type fileEntry struct {
	ID          uuid.UUID `json:"id"`
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
	UserID      uuid.UUID `json:"user_id"`
}

type FileURLRepository struct {
	file    *os.File
	encoder *json.Encoder
	wrapped domain.URLRepository
	logger  zap.SugaredLogger
}

func (f *FileURLRepository) GetByUser(ctx context.Context, userID uuid.UUID) ([]domain.UrlEntry, error) {
	return f.wrapped.GetByUser(ctx, userID)
}

func (f *FileURLRepository) BatchAdd(ctx context.Context, batch []domain.BatchItem, userID uuid.UUID) error {
	for _, item := range batch {
		err := f.Add(ctx, item.HashKey, item.URL, userID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *FileURLRepository) load() error {
	decoder := json.NewDecoder(f.file)
	var entry fileEntry
	for {
		if err := decoder.Decode(&entry); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		u, err := url.Parse(entry.OriginalURL)
		if err != nil {
			f.logger.Warn("invalid db url entry")
			continue
		}
		err = f.wrapped.Add(context.Background(), entry.ShortURL, *u, entry.UserID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *FileURLRepository) Close() error {
	return f.file.Close()
}

func (f FileURLRepository) Add(ctx context.Context, key domain.HashKey, u url.URL, userID uuid.UUID) error {
	err := f.wrapped.Add(ctx, key, u, userID)
	if err != nil {
		return err
	}
	err = f.encoder.Encode(fileEntry{
		ID:          uuid.New(),
		ShortURL:    key,
		OriginalURL: u.String(),
	})
	return err
}

func (f FileURLRepository) GetByHash(ctx context.Context, key domain.HashKey) (*url.URL, error) {
	return f.wrapped.GetByHash(ctx, key)
}
