package zstack

import (
	"github.com/shimmeringbee/bytecodec"
	"github.com/shimmeringbee/zigbee"
	"github.com/stretchr/testify/assert"
	"testing"
)

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

		data, err := bytecodec.Marshall(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x21, 0x20, 0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88, 0x01, 0xfe, 0xca, 0x01, 0xfe, 0x3f, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02}, data)
	})

	t.Run("verify ZdoUnbindReqReply marshals", func(t *testing.T) {
		req := ZdoUnbindReqReply{
			Status: 1,
		}

		data, err := bytecodec.Marshall(req)

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
			SourceAddress:     zigbee.NetworkAddress(0x2000),
			Status:            1,
		}

		data, err := bytecodec.Marshall(req)

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
