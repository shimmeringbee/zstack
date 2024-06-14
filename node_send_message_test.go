package zstack

import (
	"context"
	"github.com/shimmeringbee/bytecodec"
	"github.com/shimmeringbee/persistence/impl/memory"
	. "github.com/shimmeringbee/unpi"
	unpiTest "github.com/shimmeringbee/unpi/testing"
	"github.com/shimmeringbee/zigbee"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/semaphore"
	"testing"
	"time"
)

func Test_SendNodeMessage(t *testing.T) {
	t.Run("messages with ack wait", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		defer unpiMock.AssertCalls(t)
		zstack := New(unpiMock, memory.New())
		zstack.sem = semaphore.NewWeighted(8)
		defer unpiMock.Stop()

		zstack.nodeTable.addOrUpdate(zigbee.IEEEAddress(0x1122334455667788), zigbee.NetworkAddress(0x1000))

		c := unpiMock.On(SREQ, AF, AfDataRequestID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   AF,
			CommandID:   AfDataRequestReplyID,
			Payload:     []byte{0x00},
		})

		go func() {
			time.Sleep(10 * time.Millisecond)

			unpiMock.InjectOutgoing(Frame{
				MessageType: AREQ,
				Subsystem:   AF,
				CommandID:   AfDataConfirmID,
				Payload:     []byte{0x00, 0x04, 0x00},
			})
		}()

		appMessage := zigbee.ApplicationMessage{
			ClusterID:           0x2000,
			SourceEndpoint:      0x03,
			DestinationEndpoint: 0x04,
			Data:                []byte{0x0a, 0x0b},
		}

		err := zstack.SendApplicationMessageToNode(ctx, zigbee.IEEEAddress(0x1122334455667788), appMessage, true)
		assert.NoError(t, err)

		sentFrame := c.CapturedCalls[0].Frame

		assert.Equal(t, []byte{0x00, 0x10, 0x04, 0x03, 0x00, 0x20, 0x00, 0x10, 0x20, 0x02, 0x0a, 0x0b}, sentFrame.Payload)
	})

	t.Run("messages without ack just return", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		defer unpiMock.AssertCalls(t)
		zstack := New(unpiMock, memory.New())
		zstack.sem = semaphore.NewWeighted(8)
		defer unpiMock.Stop()

		zstack.nodeTable.addOrUpdate(zigbee.IEEEAddress(0x1122334455667788), zigbee.NetworkAddress(0x1000))

		c := unpiMock.On(SREQ, AF, AfDataRequestID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   AF,
			CommandID:   AfDataRequestReplyID,
			Payload:     []byte{0x00},
		})

		appMessage := zigbee.ApplicationMessage{
			ClusterID:           0x2000,
			SourceEndpoint:      0x03,
			DestinationEndpoint: 0x04,
			Data:                []byte{0x0a, 0x0b},
		}

		err := zstack.SendApplicationMessageToNode(ctx, zigbee.IEEEAddress(0x1122334455667788), appMessage, false)
		assert.NoError(t, err)

		sentFrame := c.CapturedCalls[0].Frame

		assert.Equal(t, []byte{0x00, 0x10, 0x04, 0x03, 0x00, 0x20, 0x00, 0x0, 0x20, 0x02, 0x0a, 0x0b}, sentFrame.Payload)
	})
}

func Test_SendMessages(t *testing.T) {
	t.Run("verify AfDataRequest marshals", func(t *testing.T) {
		req := AfDataRequest{
			DestinationAddress:  0x0102,
			DestinationEndpoint: 0x03,
			SourceEndpoint:      0x04,
			ClusterID:           0x0506,
			TransactionID:       0x07,
			Options: AfDataRequestOptions{
				EnableSecurity: true,
				DiscoveryRoute: true,
				ACKRequest:     true,
			},
			Radius: 0x09,
			Data:   []byte{0x0a, 0x0b},
		}

		data, err := bytecodec.Marshal(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x02, 0x01, 0x03, 0x04, 0x06, 0x05, 0x07, 0x70, 0x09, 0x02, 0x0a, 0x0b}, data)
	})

	t.Run("verify AfDataRequestReply marshals", func(t *testing.T) {
		req := AfDataRequestReply{
			Status: 1,
		}

		data, err := bytecodec.Marshal(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x01}, data)
	})

	t.Run("AfDataRequestReply returns true if success", func(t *testing.T) {
		g := AfDataRequestReply{Status: ZSuccess}
		assert.True(t, g.WasSuccessful())
	})

	t.Run("AfDataRequestReply returns false if not success", func(t *testing.T) {
		g := AfDataRequestReply{Status: ZFailure}
		assert.False(t, g.WasSuccessful())
	})

	t.Run("verify AfDataConfirm marshals", func(t *testing.T) {
		req := AfDataConfirm{
			Status:        0x01,
			Endpoint:      0x02,
			TransactionID: 0x03,
		}

		data, err := bytecodec.Marshal(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x01, 0x02, 0x03}, data)
	})

	t.Run("AfDataConfirm returns true if success", func(t *testing.T) {
		g := AfDataConfirm{Status: ZSuccess}
		assert.True(t, g.WasSuccessful())
	})

	t.Run("AfDataConfirm returns false if not success", func(t *testing.T) {
		g := AfDataConfirm{Status: ZFailure}
		assert.False(t, g.WasSuccessful())
	})
}
