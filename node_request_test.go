package zstack

import (
	"context"
	. "github.com/shimmeringbee/unpi"
	unpiTest "github.com/shimmeringbee/unpi/testing"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)


func Test_NodeRequest(t *testing.T)  {
	t.Run("returns an error if the response type does not support Successor", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock)
		defer unpiMock.Stop()

		type NotSuccessful struct {}

		err := zstack.nodeRequest(ctx, &NotSuccessful{}, &NotSuccessful{}, &NotSuccessful{}, AnyResponse)

		assert.Error(t, err)
		assert.Equal(t, ReplyDoesNotReportSuccess, err)
	})

	t.Run("returns an error if the request results in a failed response", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock)
		defer unpiMock.Stop()

		unpiMock.On(SREQ, ZDO, ZdoActiveEpReqID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   ZDO,
			CommandID:   ZdoActiveEpReqReplyID,
			Payload:     []byte{0x01},
		})

		err := zstack.nodeRequest(ctx, &ZdoActiveEpReq{}, &ZdoActiveEpReqReply{}, &ZdoActiveEpRsp{}, AnyResponse)
		assert.Error(t, err)
		assert.Equal(t, ErrorZFailure, err)

		unpiMock.AssertCalls(t)
	})

	t.Run("returns a success, when request was successfully replied to, and response has been received", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock)
		defer unpiMock.Stop()

		unpiMock.On(SREQ, ZDO, ZdoActiveEpReqID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   ZDO,
			CommandID:   ZdoActiveEpReqReplyID,
			Payload:     []byte{0x00},
		})

		go func() {
			time.Sleep(10 * time.Millisecond)
			unpiMock.InjectOutgoing(Frame{
				MessageType: AREQ,
				Subsystem:   ZDO,
				CommandID:   ZdoActiveEpRspID,
				Payload:     []byte{0x00, 0x20, 0x00, 0x00, 0x40, 0x03, 0x01, 0x02, 0x03},
			})
		}()

		response := &ZdoActiveEpRsp{}
		err := zstack.nodeRequest(ctx, &ZdoActiveEpReq{DestinationAddress:0x4000, OfInterestAddress:0x4000}, &ZdoActiveEpReqReply{}, response, AnyResponse)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0x01, 0x02, 0x03}, response.ActiveEndpoints)

		unpiMock.AssertCalls(t)
	})

	t.Run("returns a success, when request was successfully replied to with a response which is unwanted, then wanted", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock)
		defer unpiMock.Stop()

		unpiMock.On(SREQ, ZDO, ZdoActiveEpReqID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   ZDO,
			CommandID:   ZdoActiveEpReqReplyID,
			Payload:     []byte{0x00},
		})

		go func() {
			time.Sleep(10 * time.Millisecond)
			unpiMock.InjectOutgoing(Frame{
				MessageType: AREQ,
				Subsystem:   ZDO,
				CommandID:   ZdoActiveEpRspID,
				Payload:     []byte{0x00, 0x20, 0x00, 0x00, 0x20, 0x03, 0x01, 0x02, 0x03},
			})
			time.Sleep(10 * time.Millisecond)
			unpiMock.InjectOutgoing(Frame{
				MessageType: AREQ,
				Subsystem:   ZDO,
				CommandID:   ZdoActiveEpRspID,
				Payload:     []byte{0x00, 0x20, 0x00, 0x00, 0x40, 0x01, 0x02},
			})
		}()

		response := &ZdoActiveEpRsp{}
		err := zstack.nodeRequest(ctx, &ZdoActiveEpReq{DestinationAddress:0x4000, OfInterestAddress:0x4000}, &ZdoActiveEpReqReply{}, response, func(i interface{}) bool {
			response := i.(*ZdoActiveEpRsp)
			return response.OfInterestAddress == 0x4000
		})
		assert.NoError(t, err)
		assert.Equal(t, []byte{0x02}, response.ActiveEndpoints)

		unpiMock.AssertCalls(t)
	})

	t.Run("returns a success, when request was successfully replied to with a response which is wanted, then unwanted", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock)
		defer unpiMock.Stop()

		unpiMock.On(SREQ, ZDO, ZdoActiveEpReqID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   ZDO,
			CommandID:   ZdoActiveEpReqReplyID,
			Payload:     []byte{0x00},
		})

		go func() {
			time.Sleep(10 * time.Millisecond)
			unpiMock.InjectOutgoing(Frame{
				MessageType: AREQ,
				Subsystem:   ZDO,
				CommandID:   ZdoActiveEpRspID,
				Payload:     []byte{0x00, 0x20, 0x00, 0x00, 0x40, 0x01, 0x02},
			})
			unpiMock.InjectOutgoing(Frame{
				MessageType: AREQ,
				Subsystem:   ZDO,
				CommandID:   ZdoActiveEpRspID,
				Payload:     []byte{0x00, 0x20, 0x00, 0x00, 0x20, 0x03, 0x01, 0x02, 0x03},
			})
		}()

		response := &ZdoActiveEpRsp{}
		err := zstack.nodeRequest(ctx, &ZdoActiveEpReq{DestinationAddress:0x4000, OfInterestAddress:0x4000}, &ZdoActiveEpReqReply{}, response, func(i interface{}) bool {
			response := i.(*ZdoActiveEpRsp)
			return response.OfInterestAddress == 0x4000
		})
		assert.NoError(t, err)
		assert.Equal(t, []byte{0x02}, response.ActiveEndpoints)

		unpiMock.AssertCalls(t)
	})

	t.Run("returns an error, when request was successfully replied to, but response supports WasSuccessful and is a failure", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock)
		defer unpiMock.Stop()

		unpiMock.On(SREQ, ZDO, ZdoActiveEpReqID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   ZDO,
			CommandID:   ZdoActiveEpReqReplyID,
			Payload:     []byte{0x00},
		})

		go func() {
			time.Sleep(10 * time.Millisecond)
			unpiMock.InjectOutgoing(Frame{
				MessageType: AREQ,
				Subsystem:   ZDO,
				CommandID:   ZdoActiveEpRspID,
				Payload:     []byte{0x00, 0x20, 0x01, 0x00, 0x40, 0x03, 0x01, 0x02, 0x03},
			})
		}()

		response := &ZdoActiveEpRsp{}
		err := zstack.nodeRequest(ctx, &ZdoActiveEpReq{DestinationAddress:0x4000, OfInterestAddress:0x4000}, &ZdoActiveEpReqReply{}, response, AnyResponse)
		assert.Error(t, err)
		assert.Equal(t, NodeResponseWasNotSuccess, err)
		assert.Equal(t, []byte{0x01, 0x02, 0x03}, response.ActiveEndpoints)

		unpiMock.AssertCalls(t)
	})
}