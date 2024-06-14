package zstack

import (
	"context"
	unpiTest "github.com/shimmeringbee/unpi/testing"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_ReadEvent(t *testing.T) {
	t.Run("errors if context times out", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock)
		defer unpiMock.Stop()
		defer unpiMock.AssertCalls(t)

		_, err := zstack.ReadEvent(ctx)
		assert.Error(t, err)
	})
}
