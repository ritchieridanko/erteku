package cache

import (
	"context"
	"errors"
	"time"
)

func backoffWait(ctx context.Context, delay, attempt int) error {
	backoff := time.Duration(delay) * (1 << attempt)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(backoff):
		return nil
	}
}

func isRetryable(err error) bool {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}
	return true
}
