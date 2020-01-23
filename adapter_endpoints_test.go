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

func TestZStack_RegisterAdapterEndpoint(t *testing.T) {
	t.Run("registers the endpoint", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock)
		defer unpiMock.Stop()

		c := unpiMock.On(SREQ, AF, AFRegisterReqID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   AF,
			CommandID:   AFRegisterRespID,
			Payload:     []byte{0x00},
		})

		err := zstack.RegisterAdapterEndpoint(ctx, 0x01, 0x0104, 0x0001, 0x01, []zigbee.ZCLClusterID{0x0001}, []zigbee.ZCLClusterID{0x0002})
		assert.NoError(t, err)

		unpiMock.AssertCalls(t)

		frame := c.CapturedCalls[0].Frame

		assert.Equal(t, []byte{0x01, 0x04, 0x01, 0x01, 0x00, 0x01, 0x00, 0x01, 0x01, 0x00, 0x01, 0x02, 0x00}, frame.Payload)
	})

	t.Run("returns an error if the query fails", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock)
		defer unpiMock.Stop()

		unpiMock.On(SREQ, AF, AFRegisterReqID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   AF,
			CommandID:   AFRegisterRespID,
			Payload:     []byte{0x01},
		})

		err := zstack.RegisterAdapterEndpoint(ctx, 0x01, 0x0104, 0x0001, 0x01, []zigbee.ZCLClusterID{0x0001}, []zigbee.ZCLClusterID{0x0002})
		assert.Error(t, err)
		assert.Equal(t, ErrorZFailure, err)

		unpiMock.AssertCalls(t)
	})
}

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
