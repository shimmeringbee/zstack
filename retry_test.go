package zstack

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_Retry(t *testing.T) {
	t.Run("retry passes through a nil error", func(t *testing.T) {
		err := Retry(context.Background(), 50 * time.Millisecond, 3, func(ctx context.Context) error {
			return nil
		})

		assert.NoError(t, err)
	})

	t.Run("retry retries the number of attempts", func(t *testing.T) {
		attempts := 0
		err := Retry(context.Background(), 50 * time.Millisecond, 3, func(ctx context.Context) error {
			attempts += 1
			return errors.New("general failure")
		})

		assert.Error(t, err)
		assert.Equal(t, 3, attempts)
	})
}