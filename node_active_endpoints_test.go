package zstack

import (
	"github.com/shimmeringbee/bytecodec"
	"github.com/shimmeringbee/zigbee"
	"github.com/stretchr/testify/assert"
	"testing"
)

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

	t.Run("verify ZdoActiveEpRsp marshals", func(t *testing.T) {
		req := ZdoActiveEpRsp{
			SourceAddress:     zigbee.NetworkAddress(0x2000),
			Status:            1,
			OfInterestAddress: zigbee.NetworkAddress(0x4000),
			ActiveEndpoints:   []byte{0x01, 0x02, 0x03},
		}

		data, err := bytecodec.Marshall(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x00, 0x20, 0x01, 0x00, 0x40, 0x03, 0x01, 0x02, 0x03}, data)
	})
}
