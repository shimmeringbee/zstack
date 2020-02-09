package zstack

import (
	"context"
	"errors"
)

func AnyResponse(v interface{}) bool {
	return true
}

var ReplyDoesNotReportSuccess = errors.New("reply struct does not support Successor interface")
var NodeResponseWasNotSuccess = errors.New("response from node was not success")

func (z *ZStack) nodeRequest(ctx context.Context, request interface{}, reply interface{}, response interface{}, responseFilter func(interface{}) bool) (interface{}, error) {
	replySuccessor, replySupportsSuccessor := reply.(Successor)

	if !replySupportsSuccessor {
		return nil, ReplyDoesNotReportSuccess
	}

	ch := make(chan interface{})

	err, stop := z.subscriber.Subscribe(response, func(v interface{}) {
		if responseFilter(v) {
			select {
			case ch <- v:
			case <-ctx.Done():
			}
		}
	})
	defer stop()

	if err != nil {
		return nil, err
	}

	if err := z.requestResponder.RequestResponse(ctx, request, reply); err != nil {
		return nil, err
	}

	if !replySuccessor.WasSuccessful() {
		return nil, ErrorZFailure
	}

	select {
	case v := <-ch:
		responseSuccessor, responseSupportsSuccessor := v.(Successor)

		if responseSupportsSuccessor && !responseSuccessor.WasSuccessful() {
			return v, NodeResponseWasNotSuccess
		}
		return v, nil
	case <-ctx.Done():
		return nil, errors.New("context expired while waiting for response from node")
	}
}
