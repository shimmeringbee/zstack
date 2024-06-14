package zstack

import (
	"context"
	"github.com/shimmeringbee/logwrap"
	"github.com/shimmeringbee/zigbee"
	"reflect"
	"time"
)

const defaultPollingInterval = 30

func (z *ZStack) startNetworkManager() {
	go z.networkManager()
}

func (z *ZStack) stopNetworkManager() {
	z.networkManagerStop <- true
}

func (z *ZStack) networkManager() {
	z.nodeTable.addOrUpdate(z.NetworkProperties.IEEEAddress, z.NetworkProperties.NetworkAddress, logicalType(zigbee.Coordinator))

	immediateStart := make(chan bool, 1)
	defer close(immediateStart)
	immediateStart <- true

	_, cancel := z.subscriber.Subscribe(&ZdoMGMTLQIRsp{}, z.receiveLQIUpdate)
	defer cancel()

	_, cancel = z.subscriber.Subscribe(&ZdoEndDeviceAnnceInd{}, z.receiveEndDeviceAnnouncement)
	defer cancel()

	_, cancel = z.subscriber.Subscribe(&ZdoLeaveInd{}, z.receiveLeaveAnnouncement)
	defer cancel()

	_, cancel = z.subscriber.Subscribe(&ZdoIEEEAddrRsp{}, z.receiveIEEEAddrRsp)
	defer cancel()

	_, cancel = z.subscriber.Subscribe(&ZdoNWKAddrRsp{}, z.receiveNWKAddrRsp)
	defer cancel()

	z.nodeTable.registerCallback(z.nodeTableUpdate)

	for {
		select {
		case <-immediateStart:
			z.pollRoutersForNetworkStatus()
		case <-time.After(defaultPollingInterval * time.Second):
			z.pollRoutersForNetworkStatus()
		case <-z.networkManagerStop:
			return
		case ue := <-z.networkManagerIncoming:
			switch e := ue.(type) {
			case ZdoMGMTLQIRsp:
				z.processLQITable(e)
			case ZdoEndDeviceAnnceInd:
				z.newNode(e)
			case ZdoLeaveInd:
				z.removeNode(e.IEEEAddress)
			case ZdoIEEEAddrRsp:
				if e.WasSuccessful() {
					z.nodeTable.addOrUpdate(e.IEEEAddress, e.NetworkAddress, updateDiscovered())
				}
			case ZdoNWKAddrRsp:
				if e.WasSuccessful() {
					z.nodeTable.addOrUpdate(e.IEEEAddress, e.NetworkAddress, updateDiscovered())
				}
			default:
				z.logger.LogWarn(context.Background(), "Received unknown message type from unpi.", logwrap.Datum("Type", reflect.TypeOf(ue)))
			}
		}
	}
}

func (z *ZStack) newNode(e ZdoEndDeviceAnnceInd) {
	deviceLogicalType := zigbee.EndDevice

	if e.Capabilities.Router {
		deviceLogicalType = zigbee.Router
	}

	z.nodeTable.addOrUpdate(e.IEEEAddress, e.NetworkAddress, logicalType(deviceLogicalType), updateDiscovered(), updateReceived())
	node, _ := z.nodeTable.getByIEEE(e.IEEEAddress)

	z.sendEvent(zigbee.NodeJoinEvent{
		Node: node,
	})

	if deviceLogicalType == zigbee.Router {
		node, _ := z.nodeTable.getByIEEE(e.IEEEAddress)
		go z.pollNodeForNetworkStatus(node)
	}
}

func (z *ZStack) removeNode(ieee zigbee.IEEEAddress) bool {
	node, found := z.nodeTable.getByIEEE(ieee)
	z.nodeTable.remove(ieee)

	if found {
		z.sendEvent(zigbee.NodeLeaveEvent{
			Node: node,
		})
	}

	return found
}

func (z *ZStack) pollRoutersForNetworkStatus() {
	for _, node := range z.nodeTable.nodes() {
		if node.LogicalType == zigbee.Coordinator || node.LogicalType == zigbee.Router {
			go z.pollNodeForNetworkStatus(node)
		}
	}
}

func (z *ZStack) pollNodeForNetworkStatus(node zigbee.Node) {
	z.logger.LogDebug(context.Background(), "Polling device for network status.", logwrap.Datum("IEEEAddress", node.IEEEAddress.String()), logwrap.Datum("NetworkAddress", node.NetworkAddress))
	z.requestLQITable(node, 0)
}

func (z *ZStack) requestLQITable(node zigbee.Node, startIndex uint8) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultZStackTimeout)
	defer cancel()

	if err := z.sem.Acquire(ctx, 1); err != nil {
		z.logger.LogError(ctx, "Failed to request LQI table, failed to acquire semaphore ", logwrap.Datum("IEEEAddress", node.IEEEAddress.String()), logwrap.Datum("NetworkAddress", node.NetworkAddress), logwrap.Err(err))
		return
	}
	defer z.sem.Release(1)

	resp := ZdoMGMTLQIReqReply{}
	z.logger.LogDebug(ctx, "Requesting LQI table from device.", logwrap.Datum("IEEEAddress", node.IEEEAddress.String()), logwrap.Datum("NetworkAddress", node.NetworkAddress), logwrap.Datum("StartIndex", startIndex))
	if err := z.requestResponder.RequestResponse(ctx, ZdoMGMTLQIReq{DestinationAddress: node.NetworkAddress, StartIndex: startIndex}, &resp); err != nil {
		z.logger.LogError(ctx, "Failed to request LQI table.", logwrap.Datum("IEEEAddress", node.IEEEAddress.String()), logwrap.Datum("NetworkAddress", node.NetworkAddress), logwrap.Err(err))
	} else if resp.Status != ZSuccess {
		z.logger.LogError(ctx, "Failed to request LQI table, adapter returned error code.", logwrap.Datum("IEEEAddress", node.IEEEAddress.String()), logwrap.Datum("NetworkAddress", node.NetworkAddress), logwrap.Datum("Status", resp.Status))
	}
}

func (z *ZStack) processLQITable(lqiResp ZdoMGMTLQIRsp) {
	if lqiResp.Status != ZSuccess {
		z.logger.LogError(context.Background(), "LQI table response received, but as not success.", logwrap.Datum("NetworkAddress", lqiResp.SourceAddress), logwrap.Datum("Status", lqiResp.Status))
		return
	}

	z.logger.LogDebug(context.Background(), "LQI table response received.", logwrap.Datum("NetworkAddress", lqiResp.SourceAddress), logwrap.Datum("Status", lqiResp.Status), logwrap.Datum("StartIndex", lqiResp.StartIndex), logwrap.Datum("IncludedCount", len(lqiResp.Neighbors)), logwrap.Datum("NeighbourCount", lqiResp.NeighbourTableEntries))

	for _, neighbour := range lqiResp.Neighbors {
		if neighbour.ExtendedPANID != z.NetworkProperties.ExtendedPANID ||
			neighbour.IEEEAddress == zigbee.EmptyIEEEAddress {
			continue
		}

		z.nodeTable.addOrUpdate(neighbour.IEEEAddress, neighbour.NetworkAddress, logicalType(neighbour.Status.DeviceType), updateDiscovered())

		if neighbour.Status.Relationship == zigbee.RelationshipChild {
			z.nodeTable.update(neighbour.IEEEAddress, lqi(neighbour.LQI), depth(neighbour.Depth))
		}
	}

	nextIndex := uint8(int(lqiResp.StartIndex) + len(lqiResp.Neighbors))

	if nextIndex < lqiResp.NeighbourTableEntries {
		node, found := z.nodeTable.getByNetwork(lqiResp.SourceAddress)

		if found {
			z.logger.LogDebug(context.Background(), "LQI table response requires pagination.", logwrap.Datum("NetworkAddress", lqiResp.SourceAddress), logwrap.Datum("Status", lqiResp.Status), logwrap.Datum("StartIndex", lqiResp.StartIndex), logwrap.Datum("IncludedCount", len(lqiResp.Neighbors)), logwrap.Datum("NeighbourCount", lqiResp.NeighbourTableEntries))
			z.requestLQITable(node, nextIndex)
		}
	}
}

func (z *ZStack) receiveLQIUpdate(v interface{}) {
	msg := v.(*ZdoMGMTLQIRsp)
	z.networkManagerIncoming <- *msg
}

func (z *ZStack) receiveEndDeviceAnnouncement(v interface{}) {
	msg := v.(*ZdoEndDeviceAnnceInd)
	z.networkManagerIncoming <- *msg
}

func (z *ZStack) receiveLeaveAnnouncement(v interface{}) {
	msg := v.(*ZdoLeaveInd)
	z.networkManagerIncoming <- *msg
}

func (z *ZStack) receiveIEEEAddrRsp(v interface{}) {
	msg := v.(*ZdoIEEEAddrRsp)
	z.networkManagerIncoming <- *msg
}

func (z *ZStack) receiveNWKAddrRsp(v interface{}) {
	msg := v.(*ZdoNWKAddrRsp)
	z.networkManagerIncoming <- *msg
}

func (z *ZStack) nodeTableUpdate(node zigbee.Node) {
	z.sendEvent(zigbee.NodeUpdateEvent{
		Node: node,
	})
}

type ZdoMGMTLQIReq struct {
	DestinationAddress zigbee.NetworkAddress
	StartIndex         uint8
}

const ZdoMGMTLQIReqID uint8 = 0x31

type ZdoMGMTLQIReqReply GenericZStackStatus

const ZdoMGMTLQIReqReplyID uint8 = 0x31

type ZdoMGMTLQINeighbourStatus struct {
	Reserved     uint8               `bcfieldwidth:"1"`
	Relationship zigbee.Relationship `bcfieldwidth:"3"`
	RxOnWhenIdle uint8               `bcfieldwidth:"2"`
	DeviceType   zigbee.LogicalType  `bcfieldwidth:"2"`
}

type ZdoMGMTLQINeighbour struct {
	ExtendedPANID  zigbee.ExtendedPANID
	IEEEAddress    zigbee.IEEEAddress
	NetworkAddress zigbee.NetworkAddress
	Status         ZdoMGMTLQINeighbourStatus
	PermitJoining  bool
	Depth          uint8
	LQI            uint8
}

type ZdoMGMTLQIRsp struct {
	SourceAddress         zigbee.NetworkAddress
	Status                ZStackStatus
	NeighbourTableEntries uint8
	StartIndex            uint8
	Neighbors             []ZdoMGMTLQINeighbour `bcsliceprefix:"8"`
}

const ZdoMGMTLQIRspID uint8 = 0xb1
