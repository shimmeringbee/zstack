package zstack

import (
	"context"
	"github.com/shimmeringbee/bytecodec"
	. "github.com/shimmeringbee/unpi"
	unpiTest "github.com/shimmeringbee/unpi/testing"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestZStack_Initialise(t *testing.T) {
	t.Run("test initialisation process", func(t *testing.T) {
		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock)
		defer unpiMock.Stop()

		resetResponse, _ := bytecodec.Marshall(SysResetInd{
			Reason:            External,
			TransportRevision: 2,
			ProductID:         1,
			MajorRelease:      2,
			MinorRelease:      3,
			HardwareRevision:  4,
		})

		unpiMock.On(AREQ, SYS, SysResetReqID).Return(Frame{
			MessageType: AREQ,
			Subsystem:   SYS,
			CommandID:   SysResetIndID,
			Payload:     resetResponse,
		})

		nvramWriteResponse, _ := bytecodec.Marshall(SysOSALNVWriteResp{Status: ZSuccess})
		unpiMock.On(SREQ, SYS, SysOSALNVWriteReqID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   SYS,
			CommandID:   SysOSALNVWriteRespID,
			Payload:     nvramWriteResponse,
		})

		err := zstack.Initialise(context.Background())

		assert.NoError(t, err)
		unpiMock.AssertCalls(t)
	})
}
