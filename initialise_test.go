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

func TestZStack_Initialise(t *testing.T) {
	t.Run("test initialisation process", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
		defer cancel()

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

		resetOn := unpiMock.On(AREQ, SYS, SysResetReqID).Return(Frame{
			MessageType: AREQ,
			Subsystem:   SYS,
			CommandID:   SysResetIndID,
			Payload:     resetResponse,
		}).Times(3)

		nvramWriteResponse, _ := bytecodec.Marshall(SysOSALNVWriteResp{Status: ZSuccess})
		nvramOn := unpiMock.On(SREQ, SYS, SysOSALNVWriteReqID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   SYS,
			CommandID:   SysOSALNVWriteRespID,
			Payload:     nvramWriteResponse,
		}).Times(11)

		nc := zigbee.NetworkConfiguration{
			PANID:         [2]byte{0x01, 0x02},
			ExtendedPANID: [8]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			NetworkKey:    [16]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			Channel:       zigbee.DefaultChannel,
		}
		err := zstack.Initialise(ctx, nc)

		assert.NoError(t, err)
		unpiMock.AssertCalls(t)

		assert.Equal(t, []byte{0x01}, resetOn.CapturedCalls[0].Frame.Payload)
		assert.Equal(t, []byte{0x03, 0x00, 0x00, 0x01, 0x03}, nvramOn.CapturedCalls[0].Frame.Payload)
		assert.Equal(t, []byte{0x01}, resetOn.CapturedCalls[1].Frame.Payload)
		assert.Equal(t, []byte{0x87, 0x00, 0x00, 0x01, 0x00}, nvramOn.CapturedCalls[1].Frame.Payload)
		assert.Equal(t, []byte{0x01}, resetOn.CapturedCalls[2].Frame.Payload)
		assert.Equal(t, []byte{0x64, 0x00, 0x00, 0x01, 0x1}, nvramOn.CapturedCalls[2].Frame.Payload)
		assert.Equal(t, []byte{0x63, 0x00, 0x00, 0x01, 0x1}, nvramOn.CapturedCalls[3].Frame.Payload)
		assert.Equal(t, []byte{0x62, 0x00, 0x00, 0x10, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}, nvramOn.CapturedCalls[4].Frame.Payload)
		assert.Equal(t, []byte{0x8f, 0x00, 0x00, 0x01, 0x01}, nvramOn.CapturedCalls[5].Frame.Payload)
		assert.Equal(t, []byte{0x84, 0x00, 0x00, 0x04, 0x00, 0x00, 0x80, 0x00}, nvramOn.CapturedCalls[6].Frame.Payload)
		assert.Equal(t, []byte{0x83, 0x00, 0x00, 0x02, 0x01, 0x02}, nvramOn.CapturedCalls[7].Frame.Payload)
		assert.Equal(t, []byte{0x2d, 0x00, 0x00, 0x08, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}, nvramOn.CapturedCalls[8].Frame.Payload)
		assert.Equal(t, []byte{0x6d, 0x00, 0x00, 0x01, 0x01}, nvramOn.CapturedCalls[9].Frame.Payload)
		assert.Equal(t, []byte{0x01, 0x01, 0x00, 0x20, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x5a, 0x69, 0x67, 0x42, 0x65, 0x65, 0x41, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x30, 0x39, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, nvramOn.CapturedCalls[10].Frame.Payload)
	})
}

func TestZStack_startZigbeeStack(t *testing.T) {
	t.Run("starts zigbee stack and waits for confirmation", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock)
		defer unpiMock.Stop()

		unpiMock.On(SREQ, SAPI, SAPIZBStartRequestID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   SAPI,
			CommandID:   SAPIZBStartResponseID,
			Payload:     nil,
		})

		go func() {
			time.Sleep(50 * time.Millisecond)
			unpiMock.InjectOutgoing(Frame{
				MessageType: AREQ,
				Subsystem:   SAPI,
				CommandID:   SAPIZBStartConfirmID,
				Payload:     []byte{0x00},
			})
		}()

		err := zstack.startZigbeeStack(ctx)
		assert.NoError(t, err)

		unpiMock.AssertCalls(t)
	})

	t.Run("starts zigbee stack errors when confirmation is not ZB_SUCCESS", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock)
		defer unpiMock.Stop()

		unpiMock.On(SREQ, SAPI, SAPIZBStartRequestID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   SAPI,
			CommandID:   SAPIZBStartResponseID,
			Payload:     nil,
		})

		go func() {
			time.Sleep(50 * time.Millisecond)
			unpiMock.InjectOutgoing(Frame{
				MessageType: AREQ,
				Subsystem:   SAPI,
				CommandID:   SAPIZBStartConfirmID,
				Payload:     []byte{0x22},
			})
		}()

		err := zstack.startZigbeeStack(ctx)
		assert.Error(t, err)

		unpiMock.AssertCalls(t)
	})
}