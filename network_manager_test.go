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
			CommandID:   ZdoMGMTLQIReqReplyID,
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
			CommandID:   ZdoMGMTLQIReqReplyID,
			Payload:     []byte{0x00},
		})

		zstack.startNetworkManager()
		defer zstack.stopNetworkManager()

		time.Sleep(10 * time.Millisecond)

		device, found := zstack.devices[expectedIEEE]

		assert.True(t, found)
		assert.Equal(t, expectedAddress, device.NetworkAddress)
		assert.Equal(t, RoleCoordinator, device.Role)

		reverseIEEE, found := zstack.devicesByNetAddr[expectedAddress]
		assert.True(t, found)
		assert.Equal(t, expectedIEEE, reverseIEEE)

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
			CommandID:   ZdoMGMTLQIReqReplyID,
			Payload:     []byte{0x00},
		}).UnlimitedTimes()

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

		reverseIEEE, found := zstack.devicesByNetAddr[announce.NetworkAddress]
		assert.True(t, found)
		assert.Equal(t, announce.IEEEAddress, reverseIEEE)
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
			CommandID:   ZdoMGMTLQIReqReplyID,
			Payload:     []byte{0x00},
		}).UnlimitedTimes()

		zstack.startNetworkManager()
		defer zstack.stopNetworkManager()

		announce := ZdoLeaveInd{
			SourceAddress: zigbee.NetworkAddress(0x2000),
			IEEEAddress:   zigbee.IEEEAddress(0x0102030405060708),
		}

		zstack.devices[zigbee.IEEEAddress(0)] = &Device{ Neighbours: map[zigbee.IEEEAddress]*DeviceNeighbour{}}
		zstack.devices[zigbee.IEEEAddress(0)].Neighbours[announce.IEEEAddress] = &DeviceNeighbour{LQI:50}

		zstack.devices[announce.IEEEAddress] = &Device{
			NetworkAddress: 0x2000,
			IEEEAddress:    announce.IEEEAddress,
			Role:           RoleUnknown,
		}

		zstack.devicesByNetAddr[announce.SourceAddress] = announce.IEEEAddress

		data, _ := bytecodec.Marshall(announce)

		time.Sleep(10 * time.Millisecond)

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

		_, found = zstack.devicesByNetAddr[announce.SourceAddress]
		assert.False(t, found)

		_, found = zstack.devices[zigbee.IEEEAddress(0)].Neighbours[deviceLeave.IEEEAddress]
		assert.False(t, found)
	})
	
	t.Run("a new router will be queried for network state", func(t *testing.T) {
		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock)
		defer unpiMock.Stop()
		defer unpiMock.AssertCalls(t)

		c := unpiMock.On(SREQ, ZDO, ZdoMGMTLQIReqID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   ZDO,
			CommandID:   ZdoMGMTLQIReqReplyID,
			Payload:     []byte{0x00},
		}).Times(2)

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

		time.Sleep(10 * time.Millisecond)

		assert.Equal(t, 2, len(c.CapturedCalls))

		frame := c.CapturedCalls[1]

		lqiReq := ZdoMGMTLQIReq{}
		_ = bytecodec.Unmarshall(frame.Frame.Payload, &lqiReq)

		assert.Equal(t, zigbee.NetworkAddress(0x2000), lqiReq.DestinationAddress)
	})

	t.Run("devices in LQI query are added to network manager", func(t *testing.T) {
		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock)
		defer unpiMock.Stop()
		defer unpiMock.AssertCalls(t)

		unpiMock.On(SREQ, ZDO, ZdoMGMTLQIReqID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   ZDO,
			CommandID:   ZdoMGMTLQIReqReplyID,
			Payload:     []byte{0x00},
		}).UnlimitedTimes()

		zstack.startNetworkManager()
		defer zstack.stopNetworkManager()

		time.Sleep(10 * time.Millisecond)

		announce := ZdoMGMTLQIRsp{
			SourceAddress:         0,
			Status:                0,
			NeighbourTableEntries: 1,
			StartIndex:            0,
			Neighbors:             []ZdoMGMTLQINeighbour{
				{
					ExtendedPANID:  zstack.NetworkProperties.ExtendedPANID,
					IEEEAddress:    zigbee.IEEEAddress(0x1000),
					NetworkAddress: zigbee.NetworkAddress(0x2000),
					Status:         0b00100001,
					PermitJoining:  0,
					Depth:          0,
					LQI:            67,
				},
			},
		}

		data, _ := bytecodec.Marshall(announce)

		unpiMock.InjectOutgoing(Frame{
			MessageType: AREQ,
			Subsystem:   ZDO,
			CommandID:   ZdoMGMTLQIRspID,
			Payload:     data,
		})

		time.Sleep(10 * time.Millisecond)

		device, found := zstack.devices[zigbee.IEEEAddress(0x1000)]
		assert.True(t, found)

		if found {
			assert.Equal(t, zigbee.NetworkAddress(0x2000), device.NetworkAddress)
			assert.Equal(t, RoleRouter, device.Role)
		}

		requestingDevice, _ := zstack.devices[zigbee.IEEEAddress(0)]

		neighbourEntry, found := requestingDevice.Neighbours[zigbee.IEEEAddress(0x1000)]
		assert.True(t, found)
		assert.Equal(t, uint8(67), neighbourEntry.LQI)
		assert.Equal(t, RelationshipSibling, neighbourEntry.Relationship)
	})

	t.Run("devices in LQI query are not added if Ext PANID does not match", func(t *testing.T) {
		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock)
		defer unpiMock.Stop()
		defer unpiMock.AssertCalls(t)

		unpiMock.On(SREQ, ZDO, ZdoMGMTLQIReqID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   ZDO,
			CommandID:   ZdoMGMTLQIReqReplyID,
			Payload:     []byte{0x00},
		}).UnlimitedTimes()

		zstack.startNetworkManager()
		defer zstack.stopNetworkManager()

		time.Sleep(10 * time.Millisecond)

		announce := ZdoMGMTLQIRsp{
			SourceAddress:         0,
			Status:                0,
			NeighbourTableEntries: 1,
			StartIndex:            0,
			Neighbors:             []ZdoMGMTLQINeighbour{
				{
					ExtendedPANID:  0xfffffff,
					IEEEAddress:    zigbee.IEEEAddress(0x2000),
					NetworkAddress: zigbee.NetworkAddress(0x4000),
					Status:         0b00000001,
					PermitJoining:  0,
					Depth:          0,
					LQI:            67,
				},
			},
		}

		data, _ := bytecodec.Marshall(announce)

		unpiMock.InjectOutgoing(Frame{
			MessageType: AREQ,
			Subsystem:   ZDO,
			CommandID:   ZdoMGMTLQIRspID,
			Payload:     data,
		})

		time.Sleep(10 * time.Millisecond)

		_, found := zstack.devices[zigbee.IEEEAddress(0x2000)]
		assert.False(t, found)
	})

	t.Run("devices in LQI query are not added if it has an invalid IEEE address", func(t *testing.T) {
		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock)
		zstack.NetworkProperties.IEEEAddress = zigbee.IEEEAddress(1)

		defer unpiMock.Stop()
		defer unpiMock.AssertCalls(t)

		unpiMock.On(SREQ, ZDO, ZdoMGMTLQIReqID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   ZDO,
			CommandID:   ZdoMGMTLQIReqReplyID,
			Payload:     []byte{0x00},
		}).UnlimitedTimes()

		zstack.startNetworkManager()
		defer zstack.stopNetworkManager()

		time.Sleep(10 * time.Millisecond)

		announce := ZdoMGMTLQIRsp{
			SourceAddress:         0,
			Status:                0,
			NeighbourTableEntries: 1,
			StartIndex:            0,
			Neighbors:             []ZdoMGMTLQINeighbour{
				{
					ExtendedPANID:  zstack.NetworkProperties.ExtendedPANID,
					IEEEAddress:    zigbee.IEEEAddress(0),
					NetworkAddress: zigbee.NetworkAddress(0x4000),
					Status:         0b00000001,
					PermitJoining:  0,
					Depth:          0,
					LQI:            67,
				},
			},
		}

		data, _ := bytecodec.Marshall(announce)

		unpiMock.InjectOutgoing(Frame{
			MessageType: AREQ,
			Subsystem:   ZDO,
			CommandID:   ZdoMGMTLQIRspID,
			Payload:     data,
		})

		time.Sleep(10 * time.Millisecond)

		_, found := zstack.devices[zigbee.IEEEAddress(0)]
		assert.False(t, found)
	})

	t.Run("neighbours are removed from device if LQI does not return them", func(t *testing.T) {
		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock)
		defer unpiMock.Stop()
		defer unpiMock.AssertCalls(t)

		unpiMock.On(SREQ, ZDO, ZdoMGMTLQIReqID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   ZDO,
			CommandID:   ZdoMGMTLQIReqReplyID,
			Payload:     []byte{0x00},
		}).UnlimitedTimes()

		zstack.startNetworkManager()
		defer zstack.stopNetworkManager()

		time.Sleep(10 * time.Millisecond)

		zstack.devices[zigbee.IEEEAddress(0)].Neighbours[zigbee.IEEEAddress(0x1000)] = &DeviceNeighbour{}

		announce := ZdoMGMTLQIRsp{
			SourceAddress:         0,
			Status:                0,
			NeighbourTableEntries: 1,
			StartIndex:            0,
			Neighbors:             []ZdoMGMTLQINeighbour{},
		}

		data, _ := bytecodec.Marshall(announce)

		unpiMock.InjectOutgoing(Frame{
			MessageType: AREQ,
			Subsystem:   ZDO,
			CommandID:   ZdoMGMTLQIRspID,
			Payload:     data,
		})

		time.Sleep(10 * time.Millisecond)

		requestingDevice, _ := zstack.devices[zigbee.IEEEAddress(0)]
		_, found := requestingDevice.Neighbours[zigbee.IEEEAddress(0x1000)]
		assert.False(t, found)
	})
}
