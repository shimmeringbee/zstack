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

func TestZStack_RemoveNode(t *testing.T) {
	t.Run("returns an success on query, response for requested network address is received", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock, NewNodeTable())
		defer unpiMock.Stop()

		call := unpiMock.On(SREQ, ZDO, ZdoMgmtLeaveReqReplyID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   ZDO,
			CommandID:   ZdoMgmtLeaveReqReplyID,
			Payload:     []byte{0x00},
		})

		go func() {
			time.Sleep(10 * time.Millisecond)
			unpiMock.InjectOutgoing(Frame{
				MessageType: AREQ,
				Subsystem:   ZDO,
				CommandID:   ZdoMgmtLeaveRspID,
				Payload:     []byte{0x00, 0x40, 0x00},
			})
		}()

		zstack.nodeTable.addOrUpdate(zigbee.IEEEAddress(1), zigbee.NetworkAddress(0x4000))

		err := zstack.RemoveNode(ctx, zigbee.IEEEAddress(1))
		assert.NoError(t, err)

		leaveReq := ZdoMgmtLeaveReq{}
		bytecodec.Unmarshal(call.CapturedCalls[0].Frame.Payload, &leaveReq)

		assert.Equal(t, zigbee.IEEEAddress(1), leaveReq.IEEEAddress)
		assert.Equal(t, zigbee.NetworkAddress(0x4000), leaveReq.NetworkAddress)
		assert.False(t, leaveReq.RemoveChildren)

		unpiMock.AssertCalls(t)
	})
}

func Test_RemoveMessages(t *testing.T) {
	t.Run("verify ZdoMgmtLeaveReq marshals", func(t *testing.T) {
		req := ZdoMgmtLeaveReq{
			NetworkAddress: 0x1234,
			IEEEAddress:    zigbee.IEEEAddress(0x8899aabbccddeeff),
			RemoveChildren: true,
		}

		data, err := bytecodec.Marshal(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x34, 0x12, 0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88, 0x01}, data)
	})

	t.Run("verify ZdoMgmtLeaveReqReply marshals", func(t *testing.T) {
		req := ZdoMgmtLeaveReqReply{
			Status: 1,
		}

		data, err := bytecodec.Marshal(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x01}, data)
	})

	t.Run("ZdoMgmtLeaveReqReply returns true if success", func(t *testing.T) {
		g := ZdoMgmtLeaveReqReply{Status: ZSuccess}
		assert.True(t, g.WasSuccessful())
	})

	t.Run("ZdoMgmtLeaveReqReply returns false if not success", func(t *testing.T) {
		g := ZdoMgmtLeaveReqReply{Status: ZFailure}
		assert.False(t, g.WasSuccessful())
	})

	t.Run("verify ZdoMgmtLeaveRsp marshals", func(t *testing.T) {
		req := ZdoMgmtLeaveRsp{
			SourceAddress: zigbee.NetworkAddress(0x2000),
			Status:        1,
		}

		data, err := bytecodec.Marshal(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x00, 0x20, 0x01}, data)
	})

	t.Run("ZdoMgmtLeaveRsp returns true if success", func(t *testing.T) {
		g := ZdoMgmtLeaveRsp{Status: ZSuccess}
		assert.True(t, g.WasSuccessful())
	})

	t.Run("ZdoMgmtLeaveRsp returns false if not success", func(t *testing.T) {
		g := ZdoMgmtLeaveRsp{Status: ZFailure}
		assert.False(t, g.WasSuccessful())
	})
}
