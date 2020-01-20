package zstack

import (
	"github.com/shimmeringbee/bytecodec"
	"github.com/shimmeringbee/zigbee"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_endpointRegisterMessages(t *testing.T) {
	t.Run("verify AFRegisterReq marshals", func(t *testing.T) {
		req := AFRegisterReq{
			Endpoint:         1,
			AppProfileId:     2,
			AppDeviceId:      3,
			AppDeviceVersion: 4,
			LatencyReq:       5,
			AppInClusters:    []zigbee.ZCLClusterID{0x10},
			AppOutClusters:   []zigbee.ZCLClusterID{0x20},
		}

		data, err := bytecodec.Marshall(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x01, 0x02, 0x00, 0x03, 0x00, 0x04, 0x05, 0x01, 0x10, 0x00, 0x01, 0x20, 0x00}, data)
	})

	t.Run("verify AFRegisterResp marshals", func(t *testing.T) {
		req := AFRegisterResp{
			Status: 1,
		}

		data, err := bytecodec.Marshall(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x01}, data)
	})
}