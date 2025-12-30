package caches

import (
	"context"
	"fmt"

	"github.com/ritchieridanko/erteku/services/auth/internal/constants"
	"github.com/ritchieridanko/erteku/services/auth/internal/infra/cache"
	"github.com/ritchieridanko/erteku/services/auth/internal/utils/ce"
)

type AuthCache interface {
	IsEmailReserved(ctx context.Context, email string) (reserved bool, err *ce.Error)
}

type authCache struct {
	cache *cache.Cache
}

func NewAuthCache(c *cache.Cache) AuthCache {
	return &authCache{cache: c}
}

func (c *authCache) IsEmailReserved(ctx context.Context, email string) (bool, *ce.Error) {
	key := fmt.Sprintf("%s:%s", constants.CachePrefixEmailReservation, email)

	exists, err := c.cache.Exists(ctx, key)
	if err != nil {
		return false, ce.NewError(ce.CodeCacheQueryExec, ce.MsgInternalServer, err)
	}

	return exists, nil
}
