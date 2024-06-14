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

func Test_UnbindToNode(t *testing.T) {
	t.Run("returns an success on query, response for requested network address is received", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock, memory.New())
		zstack.sem = semaphore.NewWeighted(8)
		defer unpiMock.Stop()

		call := unpiMock.On(SREQ, ZDO, ZdoUnbindReqReplyID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   ZDO,
			CommandID:   ZdoUnbindReqReplyID,
			Payload:     []byte{0x00},
		})

		go func() {
			time.Sleep(10 * time.Millisecond)
			unpiMock.InjectOutgoing(Frame{
				MessageType: AREQ,
				Subsystem:   ZDO,
				CommandID:   ZdoUnbindRspID,
				Payload:     []byte{0x00, 0x40, 0x00},
			})
		}()

		zstack.nodeTable.addOrUpdate(zigbee.IEEEAddress(1), zigbee.NetworkAddress(0x4000))

		err := zstack.UnbindNodeFromController(ctx, zigbee.IEEEAddress(1), 2, 4, 5)
		assert.NoError(t, err)

		UnbindReq := ZdoUnbindReq{}
		bytecodec.Unmarshal(call.CapturedCalls[0].Frame.Payload, &UnbindReq)

		assert.Equal(t, zigbee.NetworkAddress(0x4000), UnbindReq.TargetAddress)
		assert.Equal(t, zigbee.IEEEAddress(1), UnbindReq.SourceAddress)
		assert.Equal(t, zigbee.Endpoint(2), UnbindReq.SourceEndpoint)
		assert.Equal(t, uint64(0), UnbindReq.DestinationAddress)
		assert.Equal(t, zigbee.Endpoint(4), UnbindReq.DestinationEndpoint)
		assert.Equal(t, zigbee.ClusterID(0x5), UnbindReq.ClusterID)
		assert.Equal(t, uint8(0x02), UnbindReq.DestinationAddressMode)

		unpiMock.AssertCalls(t)
	})
}

func Test_UnbindMessages(t *testing.T) {
	t.Run("verify ZdoUnbindReq marshals", func(t *testing.T) {
		req := ZdoUnbindReq{
			TargetAddress:          0x2021,
			SourceAddress:          zigbee.IEEEAddress(0x8899aabbccddeeff),
			SourceEndpoint:         0x01,
			ClusterID:              0xcafe,
			DestinationAddressMode: 0x01,
			DestinationAddress:     0x3ffe,
			DestinationEndpoint:    0x02,
		}

		data, err := bytecodec.Marshal(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x21, 0x20, 0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88, 0x01, 0xfe, 0xca, 0x01, 0xfe, 0x3f, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02}, data)
	})

	t.Run("verify ZdoUnbindReqReply marshals", func(t *testing.T) {
		req := ZdoUnbindReqReply{
			Status: 1,
		}

		data, err := bytecodec.Marshal(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x01}, data)
	})

	t.Run("ZdoUnbindReqReply returns true if success", func(t *testing.T) {
		g := ZdoUnbindReqReply{Status: ZSuccess}
		assert.True(t, g.WasSuccessful())
	})

	t.Run("ZdoUnbindReqReply returns false if not success", func(t *testing.T) {
		g := ZdoUnbindReqReply{Status: ZFailure}
		assert.False(t, g.WasSuccessful())
	})

	t.Run("verify ZdoUnbindRsp marshals", func(t *testing.T) {
		req := ZdoUnbindRsp{
			SourceAddress: zigbee.NetworkAddress(0x2000),
			Status:        1,
		}

		data, err := bytecodec.Marshal(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x00, 0x20, 0x01}, data)
	})

	t.Run("ZdoUnbindRsp returns true if success", func(t *testing.T) {
		g := ZdoUnbindRsp{Status: ZSuccess}
		assert.True(t, g.WasSuccessful())
	})

	t.Run("ZdoUnbindRsp returns false if not success", func(t *testing.T) {
		g := ZdoUnbindRsp{Status: ZFailure}
		assert.False(t, g.WasSuccessful())
	})
}
