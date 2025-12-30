package caches

import (
	"context"
	"fmt"

	"github.com/ritchieridanko/erteku/services/auth/configs"
	"github.com/ritchieridanko/erteku/services/auth/internal/constants"
	"github.com/ritchieridanko/erteku/services/auth/internal/infra/cache"
	"github.com/ritchieridanko/erteku/services/auth/internal/infra/logger"
	"github.com/ritchieridanko/erteku/services/auth/internal/utils/ce"
)

type TokenCache interface {
	StoreVerification(ctx context.Context, authID int64, token string) (err *ce.Error)
}

type tokenCache struct {
	config *configs.Auth
	cache  *cache.Cache
}

func NewTokenCache(cfg *configs.Auth, c *cache.Cache) TokenCache {
	return &tokenCache{config: cfg, cache: c}
}

func (c *tokenCache) StoreVerification(ctx context.Context, authID int64, token string) *ce.Error {
	prefix := constants.CachePrefixEmailVerification
	authKey := fmt.Sprintf("%s:%d", prefix, authID)
	tokenKey := fmt.Sprintf("%s:%s", prefix, token)
	duration := int(c.config.Duration.Verification.Seconds())

	script := `
		local token = redis.call("GET", KEYS[1])
		if token then
			redis.call("DEL", KEYS[1])
			redis.call("DEL", KEYS[3] .. ":" .. token)
		end

		redis.call("SET", KEYS[1], ARGV[1], "EX", ARGV[3])
		redis.call("SET", KEYS[2], ARGV[2], "EX", ARGV[3])
		return 1
	`

	_, err := c.cache.Evaluate(
		ctx, "s:stover", script,
		[]string{authKey, tokenKey, prefix},
		token, authID, duration,
	)
	if err != nil {
		return ce.NewError(
			ce.CodeCacheScriptExec,
			ce.MsgInternalServer,
			err,
			logger.NewField("auth_id", authID),
		)
	}

	return nil
}
