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

func Test_BindToNode(t *testing.T) {
	t.Run("returns an success on query, response for requested network address is received", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock)
		defer unpiMock.Stop()

		call := unpiMock.On(SREQ, ZDO, ZdoBindReqReplyID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   ZDO,
			CommandID:   ZdoBindReqReplyID,
			Payload:     []byte{0x00},
		})

		go func() {
			time.Sleep(10 * time.Millisecond)
			unpiMock.InjectOutgoing(Frame{
				MessageType: AREQ,
				Subsystem:   ZDO,
				CommandID:   ZdoBindRspID,
				Payload:     []byte{0x00, 0x40, 0x00},
			})
		}()

		zstack.nodeTable.AddOrUpdate(zigbee.IEEEAddress(1), zigbee.NetworkAddress(0x4000))

		err := zstack.BindNodeToController(ctx, zigbee.IEEEAddress(1), 2, 4, 5)
		assert.NoError(t, err)

		bindReq := ZdoBindReq{}
		bytecodec.Unmarshal(call.CapturedCalls[0].Frame.Payload, &bindReq)

		assert.Equal(t, zigbee.NetworkAddress(0x4000), bindReq.TargetAddress)
		assert.Equal(t, zigbee.IEEEAddress(1), bindReq.SourceAddress)
		assert.Equal(t, zigbee.Endpoint(2), bindReq.SourceEndpoint)
		assert.Equal(t, uint64(0), bindReq.DestinationAddress)
		assert.Equal(t, zigbee.Endpoint(4), bindReq.DestinationEndpoint)
		assert.Equal(t, zigbee.ClusterID(0x5), bindReq.ClusterID)
		assert.Equal(t, uint8(0x02), bindReq.DestinationAddressMode)

		unpiMock.AssertCalls(t)
	})
}

func Test_BindMessages(t *testing.T) {
	t.Run("verify ZdoBindReq marshals", func(t *testing.T) {
		req := ZdoBindReq{
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

	t.Run("verify ZdoBindReqReply marshals", func(t *testing.T) {
		req := ZdoBindReqReply{
			Status: 1,
		}

		data, err := bytecodec.Marshal(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x01}, data)
	})

	t.Run("ZdoBindReqReply returns true if success", func(t *testing.T) {
		g := ZdoBindReqReply{Status: ZSuccess}
		assert.True(t, g.WasSuccessful())
	})

	t.Run("ZdoBindReqReply returns false if not success", func(t *testing.T) {
		g := ZdoBindReqReply{Status: ZFailure}
		assert.False(t, g.WasSuccessful())
	})

	t.Run("verify ZdoBindRsp marshals", func(t *testing.T) {
		req := ZdoBindRsp{
			SourceAddress: zigbee.NetworkAddress(0x2000),
			Status:        1,
		}

		data, err := bytecodec.Marshal(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x00, 0x20, 0x01}, data)
	})

	t.Run("ZdoBindRsp returns true if success", func(t *testing.T) {
		g := ZdoBindRsp{Status: ZSuccess}
		assert.True(t, g.WasSuccessful())
	})

	t.Run("ZdoBindRsp returns false if not success", func(t *testing.T) {
		g := ZdoBindRsp{Status: ZFailure}
		assert.False(t, g.WasSuccessful())
	})
}
