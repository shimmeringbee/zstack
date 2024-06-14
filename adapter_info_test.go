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

func Test_GetAdapterIEEEAddress(t *testing.T) {
	t.Run("gets the IEEE address", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock, memory.New())
		zstack.sem = semaphore.NewWeighted(8)
		defer unpiMock.Stop()

		unpiMock.On(SREQ, UTIL, UtilGetDeviceInfoRequestID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   UTIL,
			CommandID:   UtilGetDeviceInfoRequestReplyID,
			Payload:     []byte{0x00, 0x09, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x12, 0x11},
		})

		address, err := zstack.GetAdapterIEEEAddress(ctx)
		assert.NoError(t, err)
		assert.Equal(t, zigbee.IEEEAddress(0x0203040506070809), address)

		unpiMock.AssertCalls(t)
	})
}

func Test_GetAdapterNetworkAddress(t *testing.T) {
	t.Run("gets the network address", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock, memory.New())
		zstack.sem = semaphore.NewWeighted(8)
		defer unpiMock.Stop()

		unpiMock.On(SREQ, UTIL, UtilGetDeviceInfoRequestID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   UTIL,
			CommandID:   UtilGetDeviceInfoRequestReplyID,
			Payload:     []byte{0x00, 0x09, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x12, 0x11},
		})

		address, err := zstack.GetAdapterNetworkAddress(ctx)
		assert.NoError(t, err)
		assert.Equal(t, zigbee.NetworkAddress(0x1112), address)

		unpiMock.AssertCalls(t)
	})
}

func Test_UtilGetDeviceInfoStructs(t *testing.T) {
	t.Run("UtilGetDeviceInfoRequestReply", func(t *testing.T) {
		s := UtilGetDeviceInfoRequestReply{
			Status:         0x01,
			IEEEAddress:    0x0203040506070809,
			NetworkAddress: 0x1112,
		}

		actualBytes, err := bytecodec.Marshal(s)

		expectedBytes := []byte{0x01, 0x09, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x12, 0x11}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})
}
