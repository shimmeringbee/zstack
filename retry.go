package zstack

import (
	"context"
	"time"
)

func Retry(parent context.Context, duration time.Duration, attempts int, f func(ctx context.Context) error) (err error) {
	for i := 0; i < attempts; i++ {
		ctx, cancel := context.WithTimeout(parent, duration)

		if err = f(ctx); err == nil {
			cancel()
			return nil
		}

		cancel()
	}

	return
}
