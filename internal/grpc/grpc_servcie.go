// Package grpc - grpc сервис
package grpc

import (
	"context"
	"github.com/google/uuid"
	"github.com/sashaaro/url-shortener/internal/adapters"
	"github.com/sashaaro/url-shortener/internal/domain"
	"github.com/sashaaro/url-shortener/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/url"
	"strings"
)

// GrpcService - сервис
type GrpcService struct {
	proto.UnimplementedURLShortenerServer
	genShortURLToken domain.GenShortURLToken
	service          *domain.ShortenerService
}

// NewGrpcService конструктор
func NewGrpcService(service *domain.ShortenerService, genShortURLToken domain.GenShortURLToken) *GrpcService {
	return &GrpcService{service: service, genShortURLToken: genShortURLToken}
}

// Shorten создание
func (s *GrpcService) Shorten(ctx context.Context, req *proto.ShortenRequest) (*proto.ShortenResponse, error) {
	key := s.genShortURLToken()
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	originURL, err := url.Parse(req.Url)
	if err != nil || !strings.HasPrefix(originURL.Scheme, "http") {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.service.CreateShort(ctx, key, *originURL, userID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &proto.ShortenResponse{Result: adapters.CreatePublicURL(key)}, nil
}

// GetOriginLink получение
func (s *GrpcService) GetOriginLink(ctx context.Context, req *proto.GetOriginLinkRequest) (*proto.GetOriginLinkResponse, error) {
	originLink, err := s.service.GetOriginLink(ctx, req.Hash)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &proto.GetOriginLinkResponse{OriginalUrl: originLink.String()}, nil
}

// GetUserUrls получение
func (s *GrpcService) GetUserUrls(ctx context.Context, req *proto.GetUserUrlsRequest) (*proto.GetUserUrlsResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	list, err := s.service.GetByUser(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	res := &proto.GetUserUrlsResponse{
		Urls: make([]string, 0, len(list)),
	}

	for _, i := range list {
		res.Urls = append(res.Urls, i.OriginalURL)
	}

	return res, nil
}

// DeleteUrls удаление
func (s *GrpcService) DeleteUrls(ctx context.Context, req *proto.DeleteUrlsRequest) (*proto.DeleteUrlsResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	keys := make([]domain.HashKey, 0, len(req.Keys))
	for _, key := range req.Keys {
		keys = append(keys, domain.HashKey(key))
	}

	_, err = s.service.DeleteByUser(ctx, keys, userID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &proto.DeleteUrlsResponse{}, nil
}

// GetStats статистика
func (s *GrpcService) GetStats(ctx context.Context, req *proto.StatsRequest) (*proto.StatsResponse, error) {
	res, err := s.service.Stats(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &proto.StatsResponse{
		Urls:  res.Urls,
		Users: res.Users,
	}, nil
}

// Ping пинг
func (s *GrpcService) Ping(ctx context.Context, req *proto.PingRequest) (*proto.PongResponse, error) {
	return &proto.PongResponse{}, nil
}
