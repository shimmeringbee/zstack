package zstack

import (
	"context"
	"github.com/shimmeringbee/bytecodec"
	. "github.com/shimmeringbee/unpi"
	unpiTest "github.com/shimmeringbee/unpi/testing"
	"github.com/shimmeringbee/zigbee"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_QueryNodeEndpoints(t *testing.T) {
	t.Run("returns an success on query, response for requested network address is received", func(t *testing.T) {
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

		zstack.nodeTable.AddOrUpdate(zigbee.IEEEAddress(0x11223344556677), zigbee.NetworkAddress(0x4000))

		endpoints, err := zstack.QueryNodeEndpoints(ctx, zigbee.IEEEAddress(0x11223344556677))
		assert.NoError(t, err)
		assert.Equal(t, []zigbee.Endpoint{0x01, 0x02, 0x03}, endpoints)

		unpiMock.AssertCalls(t)
	})
}

func Test_ActiveEndpointMessages(t *testing.T) {
	t.Run("verify ZdoActiveEpReq marshals", func(t *testing.T) {
		req := ZdoActiveEpReq{
			DestinationAddress: zigbee.NetworkAddress(0x2000),
			OfInterestAddress:  zigbee.NetworkAddress(0x4000),
		}

		data, err := bytecodec.Marshall(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x00, 0x20, 0x00, 0x40}, data)
	})

	t.Run("verify ZdoActiveEpReqReply marshals", func(t *testing.T) {
		req := ZdoActiveEpReqReply{
			Status: 1,
		}

		data, err := bytecodec.Marshall(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x01}, data)
	})

	t.Run("generic ZdoActiveEpReqReply returns true if success", func(t *testing.T) {
		g := ZdoActiveEpReqReply{Status: ZSuccess}
		assert.True(t, g.WasSuccessful())
	})

	t.Run("generic ZdoActiveEpReqReply returns false if not success", func(t *testing.T) {
		g := ZdoActiveEpReqReply{Status: ZFailure}
		assert.False(t, g.WasSuccessful())
	})

	t.Run("verify ZdoActiveEpRsp marshals", func(t *testing.T) {
		req := ZdoActiveEpRsp{
			SourceAddress:     zigbee.NetworkAddress(0x2000),
			Status:            1,
			OfInterestAddress: zigbee.NetworkAddress(0x4000),
			ActiveEndpoints:   []zigbee.Endpoint{0x01, 0x02, 0x03},
		}

		data, err := bytecodec.Marshall(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x00, 0x20, 0x01, 0x00, 0x40, 0x03, 0x01, 0x02, 0x03}, data)
	})

	t.Run("generic ZdoActiveEpRsp returns true if success", func(t *testing.T) {
		g := ZdoActiveEpRsp{Status: ZSuccess}
		assert.True(t, g.WasSuccessful())
	})

	t.Run("generic ZdoActiveEpRsp returns false if not success", func(t *testing.T) {
		g := ZdoActiveEpRsp{Status: ZFailure}
		assert.False(t, g.WasSuccessful())
	})
}
