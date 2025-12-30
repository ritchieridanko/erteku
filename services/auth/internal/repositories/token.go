package repositories

import (
	"context"

	"github.com/ritchieridanko/erteku/services/auth/internal/repositories/caches"
	"github.com/ritchieridanko/erteku/services/auth/internal/utils/ce"
)

type TokenRepository interface {
	CreateVerification(ctx context.Context, authID int64, token string) (err *ce.Error)
}

type tokenRepository struct {
	cache caches.TokenCache
}

func NewTokenRepository(c caches.TokenCache) TokenRepository {
	return &tokenRepository{cache: c}
}

func (r *tokenRepository) CreateVerification(ctx context.Context, authID int64, token string) *ce.Error {
	return r.cache.StoreVerification(ctx, authID, token)
}
