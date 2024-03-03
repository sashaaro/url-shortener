package infra

import (
	"github.com/sashaaro/url-shortener/internal/domain"
	"net/url"
	"sync"
)

var _ domain.URLRepository = &memURLRepository{}

type memURLRepository struct {
	mx       sync.Mutex
	urlStore map[domain.URLKey]url.URL
}

func NewMemURLRepository() domain.URLRepository {
	return &memURLRepository{
		urlStore: map[domain.URLKey]url.URL{},
	}
}

func (m *memURLRepository) Add(key domain.URLKey, u url.URL) {
	m.mx.Lock()
	defer m.mx.Unlock()
	m.urlStore[key] = u
}

func (m *memURLRepository) GetByHash(key domain.URLKey) (url.URL, bool) {
	m.mx.Lock()
	defer m.mx.Unlock()
	u, ok := m.urlStore[key]
	return u, ok
}
