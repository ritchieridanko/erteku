package interceptors

import (
	"context"

	"github.com/ritchieridanko/erteku/services/auth/internal/infra/logger"
	"google.golang.org/grpc"
)

func RecoveryInterceptor(l *logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		defer func() {
			if r := recover(); r != nil {
				l.Error(
					ctx,
					"PANIC RECOVERED",
					logger.NewField("method", info.FullMethod),
					logger.NewField("panic", r),
				)
			}
		}()

		return handler(ctx, req)
	}
}
