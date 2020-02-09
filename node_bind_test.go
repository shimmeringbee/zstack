package zstack

import (
	"github.com/shimmeringbee/bytecodec"
	"github.com/shimmeringbee/zigbee"
	"github.com/stretchr/testify/assert"
	"testing"
)

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

		data, err := bytecodec.Marshall(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x21, 0x20, 0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88, 0x01, 0xfe, 0xca, 0x01, 0xfe, 0x3f, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02}, data)
	})

	t.Run("verify ZdoBindReqReply marshals", func(t *testing.T) {
		req := ZdoBindReqReply{
			Status: 1,
		}

		data, err := bytecodec.Marshall(req)

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
			SourceAddress:     zigbee.NetworkAddress(0x2000),
			Status:            1,
		}

		data, err := bytecodec.Marshall(req)

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
