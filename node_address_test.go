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

func Test_ResolveNodeIEEEAddress(t *testing.T) {
	t.Run("returns immediately if result is in cache, with no interaction unpi", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		defer unpiMock.AssertCalls(t)
		zstack := New(unpiMock, NewNodeTable())
		defer unpiMock.Stop()

		zstack.nodeTable.addOrUpdate(0x1122334455667788, 0xaabb)

		ieee, err := zstack.ResolveNodeIEEEAddress(ctx, zigbee.NetworkAddress(0xaabb))
		assert.NoError(t, err)
		assert.Equal(t, zigbee.IEEEAddress(0x1122334455667788), ieee)
	})

	t.Run("returns result of successful call to QueryNodeIEEEAddress if not in cache", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		defer unpiMock.AssertCalls(t)
		zstack := New(unpiMock, NewNodeTable())
		defer unpiMock.Stop()

		call := unpiMock.On(SREQ, ZDO, ZdoIEEEAddrReqID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   ZDO,
			CommandID:   ZdoIEEEAddrReqReplyID,
			Payload:     []byte{0x00},
		})

		go func() {
			time.Sleep(10 * time.Millisecond)
			unpiMock.InjectOutgoing(Frame{
				MessageType: AREQ,
				Subsystem:   ZDO,
				CommandID:   ZdoIEEEAddrRspID,
				Payload:     []byte{0x00, 0x88, 0x77, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11, 0x00, 0x40, 0x00, 0x00},
			})
		}()

		ieee, err := zstack.ResolveNodeIEEEAddress(ctx, zigbee.NetworkAddress(0x4000))
		assert.NoError(t, err)
		assert.Equal(t, zigbee.IEEEAddress(0x1122334455667788), ieee)

		addressReq := ZdoIEEEAddrReq{}
		bytecodec.Unmarshal(call.CapturedCalls[0].Frame.Payload, &addressReq)

		assert.Equal(t, zigbee.NetworkAddress(0x4000), addressReq.NetworkAddress)
		assert.Equal(t, uint8(0), addressReq.ReqType)
		assert.Equal(t, uint8(0), addressReq.StartIndex)
	})
}

func Test_QueryNodeIEEEAddress(t *testing.T) {
	t.Run("returns an success on query, response for requested network address is received", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		defer unpiMock.AssertCalls(t)
		zstack := New(unpiMock, NewNodeTable())
		defer unpiMock.Stop()

		call := unpiMock.On(SREQ, ZDO, ZdoIEEEAddrReqID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   ZDO,
			CommandID:   ZdoIEEEAddrReqReplyID,
			Payload:     []byte{0x00},
		})

		go func() {
			time.Sleep(10 * time.Millisecond)
			unpiMock.InjectOutgoing(Frame{
				MessageType: AREQ,
				Subsystem:   ZDO,
				CommandID:   ZdoIEEEAddrRspID,
				Payload:     []byte{0x00, 0x88, 0x77, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11, 0x00, 0x40, 0x00, 0x00},
			})
		}()

		ieee, err := zstack.QueryNodeIEEEAddress(ctx, zigbee.NetworkAddress(0x4000))
		assert.NoError(t, err)
		assert.Equal(t, zigbee.IEEEAddress(0x1122334455667788), ieee)

		addressReq := ZdoIEEEAddrReq{}
		bytecodec.Unmarshal(call.CapturedCalls[0].Frame.Payload, &addressReq)

		assert.Equal(t, zigbee.NetworkAddress(0x4000), addressReq.NetworkAddress)
		assert.Equal(t, uint8(0), addressReq.ReqType)
		assert.Equal(t, uint8(0), addressReq.StartIndex)
	})
}

func Test_ResolveNodeNWKAddress(t *testing.T) {
	t.Run("returns immediately if result is in cache, with no interaction unpi", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		defer unpiMock.AssertCalls(t)
		zstack := New(unpiMock, NewNodeTable())
		defer unpiMock.Stop()

		zstack.nodeTable.addOrUpdate(0x1122334455667788, 0xaabb)

		nwk, err := zstack.ResolveNodeNWKAddress(ctx, zigbee.IEEEAddress(0x1122334455667788))
		assert.NoError(t, err)
		assert.Equal(t, zigbee.NetworkAddress(0xaabb), nwk)
	})

	t.Run("returns result of successful call to QueryNodeNWKAddress if not in cache", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		defer unpiMock.AssertCalls(t)
		zstack := New(unpiMock, NewNodeTable())
		defer unpiMock.Stop()

		call := unpiMock.On(SREQ, ZDO, ZdoNWKAddrReqID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   ZDO,
			CommandID:   ZdoNWKAddrReqReplyID,
			Payload:     []byte{0x00},
		})

		go func() {
			time.Sleep(10 * time.Millisecond)
			unpiMock.InjectOutgoing(Frame{
				MessageType: AREQ,
				Subsystem:   ZDO,
				CommandID:   ZdoNWKAddrRspID,
				Payload:     []byte{0x00, 0x88, 0x77, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11, 0x00, 0x40, 0x00, 0x00},
			})
		}()

		NWK, err := zstack.ResolveNodeNWKAddress(ctx, zigbee.IEEEAddress(0x1122334455667788))
		assert.NoError(t, err)
		assert.Equal(t, zigbee.NetworkAddress(0x4000), NWK)

		addressReq := ZdoNWKAddrReq{}
		bytecodec.Unmarshal(call.CapturedCalls[0].Frame.Payload, &addressReq)

		assert.Equal(t, zigbee.IEEEAddress(0x1122334455667788), addressReq.IEEEAddress)
		assert.Equal(t, uint8(0), addressReq.ReqType)
		assert.Equal(t, uint8(0), addressReq.StartIndex)
	})
}

func Test_QueryNodeNWKAddress(t *testing.T) {
	t.Run("returns an success on query, response for requested network address is received", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		defer unpiMock.AssertCalls(t)
		zstack := New(unpiMock, NewNodeTable())
		defer unpiMock.Stop()

		call := unpiMock.On(SREQ, ZDO, ZdoNWKAddrReqID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   ZDO,
			CommandID:   ZdoNWKAddrReqReplyID,
			Payload:     []byte{0x00},
		})

		go func() {
			time.Sleep(10 * time.Millisecond)
			unpiMock.InjectOutgoing(Frame{
				MessageType: AREQ,
				Subsystem:   ZDO,
				CommandID:   ZdoNWKAddrRspID,
				Payload:     []byte{0x00, 0x88, 0x77, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11, 0x00, 0x40, 0x00, 0x00},
			})
		}()

		ieee, err := zstack.QueryNodeNWKAddress(ctx, zigbee.IEEEAddress(0x1122334455667788))
		assert.NoError(t, err)
		assert.Equal(t, zigbee.NetworkAddress(0x4000), ieee)

		addressReq := ZdoNWKAddrReq{}
		bytecodec.Unmarshal(call.CapturedCalls[0].Frame.Payload, &addressReq)

		assert.Equal(t, zigbee.IEEEAddress(0x1122334455667788), addressReq.IEEEAddress)
		assert.Equal(t, uint8(0), addressReq.ReqType)
		assert.Equal(t, uint8(0), addressReq.StartIndex)
	})
}

func Test_IEEEMessages(t *testing.T) {
	t.Run("verify ZdoIEEEAddrReq marshals", func(t *testing.T) {
		req := ZdoIEEEAddrReq{
			NetworkAddress: 0x2040,
			ReqType:        0x01,
			StartIndex:     0x02,
		}

		data, err := bytecodec.Marshal(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x40, 0x20, 0x01, 0x02}, data)
	})

	t.Run("verify ZdoIEEEAddrReqReply marshals", func(t *testing.T) {
		req := ZdoIEEEAddrReqReply{
			Status: 1,
		}

		data, err := bytecodec.Marshal(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x01}, data)
	})

	t.Run("ZdoIEEEAddrReqReply returns true if success", func(t *testing.T) {
		g := ZdoIEEEAddrReqReply{Status: ZSuccess}
		assert.True(t, g.WasSuccessful())
	})

	t.Run("ZdoIEEEAddrReqReply returns false if not success", func(t *testing.T) {
		g := ZdoIEEEAddrReqReply{Status: ZFailure}
		assert.False(t, g.WasSuccessful())
	})

	t.Run("verify ZdoIEEEAddrRsp marshals", func(t *testing.T) {
		req := ZdoIEEEAddrRsp{
			Status:            0x01,
			IEEEAddress:       0x1122334455667788,
			NetworkAddress:    0xaabb,
			StartIndex:        0x02,
			AssociatedDevices: []zigbee.NetworkAddress{0x2002, 0x3003},
		}

		data, err := bytecodec.Marshal(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x01, 0x88, 0x77, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11, 0xbb, 0xaa, 0x02, 0x02, 0x02, 0x20, 0x03, 0x30}, data)
	})

	t.Run("ZdoIEEEAddrRsp returns true if success", func(t *testing.T) {
		g := ZdoIEEEAddrRsp{Status: ZSuccess}
		assert.True(t, g.WasSuccessful())
	})

	t.Run("ZdoIEEEAddrRsp returns false if not success", func(t *testing.T) {
		g := ZdoIEEEAddrRsp{Status: ZFailure}
		assert.False(t, g.WasSuccessful())
	})
}

func Test_NWKMessages(t *testing.T) {
	t.Run("verify ZdoNWKAddrReq marshals", func(t *testing.T) {
		req := ZdoNWKAddrReq{
			IEEEAddress: 0x1122334455667788,
			ReqType:     0x01,
			StartIndex:  0x02,
		}

		data, err := bytecodec.Marshal(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x88, 0x77, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11, 0x01, 0x02}, data)
	})

	t.Run("verify ZdoNWKAddrReqReply marshals", func(t *testing.T) {
		req := ZdoNWKAddrReqReply{
			Status: 1,
		}

		data, err := bytecodec.Marshal(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x01}, data)
	})

	t.Run("ZdoNWKAddrReqReply returns true if success", func(t *testing.T) {
		g := ZdoNWKAddrReqReply{Status: ZSuccess}
		assert.True(t, g.WasSuccessful())
	})

	t.Run("ZdoNWKAddrReqReply returns false if not success", func(t *testing.T) {
		g := ZdoNWKAddrReqReply{Status: ZFailure}
		assert.False(t, g.WasSuccessful())
	})

	t.Run("verify ZdoNWKAddrRsp marshals", func(t *testing.T) {
		req := ZdoNWKAddrRsp{
			Status:            0x01,
			IEEEAddress:       0x1122334455667788,
			NetworkAddress:    0xaabb,
			StartIndex:        0x02,
			AssociatedDevices: []zigbee.NetworkAddress{0x2002, 0x3003},
		}

		data, err := bytecodec.Marshal(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x01, 0x88, 0x77, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11, 0xbb, 0xaa, 0x02, 0x02, 0x02, 0x20, 0x03, 0x30}, data)
	})

	t.Run("ZdoNWKAddrRsp returns true if success", func(t *testing.T) {
		g := ZdoNWKAddrRsp{Status: ZSuccess}
		assert.True(t, g.WasSuccessful())
	})

	t.Run("ZdoNWKAddrRsp returns false if not success", func(t *testing.T) {
		g := ZdoNWKAddrRsp{Status: ZFailure}
		assert.False(t, g.WasSuccessful())
	})
}
