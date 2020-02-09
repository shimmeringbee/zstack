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

func Test_QueryNodeDescription(t *testing.T) {
	t.Run("returns an success on query, response for requested network address is received", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock)
		defer unpiMock.Stop()

		unpiMock.On(SREQ, ZDO, ZdoNodeDescReqID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   ZDO,
			CommandID:   ZdoNodeDescReqReplyID,
			Payload:     []byte{0x00},
		})

		go func() {
			time.Sleep(10 * time.Millisecond)
			unpiMock.InjectOutgoing(Frame{
				MessageType: AREQ,
				Subsystem:   ZDO,
				CommandID:   ZdoNodeDescRspID,
				Payload:     []byte{0x00, 0x20, 0x01, 0x00, 0x40, 0x40, 0x02, 0x03, 0x05, 0x04, 0x06, 0x08, 0x07, 0x0a, 0x09, 0x0c, 0x0b, 0x0d},
			})
		}()

		nodeDescription, err := zstack.QueryNodeDescription(ctx, zigbee.NetworkAddress(0x4000))
		assert.NoError(t, err)
		assert.Equal(t, zigbee.NodeDescription{
			LogicalType:      zigbee.EndDevice,
			ManufacturerCode: 0x0405,
		}, nodeDescription)

		unpiMock.AssertCalls(t)
	})
}

func Test_NodeDescriptionMessages(t *testing.T) {
	t.Run("verify ZdoNodeDescReq marshals", func(t *testing.T) {
		req := ZdoNodeDescReq{
			DestinationAddress: zigbee.NetworkAddress(0x2000),
			OfInterestAddress:  zigbee.NetworkAddress(0x4000),
		}

		data, err := bytecodec.Marshall(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x00, 0x20, 0x00, 0x40}, data)
	})

	t.Run("verify ZdoNodeDescReqReply marshals", func(t *testing.T) {
		req := ZdoNodeDescReqReply{
			Status: 1,
		}

		data, err := bytecodec.Marshall(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x01}, data)
	})

	t.Run("ZdoNodeDescReqReply returns true if success", func(t *testing.T) {
		g := ZdoNodeDescReqReply{Status: ZSuccess}
		assert.True(t, g.WasSuccessful())
	})

	t.Run("ZdoNodeDescReqReply returns false if not success", func(t *testing.T) {
		g := ZdoNodeDescReqReply{Status: ZFailure}
		assert.False(t, g.WasSuccessful())
	})

	t.Run("verify ZdoNodeDescRsp marshals", func(t *testing.T) {
		req := ZdoNodeDescRsp{
			SourceAddress:          zigbee.NetworkAddress(0x2000),
			Status:                 1,
			OfInterestAddress:      zigbee.NetworkAddress(0x4000),
			LogicalTypeDescriptor:  0x01,
			APSFlagsFrequency:      0x02,
			MacCapabilitiesFlags:   0x03,
			ManufacturerCode:       0x0405,
			MaxBufferSize:          0x06,
			MaxInTransferSize:      0x0708,
			ServerMask:             0x090a,
			MaxOutTransferSize:     0x0b0c,
			DescriptorCapabilities: 0x0d,
		}

		data, err := bytecodec.Marshall(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x00, 0x20, 0x01, 0x00, 0x40, 0x01, 0x02, 0x03, 0x05, 0x04, 0x06, 0x08, 0x07, 0x0a, 0x09, 0x0c, 0x0b, 0x0d}, data)
	})

	t.Run("ZdoNodeDescRsp returns true if success", func(t *testing.T) {
		g := ZdoNodeDescRsp{Status: ZSuccess}
		assert.True(t, g.WasSuccessful())
	})

	t.Run("ZdoNodeDescRsp returns false if not success", func(t *testing.T) {
		g := ZdoNodeDescRsp{Status: ZFailure}
		assert.False(t, g.WasSuccessful())
	})
}
