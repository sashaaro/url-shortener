package domain

import (
	"context"
	"github.com/google/uuid"
	"net/url"
)

type StatsResponse struct {
	Urls  int64 `json:"urls"`
	Users int64 `json:"users"`
}

type ShortenerService struct {
	urlRepo URLRepository
}

func NewShortenerService(urlRepo URLRepository) *ShortenerService {
	return &ShortenerService{urlRepo: urlRepo}
}

func (r *ShortenerService) GetOriginLink(ctx context.Context, hashkey string) (*url.URL, error) {
	originURL, err := r.urlRepo.GetByHash(ctx, hashkey)
	return originURL, err
}

func (r *ShortenerService) BatchAdd(ctx context.Context, batch []BatchItem, userID uuid.UUID) error {
	return r.urlRepo.BatchAdd(ctx, batch, userID)
}

func (r *ShortenerService) DeleteByUser(ctx context.Context, keys []HashKey, userID uuid.UUID) (bool, error) {
	return r.urlRepo.DeleteByUser(ctx, keys, userID)
}

func (r *ShortenerService) GetByUser(ctx context.Context, userID uuid.UUID) ([]URLEntry, error) {
	return r.urlRepo.GetByUser(ctx, userID)
}

func (r *ShortenerService) CreateShort(ctx context.Context, key HashKey, u url.URL, userID uuid.UUID) error {
	return r.urlRepo.Add(ctx, key, u, userID)
}

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