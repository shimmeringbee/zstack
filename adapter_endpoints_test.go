package zstack

import (
	"context"
	"github.com/shimmeringbee/bytecodec"
	"github.com/shimmeringbee/persistence/impl/memory"
	. "github.com/shimmeringbee/unpi"
	unpiTest "github.com/shimmeringbee/unpi/testing"
	"github.com/shimmeringbee/zigbee"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/semaphore"
	"testing"
	"time"
)

func Test_RegisterAdapterEndpoint(t *testing.T) {
	t.Run("registers the endpoint", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock, memory.New())
		zstack.sem = semaphore.NewWeighted(8)
		defer unpiMock.Stop()

		c := unpiMock.On(SREQ, AF, AFRegisterID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   AF,
			CommandID:   AFRegisterReplyID,
			Payload:     []byte{0x00},
		})

		err := zstack.RegisterAdapterEndpoint(ctx, 0x01, 0x0104, 0x0001, 0x01, []zigbee.ClusterID{0x0001}, []zigbee.ClusterID{0x0002})
		assert.NoError(t, err)

		unpiMock.AssertCalls(t)

		frame := c.CapturedCalls[0].Frame

		assert.Equal(t, []byte{0x01, 0x04, 0x01, 0x01, 0x00, 0x01, 0x00, 0x01, 0x01, 0x00, 0x01, 0x02, 0x00}, frame.Payload)
	})

	t.Run("returns an error if the query fails", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock, memory.New())
		zstack.sem = semaphore.NewWeighted(8)
		defer unpiMock.Stop()

		unpiMock.On(SREQ, AF, AFRegisterID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   AF,
			CommandID:   AFRegisterReplyID,
			Payload:     []byte{0x01},
		})

		err := zstack.RegisterAdapterEndpoint(ctx, 0x01, 0x0104, 0x0001, 0x01, []zigbee.ClusterID{0x0001}, []zigbee.ClusterID{0x0002})
		assert.Error(t, err)
		assert.Equal(t, ErrorZFailure, err)

		unpiMock.AssertCalls(t)
	})
}

func Test_EndpointRegisterMessages(t *testing.T) {
	t.Run("verify AFRegister marshals", func(t *testing.T) {
		req := AFRegister{
			Endpoint:         1,
			AppProfileId:     2,
			AppDeviceId:      3,
			AppDeviceVersion: 4,
			LatencyReq:       5,
			AppInClusters:    []zigbee.ClusterID{0x10},
			AppOutClusters:   []zigbee.ClusterID{0x20},
		}

		data, err := bytecodec.Marshal(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x01, 0x02, 0x00, 0x03, 0x00, 0x04, 0x05, 0x01, 0x10, 0x00, 0x01, 0x20, 0x00}, data)
	})

	t.Run("verify AFRegisterReply marshals", func(t *testing.T) {
		req := AFRegisterReply{
			Status: 1,
		}

		data, err := bytecodec.Marshal(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x01}, data)
	})
}
