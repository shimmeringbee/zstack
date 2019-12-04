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
	t.Run("issues a LQI poll request only for coordinators or routers", func(t *testing.T) {
		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock)
		defer unpiMock.Stop()
		defer zstack.Stop()

		unpiMock.On(SREQ, ZDO, ZdoMGMTLQIReqID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   ZDO,
			CommandID:   ZdoMGMTLQIReqRespID,
			Payload:     []byte{0x00},
		}).Times(2)

		zstack.devices[zigbee.IEEEAddress(1)] = &Device{NetworkAddress: 1, IEEEAddress: 1, Role: RoleRouter}
		zstack.devices[zigbee.IEEEAddress(2)] = &Device{NetworkAddress: 2, IEEEAddress: 2, Role: RoleUnknown}

		zstack.startNetworkManager()
		defer zstack.stopNetworkManager()

		time.Sleep(10 * time.Millisecond)

		unpiMock.AssertCalls(t)
	})

	t.Run("the coordinator is added to the device list as a coordinator", func(t *testing.T) {
		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock)
		defer unpiMock.Stop()
		defer zstack.Stop()

		expectedIEEE := zigbee.IEEEAddress(0x0002)
		expectedAddress := zigbee.NetworkAddress(0x0001)

		zstack.NetworkProperties.NetworkAddress = expectedAddress
		zstack.NetworkProperties.IEEEAddress = expectedIEEE

		unpiMock.On(SREQ, ZDO, ZdoMGMTLQIReqID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   ZDO,
			CommandID:   ZdoMGMTLQIReqRespID,
			Payload:     []byte{0x00},
		})

		zstack.startNetworkManager()
		defer zstack.stopNetworkManager()

		time.Sleep(10 * time.Millisecond)

		device, found := zstack.devices[expectedIEEE]

		assert.True(t, found)
		assert.Equal(t, expectedAddress, device.NetworkAddress)
		assert.Equal(t, RoleCoordinator, device.Role)

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
			Capabilities:   0b00000010,
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

		device, found := zstack.devices[announce.IEEEAddress]
		assert.True(t, found)
		assert.Equal(t, announce.IEEEAddress, device.IEEEAddress)
		assert.Equal(t, announce.NetworkAddress, device.NetworkAddress)
		assert.Equal(t, RoleRouter, device.Role)
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
		}).UnlimitedTimes()

		zstack.startNetworkManager()
		defer zstack.stopNetworkManager()

		announce := ZdoLeaveInd{
			SourceAddress: zigbee.NetworkAddress(0x2000),
			IEEEAddress:   zigbee.IEEEAddress(0x0102030405060708),
		}

		zstack.devices[announce.IEEEAddress] = &Device{
			NetworkAddress: 1234,
			IEEEAddress:    announce.IEEEAddress,
			Role:           RoleUnknown,
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

		_, found := zstack.devices[announce.IEEEAddress]
		assert.False(t, found)
	})

}
