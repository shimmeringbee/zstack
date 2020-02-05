package zstack

import (
	"github.com/shimmeringbee/bytecodec"
	"github.com/shimmeringbee/zigbee"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_SimpleDescriptionMessages(t *testing.T) {
	t.Run("verify ZdoSimpleDescReq marshals", func(t *testing.T) {
		req := ZdoSimpleDescReq{
			DestinationAddress: zigbee.NetworkAddress(0x2000),
			OfInterestAddress:  zigbee.NetworkAddress(0x4000),
			Endpoint:           0x08,
		}

		data, err := bytecodec.Marshall(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x00, 0x20, 0x00, 0x40, 0x08}, data)
	})

	t.Run("verify ZdoSimpleDescReqReply marshals", func(t *testing.T) {
		req := ZdoSimpleDescReqReply{
			Status: 1,
		}

		data, err := bytecodec.Marshall(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x01}, data)
	})

	t.Run("verify ZdoSimpleDescRsp marshals", func(t *testing.T) {
		req := ZdoSimpleDescRsp{
			SourceAddress:     zigbee.NetworkAddress(0x2000),
			Status:            1,
			OfInterestAddress: zigbee.NetworkAddress(0x4000),
			Length:            0x0a,
			Endpoint:          0x08,
			ProfileID:         0x1234,
			DeviceID:          0x5678,
			DeviceVersion:     0,
			InClusterList:     []zigbee.ZCLClusterID{0x1234},
			OutClusterList:    []zigbee.ZCLClusterID{0x5678},
		}

		data, err := bytecodec.Marshall(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x00, 0x20, 0x01, 0x00, 0x40, 0x0a, 0x08, 0x34, 0x12, 0x78, 0x56, 0x00, 0x01, 0x34, 0x12, 0x01, 0x78, 0x56}, data)
	})
}
