package zstack

import (
	. "github.com/shimmeringbee/unpi"
	unpiTest "github.com/shimmeringbee/unpi/testing"
	"testing"
	"time"
)

func Test_NetworkManager(t *testing.T) {
	t.Run("issues a LQI poll request", func(t *testing.T) {
		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock)
		defer unpiMock.Stop()
		defer zstack.Stop()

		unpiMock.On(SREQ, ZDO, ZdoMGMTLQIReqID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   ZDO,
			CommandID:   ZdoMGMTLQIReqRespID,
			Payload:     []byte{0x00},
		})

		zstack.startNetworkManager()
		defer zstack.stopNetworkManager()

		time.Sleep(10 * time.Millisecond)

		unpiMock.AssertCalls(t)
	})
}