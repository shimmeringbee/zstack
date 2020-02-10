package zstack

import (
	"github.com/shimmeringbee/bytecodec"
	"github.com/shimmeringbee/zigbee"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_IEEEMessages(t *testing.T) {
	t.Run("verify ZdoIEEEAddrReq marshals", func(t *testing.T) {
		req := ZdoIEEEAddrReq{
			NetworkAddress: 0x2040,
			ReqType:        0x01,
			StartIndex:     0x02,
		}

		data, err := bytecodec.Marshall(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x40, 0x20, 0x01, 0x02}, data)
	})

	t.Run("verify ZdoIEEEAddrReqReply marshals", func(t *testing.T) {
		req := ZdoIEEEAddrReqReply{
			Status: 1,
		}

		data, err := bytecodec.Marshall(req)

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
			AssociatedDevices: []zigbee.NetworkAddress{ 0x2002, 0x3003 },
		}

		data, err := bytecodec.Marshall(req)

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
