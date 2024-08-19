package adapters

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/sashaaro/url-shortener/internal/utils"
	"net/http"
)

type userContext struct{}

func userIDFromReq(req *http.Request) (uuid.UUID, error) {
	return userIDFromCtx(req.Context())
}

// получение пользователя из http запроса
func MustUserIDFromReq(req *http.Request) uuid.UUID {
	return utils.Must(userIDFromReq(req))
}

func userIDFromCtx(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(&userContext{}).(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("user id not found")
	}
	return userID, nil
}

// получение пользователя из контекста
func UserIDToCxt(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, &userContext{}, userID)
}
