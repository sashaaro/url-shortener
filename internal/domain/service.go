// Package domain - сервис
package domain

import (
	"context"
	"github.com/google/uuid"
	"net/url"
)

// StatsResponse - dto
type StatsResponse struct {
	Urls  int64 `json:"urls"`
	Users int64 `json:"users"`
}

// ShortenerService - сервис
type ShortenerService struct {
	urlRepo          URLRepository
	genShortURLToken GenShortURLToken
}

// NewShortenerService конструктор
func NewShortenerService(urlRepo URLRepository, genShortURLToken GenShortURLToken) *ShortenerService {
	return &ShortenerService{urlRepo: urlRepo, genShortURLToken: genShortURLToken}
}

// GetOriginLink получение
func (r *ShortenerService) GetOriginLink(ctx context.Context, hashkey string) (*url.URL, error) {
	originURL, err := r.urlRepo.GetByHash(ctx, hashkey)
	return originURL, err
}

// BatchAdd создание
func (r *ShortenerService) BatchAdd(ctx context.Context, batch []BatchItem, userID uuid.UUID) ([]BatchItem, error) {
	for i, item := range batch {
		item.HashKey = r.genShortURLToken()
		batch[i] = item
	}
	return batch, r.urlRepo.BatchAdd(ctx, batch, userID)
}

// DeleteByUser удаление
func (r *ShortenerService) DeleteByUser(ctx context.Context, keys []HashKey, userID uuid.UUID) (bool, error) {
	return r.urlRepo.DeleteByUser(ctx, keys, userID)
}

// GetByUser получение
func (r *ShortenerService) GetByUser(ctx context.Context, userID uuid.UUID) ([]URLEntry, error) {
	return r.urlRepo.GetByUser(ctx, userID)
}

// CreateShort создание
func (r *ShortenerService) CreateShort(ctx context.Context, u url.URL, userID uuid.UUID) (HashKey, error) {
	key := r.genShortURLToken()

	return key, r.urlRepo.Add(ctx, key, u, userID)
}

// Stats статистика
func (r *ShortenerService) Stats(ctx context.Context) (*StatsResponse, error) {
	resp := &StatsResponse{}
	var err error
	resp.Urls, err = r.urlRepo.CountUrls(ctx)
	if err != nil {
		return nil, err
	}

	resp.Users, err = r.urlRepo.CountUsers(ctx)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
