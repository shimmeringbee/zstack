package zstack

import (
	"github.com/shimmeringbee/bytecodec"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_SendMessages(t *testing.T) {
	t.Run("verify AfDataRequest marshals", func(t *testing.T) {
		req := AfDataRequest{
			DestinationAddress:  0x0102,
			DestinationEndpoint: 0x03,
			SourceEndpoint:      0x04,
			ClusterID:           0x0506,
			TransactionID:       0x07,
			Options:             0x08,
			Radius:              0x09,
			Data:                []byte{0x0a, 0x0b},
		}

		data, err := bytecodec.Marshall(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x02, 0x01, 0x03, 0x04, 0x06, 0x05, 0x07, 0x08, 0x09, 0x02, 0x0a, 0x0b}, data)
	})

	t.Run("verify AfDataRequestReply marshals", func(t *testing.T) {
		req := AfDataRequestReply{
			Status: 1,
		}

		data, err := bytecodec.Marshall(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x01}, data)
	})

	t.Run("AfDataRequestReply returns true if success", func(t *testing.T) {
		g := AfDataRequestReply{Status: ZSuccess}
		assert.True(t, g.WasSuccessful())
	})

	t.Run("AfDataRequestReply returns false if not success", func(t *testing.T) {
		g := AfDataRequestReply{Status: ZFailure}
		assert.False(t, g.WasSuccessful())
	})

	t.Run("verify AfDataConfirm marshals", func(t *testing.T) {
		req := AfDataConfirm{
			Status:        0x01,
			Endpoint:      0x02,
			TransactionID: 0x03,
		}

		data, err := bytecodec.Marshall(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x01, 0x02, 0x03}, data)
	})

	t.Run("AfDataConfirm returns true if success", func(t *testing.T) {
		g := AfDataConfirm{Status: ZSuccess}
		assert.True(t, g.WasSuccessful())
	})

	t.Run("AfDataConfirm returns false if not success", func(t *testing.T) {
		g := AfDataConfirm{Status: ZFailure}
		assert.False(t, g.WasSuccessful())
	})
}
