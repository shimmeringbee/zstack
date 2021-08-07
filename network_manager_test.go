package zstack

import (
	"context"
	"github.com/shimmeringbee/bytecodec"
	. "github.com/shimmeringbee/unpi"
	unpiTest "github.com/shimmeringbee/unpi/testing"
	"github.com/shimmeringbee/zigbee"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/semaphore"
	"testing"
	"time"
)

func Test_NetworkManager(t *testing.T) {
	t.Run("issues a lqi poll request only for coordinators or routers", func(t *testing.T) {
		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock, NewNodeTable())
		zstack.sem = semaphore.NewWeighted(8)
		defer unpiMock.Stop()
		defer zstack.Stop()

		unpiMock.On(SREQ, ZDO, ZdoMGMTLQIReqID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   ZDO,
			CommandID:   ZdoMGMTLQIReqReplyID,
			Payload:     []byte{0x00},
		}).Times(2)

		zstack.nodeTable.addOrUpdate(zigbee.IEEEAddress(1), zigbee.NetworkAddress(1), logicalType(zigbee.Router))
		zstack.nodeTable.addOrUpdate(zigbee.IEEEAddress(2), zigbee.NetworkAddress(2), logicalType(zigbee.Unknown))

		zstack.startNetworkManager()
		defer zstack.stopNetworkManager()

		time.Sleep(10 * time.Millisecond)

		unpiMock.AssertCalls(t)
	})

	t.Run("the coordinator is added to the node list as a coordinator", func(t *testing.T) {
		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock, NewNodeTable())
		zstack.sem = semaphore.NewWeighted(8)
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

		node, found := zstack.nodeTable.getByIEEE(expectedIEEE)

		assert.True(t, found)
		assert.Equal(t, expectedAddress, node.NetworkAddress)
		assert.Equal(t, zigbee.Coordinator, node.LogicalType)

		unpiMock.AssertCalls(t)
	})

	t.Run("a node is added to the node table when an ZdoIEEEAddrRsp messages are received", func(t *testing.T) {
		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock, NewNodeTable())
		zstack.sem = semaphore.NewWeighted(8)
		defer unpiMock.Stop()
		defer unpiMock.AssertCalls(t)

		zstack.startNetworkManager()
		defer zstack.stopNetworkManager()

		unpiMock.On(SREQ, ZDO, ZdoMGMTLQIReqID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   ZDO,
			CommandID:   ZdoMGMTLQIReqReplyID,
			Payload:     []byte{0x00},
		})

		time.Sleep(10 * time.Millisecond)

		unpiMock.InjectOutgoing(Frame{
			MessageType: AREQ,
			Subsystem:   ZDO,
			CommandID:   ZdoIEEEAddrRspID,
			Payload:     []byte{0x00, 0x88, 0x77, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11, 0x00, 0x40, 0x00, 0x00},
		})

		time.Sleep(10 * time.Millisecond)

		node, found := zstack.nodeTable.getByIEEE(0x1122334455667788)

		assert.True(t, found)
		assert.Equal(t, zigbee.IEEEAddress(0x1122334455667788), node.IEEEAddress)
		assert.Equal(t, zigbee.NetworkAddress(0x4000), node.NetworkAddress)
	})

	t.Run("a node is added to the node table when an ZdoNWKAddrRsp messages are received", func(t *testing.T) {
		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock, NewNodeTable())
		zstack.sem = semaphore.NewWeighted(8)
		defer unpiMock.Stop()
		defer unpiMock.AssertCalls(t)

		zstack.startNetworkManager()
		defer zstack.stopNetworkManager()

		unpiMock.On(SREQ, ZDO, ZdoMGMTLQIReqID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   ZDO,
			CommandID:   ZdoMGMTLQIReqReplyID,
			Payload:     []byte{0x00},
		})

		time.Sleep(10 * time.Millisecond)

		unpiMock.InjectOutgoing(Frame{
			MessageType: AREQ,
			Subsystem:   ZDO,
			CommandID:   ZdoNWKAddrRspID,
			Payload:     []byte{0x00, 0x88, 0x77, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11, 0x00, 0x40, 0x00, 0x00},
		})

		time.Sleep(10 * time.Millisecond)

		node, found := zstack.nodeTable.getByIEEE(0x1122334455667788)

		assert.True(t, found)
		assert.Equal(t, zigbee.IEEEAddress(0x1122334455667788), node.IEEEAddress)
		assert.Equal(t, zigbee.NetworkAddress(0x4000), node.NetworkAddress)
	})

	t.Run("emits NodeJoinEvent event when node join announcement received", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock, NewNodeTable())
		zstack.sem = semaphore.NewWeighted(8)
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
			Capabilities: ZdoEndDeviceAnnceIndCapabilities{
				AltPANController:   false,
				Router:             true,
				PowerSource:        false,
				ReceiveOnIdle:      false,
				Reserved:           0,
				SecurityCapability: false,
				AddressAllocated:   false,
			},
		}

		data, _ := bytecodec.Marshal(announce)

		unpiMock.InjectOutgoing(Frame{
			MessageType: AREQ,
			Subsystem:   ZDO,
			CommandID:   ZdoEndDeviceAnnceIndID,
			Payload:     data,
		})

		// Throw away the NodeUpdateEvent.
		zstack.ReadEvent(ctx)

		event, err := zstack.ReadEvent(ctx)
		assert.NoError(t, err)

		nodeJoin, ok := event.(zigbee.NodeJoinEvent)

		assert.True(t, ok)

		assert.Equal(t, announce.NetworkAddress, nodeJoin.NetworkAddress)
		assert.Equal(t, announce.IEEEAddress, nodeJoin.IEEEAddress)

		node, found := zstack.nodeTable.getByIEEE(announce.IEEEAddress)

		assert.True(t, found)
		assert.Equal(t, announce.IEEEAddress, node.IEEEAddress)
		assert.Equal(t, announce.NetworkAddress, node.NetworkAddress)
		assert.Equal(t, zigbee.Router, node.LogicalType)
	})

	t.Run("emits NodeLeaveEvent event when node leave announcement received", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock, NewNodeTable())
		zstack.sem = semaphore.NewWeighted(8)
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

		zstack.nodeTable.addOrUpdate(zigbee.IEEEAddress(0x0102030405060708), zigbee.NetworkAddress(0x2000))

		data, _ := bytecodec.Marshal(announce)

		time.Sleep(10 * time.Millisecond)

		unpiMock.InjectOutgoing(Frame{
			MessageType: AREQ,
			Subsystem:   ZDO,
			CommandID:   ZdoLeaveIndID,
			Payload:     data,
		})

		event, err := zstack.ReadEvent(ctx)
		assert.NoError(t, err)

		nodeLeave, ok := event.(zigbee.NodeLeaveEvent)

		assert.True(t, ok)

		assert.Equal(t, announce.SourceAddress, nodeLeave.NetworkAddress)
		assert.Equal(t, announce.IEEEAddress, nodeLeave.IEEEAddress)

		_, found := zstack.nodeTable.getByIEEE(announce.IEEEAddress)
		assert.False(t, found)
	})

	t.Run("a new router will be queried for network state", func(t *testing.T) {
		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock, NewNodeTable())
		zstack.sem = semaphore.NewWeighted(8)
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
			Capabilities: ZdoEndDeviceAnnceIndCapabilities{
				AltPANController:   false,
				Router:             true,
				PowerSource:        false,
				ReceiveOnIdle:      false,
				Reserved:           0,
				SecurityCapability: false,
				AddressAllocated:   false,
			},
		}

		data, _ := bytecodec.Marshal(announce)

		unpiMock.InjectOutgoing(Frame{
			MessageType: AREQ,
			Subsystem:   ZDO,
			CommandID:   ZdoEndDeviceAnnceIndID,
			Payload:     data,
		})

		time.Sleep(20 * time.Millisecond)

		assert.Equal(t, 2, len(c.CapturedCalls))

		frame := c.CapturedCalls[1]

		lqiReq := ZdoMGMTLQIReq{}
		_ = bytecodec.Unmarshal(frame.Frame.Payload, &lqiReq)

		assert.Equal(t, zigbee.NetworkAddress(0x2000), lqiReq.DestinationAddress)
	})

	t.Run("nodes in lqi query are added to network manager", func(t *testing.T) {
		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock, NewNodeTable())
		zstack.sem = semaphore.NewWeighted(8)
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
			Neighbors: []ZdoMGMTLQINeighbour{
				{
					ExtendedPANID:  zstack.NetworkProperties.ExtendedPANID,
					IEEEAddress:    zigbee.IEEEAddress(0x1000),
					NetworkAddress: zigbee.NetworkAddress(0x2000),
					Status: ZdoMGMTLQINeighbourStatus{
						Reserved:     0,
						Relationship: zigbee.RelationshipChild,
						RxOnWhenIdle: 0,
						DeviceType:   zigbee.Router,
					},
					PermitJoining: false,
					Depth:         1,
					LQI:           67,
				},
			},
		}

		data, _ := bytecodec.Marshal(announce)

		unpiMock.InjectOutgoing(Frame{
			MessageType: AREQ,
			Subsystem:   ZDO,
			CommandID:   ZdoMGMTLQIRspID,
			Payload:     data,
		})

		time.Sleep(10 * time.Millisecond)

		node, found := zstack.nodeTable.getByIEEE(zigbee.IEEEAddress(0x1000))
		assert.True(t, found)

		assert.Equal(t, zigbee.NetworkAddress(0x2000), node.NetworkAddress)
		assert.Equal(t, zigbee.Router, node.LogicalType)
		assert.Equal(t, uint8(0x43), node.LQI)
		assert.Equal(t, uint8(0x01), node.Depth)
	})

	t.Run("nodes in lqi query are added to network manager", func(t *testing.T) {
		unpiMock := unpiTest.NewMockAdapter()
		nt := NewNodeTable()
		zstack := New(unpiMock, nt)
		zstack.sem = semaphore.NewWeighted(8)
		defer unpiMock.Stop()
		defer unpiMock.AssertCalls(t)

		nt.addOrUpdate(zigbee.GenerateLocalAdministeredIEEEAddress(), 0x1122)

		lqiReqOn := unpiMock.On(SREQ, ZDO, ZdoMGMTLQIReqID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   ZDO,
			CommandID:   ZdoMGMTLQIReqReplyID,
			Payload:     []byte{0x00},
		}).Times(1)

		lqiRespWithMore := ZdoMGMTLQIRsp{
			SourceAddress:         0x1122,
			Status:                0,
			NeighbourTableEntries: 2,
			StartIndex:            0,
			Neighbors: []ZdoMGMTLQINeighbour{
				{
					ExtendedPANID:  zstack.NetworkProperties.ExtendedPANID,
					IEEEAddress:    zigbee.IEEEAddress(0x1000),
					NetworkAddress: zigbee.NetworkAddress(0x2000),
					Status: ZdoMGMTLQINeighbourStatus{
						Reserved:     0,
						Relationship: zigbee.RelationshipChild,
						RxOnWhenIdle: 0,
						DeviceType:   zigbee.Router,
					},
					PermitJoining: false,
					Depth:         1,
					LQI:           67,
				},
			},
		}

		zstack.processLQITable(lqiRespWithMore)

		assert.Equal(t, []byte{0x22, 0x11, 0x01}, lqiReqOn.CapturedCalls[0].Frame.Payload)
	})

	t.Run("nodes in lqi query are not added if Ext PANID does not match", func(t *testing.T) {
		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock, NewNodeTable())
		zstack.sem = semaphore.NewWeighted(8)
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
			Neighbors: []ZdoMGMTLQINeighbour{
				{
					ExtendedPANID:  0xfffffff,
					IEEEAddress:    zigbee.IEEEAddress(0x2000),
					NetworkAddress: zigbee.NetworkAddress(0x4000),
					Status: ZdoMGMTLQINeighbourStatus{
						Reserved:     0,
						Relationship: zigbee.RelationshipParent,
						RxOnWhenIdle: 0,
						DeviceType:   zigbee.Router,
					},
					PermitJoining: false,
					Depth:         0,
					LQI:           67,
				},
			},
		}

		data, _ := bytecodec.Marshal(announce)

		unpiMock.InjectOutgoing(Frame{
			MessageType: AREQ,
			Subsystem:   ZDO,
			CommandID:   ZdoMGMTLQIRspID,
			Payload:     data,
		})

		time.Sleep(10 * time.Millisecond)

		_, found := zstack.nodeTable.getByIEEE(zigbee.IEEEAddress(0x2000))
		assert.False(t, found)
	})

	t.Run("nodes in lqi query are not added if it has an invalid IEEE address", func(t *testing.T) {
		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock, NewNodeTable())
		zstack.sem = semaphore.NewWeighted(8)
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
			Neighbors: []ZdoMGMTLQINeighbour{
				{
					ExtendedPANID:  zstack.NetworkProperties.ExtendedPANID,
					IEEEAddress:    zigbee.IEEEAddress(0),
					NetworkAddress: zigbee.NetworkAddress(0x4000),
					Status: ZdoMGMTLQINeighbourStatus{
						Reserved:     0,
						Relationship: zigbee.RelationshipParent,
						RxOnWhenIdle: 0,
						DeviceType:   zigbee.Router,
					},
					PermitJoining: false,
					Depth:         0,
					LQI:           67,
				},
			},
		}

		data, _ := bytecodec.Marshal(announce)

		unpiMock.InjectOutgoing(Frame{
			MessageType: AREQ,
			Subsystem:   ZDO,
			CommandID:   ZdoMGMTLQIRspID,
			Payload:     data,
		})

		time.Sleep(10 * time.Millisecond)

		_, found := zstack.nodeTable.getByIEEE(zigbee.IEEEAddress(0))
		assert.False(t, found)
	})

	t.Run("updates to the node table sends a node update event", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock, NewNodeTable())
		zstack.sem = semaphore.NewWeighted(8)
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

		go func() {
			time.Sleep(10 * time.Millisecond)
			zstack.nodeTable.addOrUpdate(zigbee.IEEEAddress(0x01), zigbee.NetworkAddress(0x02))
		}()

		event, err := zstack.ReadEvent(ctx)
		assert.NoError(t, err)

		nodeUpdateEvent, ok := event.(zigbee.NodeUpdateEvent)

		assert.True(t, ok)

		assert.Equal(t, zigbee.IEEEAddress(0x01), nodeUpdateEvent.Node.IEEEAddress)
		assert.Equal(t, zigbee.NetworkAddress(0x02), nodeUpdateEvent.Node.NetworkAddress)
	})
}
