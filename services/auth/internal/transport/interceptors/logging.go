package interceptors

import (
	"context"
	"errors"
	"time"

	"github.com/ritchieridanko/erteku/services/auth/internal/infra/logger"
	"github.com/ritchieridanko/erteku/services/auth/internal/utils/ce"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func LoggingInterceptor(l *logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now().UTC()
		resp, err := handler(ctx, req)

		st, _ := status.FromError(err)
		fields := []logger.Field{
			logger.NewField("method", info.FullMethod),
			logger.NewField("status", st.Code().String()),
			logger.NewField("latency", time.Since(start).String()),
		}

		if err == nil {
			l.Info(ctx, "REQUEST OK", fields...)
			return resp, nil
		}

		var e *ce.Error
		if errors.As(err, &e) {
			fields = append(fields, e.Fields()...)
			fields = append(
				fields,
				logger.NewField("error_code", e.Code()),
				logger.NewField("error", e),
			)

			l.Error(ctx, "REQUEST ERROR", fields...)
			return nil, e.ToGRPCStatus()
		}

		// fallback
		fields = append(
			fields,
			logger.NewField("error_code", ce.CodeUnknown),
			logger.NewField("error", err),
		)

		l.Error(ctx, "UNKNOWN ERROR", fields...)
		return nil, status.Error(codes.Internal, ce.MsgInternalServer)
	}
}
