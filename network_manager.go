package zstack

import (
	"context"
	"fmt"
	"github.com/shimmeringbee/zigbee"
	"log"
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
				deviceLogicalType := zigbee.EndDevice

				if e.Capabilities.Router {
					deviceLogicalType = zigbee.Router
				}

				z.nodeTable.addOrUpdate(e.IEEEAddress, e.NetworkAddress, logicalType(deviceLogicalType), updateDiscovered, updateReceived)
				node, _ := z.nodeTable.getByIEEE(e.IEEEAddress)

				z.sendEvent(zigbee.NodeJoinEvent{
					Node: node,
				})

				if deviceLogicalType == zigbee.Router {
					node, _ := z.nodeTable.getByIEEE(e.IEEEAddress)
					go z.pollNodeForNetworkStatus(node)
				}
			case ZdoLeaveInd:
				node, found := z.nodeTable.getByIEEE(e.IEEEAddress)
				z.nodeTable.remove(e.IEEEAddress)

				if found {
					z.sendEvent(zigbee.NodeLeaveEvent{
						Node: node,
					})
				}
			case ZdoIEEEAddrRsp:
				if e.WasSuccessful() {
					z.nodeTable.addOrUpdate(e.IEEEAddress, e.NetworkAddress, updateDiscovered)
				}
			case ZdoNWKAddrRsp:
				if e.WasSuccessful() {
					z.nodeTable.addOrUpdate(e.IEEEAddress, e.NetworkAddress, updateDiscovered)
				}
			default:
				fmt.Printf("received unknown %+v", reflect.TypeOf(ue))
			}
		}
	}
}

func (z *ZStack) pollRoutersForNetworkStatus() {
	for _, node := range z.nodeTable.Nodes() {
		if node.LogicalType == zigbee.Coordinator || node.LogicalType == zigbee.Router {
			go z.pollNodeForNetworkStatus(node)
		}
	}
}

func (z *ZStack) pollNodeForNetworkStatus(node zigbee.Node) {
	log.Printf("polling %v (%d) for network status\n", node.IEEEAddress, node.NetworkAddress)
	z.requestLQITable(node)
}

func (z *ZStack) requestLQITable(node zigbee.Node) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultZStackTimeout)
	defer cancel()

	resp := ZdoMGMTLQIReqReply{}
	if err := z.requestResponder.RequestResponse(ctx, ZdoMGMTLQIReq{DestinationAddress: node.NetworkAddress, StartIndex: 0}, &resp); err != nil {
		log.Printf("failed to request lqi tables: %v\n", err)
	}

	if resp.Status != ZSuccess {
		log.Printf("failed to request lqi tables: from the adapter\n")
	}
}

func (z *ZStack) processLQITable(lqiResp ZdoMGMTLQIRsp) {
	if lqiResp.Status != ZSuccess {
		log.Printf("failed lqi response from %+v\n", lqiResp.SourceAddress)
		return
	}

	for _, neighbour := range lqiResp.Neighbors {
		if neighbour.ExtendedPANID != z.NetworkProperties.ExtendedPANID ||
			neighbour.IEEEAddress == zigbee.EmptyIEEEAddress {
			continue
		}

		z.nodeTable.addOrUpdate(neighbour.IEEEAddress, neighbour.NetworkAddress, logicalType(neighbour.Status.DeviceType), updateDiscovered)

		if neighbour.Status.Relationship == zigbee.RelationshipChild {
			z.nodeTable.update(neighbour.IEEEAddress, lqi(neighbour.LQI), depth(neighbour.Depth))
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
