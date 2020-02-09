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

func (z *ZStack) nodeRequest(ctx context.Context, request interface{}, reply interface{}, response interface{}, responseFilter func(interface{}) bool) error {
	replySuccessor, replySupportsSuccessor := reply.(Successor)

	if !replySupportsSuccessor {
		return ReplyDoesNotReportSuccess
	}

	ch := make(chan bool)

	err, stop := z.subscriber.Subscribe(response, func(v interface{}) {
		if responseFilter(v) {
			select {
			case ch <- true:
			case <-ctx.Done():
			}
		}
	})
	defer stop()

	if err != nil {
		return err
	}

	if err := z.requestResponder.RequestResponse(ctx, request, reply); err != nil {
		return err
	}

	if !replySuccessor.WasSuccessful() {
		return ErrorZFailure
	}

	select {
	case <-ch:
		responseSuccessor, responseSupportsSuccessor := response.(Successor)

		if responseSupportsSuccessor && !responseSuccessor.WasSuccessful() {
			return NodeResponseWasNotSuccess
		}
		return nil
	case <-ctx.Done():
		return errors.New("context expired while waiting for response from node")
	}
}
