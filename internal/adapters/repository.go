package adapters

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sashaaro/url-shortener/internal/domain"
	"go.uber.org/zap"
	"io"
	"log"
	"net/url"
	"os"
	"sync"
)

var _ domain.URLRepository = &memURLRepository{}

type memURLRepository struct {
	mx       sync.Mutex
	urlStore map[domain.HashKey]url.URL
}

func (m *memURLRepository) BatchAdd(batch []domain.BatchItem) {
	for _, item := range batch {
		m.Add(item.HashKey, item.URL)
	}
}

func NewMemURLRepository() domain.URLRepository {
	return &memURLRepository{
		urlStore: map[domain.HashKey]url.URL{},
	}
}

func (m *memURLRepository) Add(key domain.HashKey, u url.URL) {
	m.mx.Lock()
	defer m.mx.Unlock()
	m.urlStore[key] = u
}

func (m *memURLRepository) GetByHash(key domain.HashKey) (url.URL, bool) {
	m.mx.Lock()
	defer m.mx.Unlock()
	u, ok := m.urlStore[key]
	return u, ok
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
}

type FileURLRepository struct {
	file    *os.File
	encoder *json.Encoder
	wrapped domain.URLRepository
	logger  zap.SugaredLogger
}

func (f *FileURLRepository) BatchAdd(batch []domain.BatchItem) {
	for _, item := range batch {
		f.Add(item.HashKey, item.URL)
	}
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
		f.wrapped.Add(entry.ShortURL, *u)
	}
	return nil
}

func (f *FileURLRepository) Close() error {
	return f.file.Close()
}

func (f FileURLRepository) Add(key domain.HashKey, u url.URL) {
	err := f.encoder.Encode(fileEntry{
		ID:          uuid.New(),
		ShortURL:    key,
		OriginalURL: u.String(),
	})
	if err != nil {
		log.Fatal(err)
	}
	f.wrapped.Add(key, u)
}

func (f FileURLRepository) GetByHash(key domain.HashKey) (url.URL, bool) {
	return f.wrapped.GetByHash(key)
}
