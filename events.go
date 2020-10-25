package zstack

import (
	"context"
)

func (z *ZStack) sendEvent(event interface{}) {
	z.events <- event
}

func (z *ZStack) ReadEvent(ctx context.Context) (interface{}, error) {
	select {
	case event := <-z.events:
		return event, nil
	case <-ctx.Done():
		return nil, context.DeadlineExceeded
	}
}
