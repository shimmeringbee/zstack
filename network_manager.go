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
	z.deviceTable.AddOrUpdate(z.NetworkProperties.IEEEAddress, z.NetworkProperties.NetworkAddress, LogicalType(zigbee.Coordinator))

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
				logicalType := zigbee.EndDevice

				if e.Capabilities&0x02 == 0x02 {
					logicalType = zigbee.Router
				}

				z.deviceTable.AddOrUpdate(e.IEEEAddress, e.NetworkAddress, LogicalType(logicalType))
				z.sendEvent(zigbee.DeviceJoinEvent{
					NetworkAddress: e.NetworkAddress,
					IEEEAddress:    e.IEEEAddress,
				})

				if logicalType == zigbee.Router {
					device, _ := z.deviceTable.GetByIEEE(e.IEEEAddress)
					go z.pollDeviceForNetworkStatus(device)
				}
			case ZdoLeaveInd:
				z.deviceTable.Remove(e.IEEEAddress)
				z.sendEvent(zigbee.DeviceLeaveEvent{
					NetworkAddress: e.SourceAddress,
					IEEEAddress:    e.IEEEAddress,
				})
			case ZdoIEEEAddrRsp:
				if e.WasSuccessful() {
					z.deviceTable.AddOrUpdate(e.IEEEAddress, e.NetworkAddress)
				}
			case ZdoNWKAddrRsp:
				if e.WasSuccessful() {
					z.deviceTable.AddOrUpdate(e.IEEEAddress, e.NetworkAddress)
				}
			default:
				fmt.Printf("received unknown %+v", reflect.TypeOf(ue))
			}
		}
	}
}

func (z *ZStack) pollRoutersForNetworkStatus() {
	for _, device := range z.deviceTable.GetAllDevices() {
		if device.LogicalType == zigbee.Coordinator || device.LogicalType == zigbee.Router {
			go z.pollDeviceForNetworkStatus(device)
		}
	}
}

func (z *ZStack) pollDeviceForNetworkStatus(device Device) {
	log.Printf("polling %v (%d) for network status\n", device.IEEEAddress, device.NetworkAddress)
	z.requestLQITable(device)
}

func (z *ZStack) requestLQITable(device Device) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultZStackTimeout)
	defer cancel()

	resp := ZdoMGMTLQIReqReply{}
	if err := z.requestResponder.RequestResponse(ctx, ZdoMGMTLQIReq{DestinationAddress: device.NetworkAddress, StartIndex: 0}, &resp); err != nil {
		log.Printf("failed to request lqi tables: %v\n", err)
	}

	if resp.Status != ZSuccess {
		log.Printf("failed to request lqi tables: from the adapter\n")
	}
}

func (z *ZStack) processLQITable(lqi ZdoMGMTLQIRsp) {
	if lqi.Status != ZSuccess {
		log.Printf("failed lqi response from %+v\n", lqi.SourceAddress)
		return
	}

	for _, neighbour := range lqi.Neighbors {
		if neighbour.ExtendedPANID != z.NetworkProperties.ExtendedPANID ||
			neighbour.IEEEAddress == zigbee.EmptyIEEEAddress {
			continue
		}

		logicalType := zigbee.LogicalType(neighbour.Status & 0x03)
		relationship := zigbee.Relationship((neighbour.Status >> 4) & 0x07)

		z.deviceTable.AddOrUpdate(neighbour.IEEEAddress, neighbour.NetworkAddress, LogicalType(logicalType))

		if relationship == zigbee.RelationshipChild {
			z.deviceTable.Update(neighbour.IEEEAddress, LQI(neighbour.LQI), Depth(neighbour.Depth))
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

type ZdoMGMTLQIReq struct {
	DestinationAddress zigbee.NetworkAddress
	StartIndex         uint8
}

const ZdoMGMTLQIReqID uint8 = 0x31

type ZdoMGMTLQIReqReply GenericZStackStatus

const ZdoMGMTLQIReqReplyID uint8 = 0x31

type ZdoMGMTLQINeighbour struct {
	ExtendedPANID  zigbee.ExtendedPANID
	IEEEAddress    zigbee.IEEEAddress
	NetworkAddress zigbee.NetworkAddress
	Status         uint8
	PermitJoining  uint8
	Depth          uint8
	LQI            uint8
}

type ZdoMGMTLQIRsp struct {
	SourceAddress         zigbee.NetworkAddress
	Status                ZStackStatus
	NeighbourTableEntries uint8
	StartIndex            uint8
	Neighbors             []ZdoMGMTLQINeighbour `bclength:"8"`
}

const ZdoMGMTLQIRspID uint8 = 0xb1
