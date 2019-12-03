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

	t.Run("emits DeviceJoin event when device join announcement received", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock)
		defer unpiMock.Stop()
		defer unpiMock.AssertCalls(t)

		unpiMock.On(SREQ, ZDO, ZdoMGMTLQIReqID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   ZDO,
			CommandID:   ZdoMGMTLQIReqRespID,
			Payload:     []byte{0x00},
		})

		zstack.startNetworkManager()
		defer zstack.stopNetworkManager()

		time.Sleep(10 * time.Millisecond)

		announce := ZdoEndDeviceAnnceInd{
			SourceAddress:  zigbee.NetworkAddress(0x1000),
			NetworkAddress: zigbee.NetworkAddress(0x2000),
			IEEEAddress:    zigbee.IEEEAddress(0x0102030405060708),
			Capabilities:   181,
		}

		data, _ := bytecodec.Marshall(announce)

		unpiMock.InjectOutgoing(Frame{
			MessageType: AREQ,
			Subsystem:   ZDO,
			CommandID:   ZdoEndDeviceAnnceIndID,
			Payload:     data,
		})

		event, err := zstack.ReadEvent(ctx)
		assert.NoError(t, err)

		deviceJoin, ok := event.(DeviceJoinEvent)

		assert.True(t, ok)

		assert.Equal(t, announce.NetworkAddress, deviceJoin.NetworkAddress)
		assert.Equal(t, announce.IEEEAddress, deviceJoin.IEEEAddress)
	})

	t.Run("emits DeviceLeave event when device leave announcement received", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock)
		defer unpiMock.Stop()
		defer unpiMock.AssertCalls(t)

		unpiMock.On(SREQ, ZDO, ZdoMGMTLQIReqID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   ZDO,
			CommandID:   ZdoMGMTLQIReqRespID,
			Payload:     []byte{0x00},
		})

		zstack.startNetworkManager()
		defer zstack.stopNetworkManager()

		announce := ZdoLeaveInd{
			SourceAddress: zigbee.NetworkAddress(0x2000),
			IEEEAddress:   zigbee.IEEEAddress(0x0102030405060708),
		}

		data, _ := bytecodec.Marshall(announce)

		unpiMock.InjectOutgoing(Frame{
			MessageType: AREQ,
			Subsystem:   ZDO,
			CommandID:   ZdoLeaveIndID,
			Payload:     data,
		})

		event, err := zstack.ReadEvent(ctx)
		assert.NoError(t, err)

		deviceLeave, ok := event.(DeviceLeaveEvent)

		assert.True(t, ok)

		assert.Equal(t, announce.SourceAddress, deviceLeave.NetworkAddress)
		assert.Equal(t, announce.IEEEAddress, deviceLeave.IEEEAddress)
	})

}
