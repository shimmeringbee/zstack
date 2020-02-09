package zstack

import (
	"github.com/shimmeringbee/bytecodec"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_IncommingMessage(t *testing.T) {
	t.Run("verify AfIncomingMsg marshals", func(t *testing.T) {
		req := AfIncomingMsg{
			GroupID:             0x0102,
			ClusterID:           0x0304,
			SourceAddress:       0x0506,
			SourceEndpoint:      0x07,
			DestinationEndpoint: 0x08,
			WasBroadcast:        0x09,
			LinkQuality:         0x0a,
			SecurityUse:         0x0b,
			TimeStamp:           0x10111213,
			Sequence:            0x0d,
			Data:                []byte { 0x00, 0x01, 0x02 },
		}

		data, err := bytecodec.Marshall(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x02, 0x01, 0x04, 0x03, 0x06, 0x05, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x13, 0x12, 0x11, 0x10, 0x0d, 0x03, 0x00, 0x01, 0x02}, data)
	})
}