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
	z.addOrUpdateDevice(z.NetworkProperties.IEEEAddress, z.NetAddr(z.NetworkProperties.NetworkAddress), z.Role(RoleCoordinator))

	immediateStart := make(chan bool, 1)
	defer close(immediateStart)
	immediateStart <- true

	_, cancel := z.subscriber.Subscribe(&ZdoMGMTLQIRsp{}, z.receiveLQIUpdate)
	defer cancel()

	_, cancel = z.subscriber.Subscribe(&ZdoEndDeviceAnnceInd{}, z.receiveEndDeviceAnnouncement)
	defer cancel()

	_, cancel = z.subscriber.Subscribe(&ZdoLeaveInd{}, z.receiveLeaveAnnouncement)
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
				role := RoleEndDevice

				if e.Capabilities&0x02 == 0x02 {
					role = RoleRouter
				}

				z.addOrUpdateDevice(e.IEEEAddress, z.NetAddr(e.NetworkAddress), z.Role(role))
				z.events <- DeviceJoinEvent{
					NetworkAddress: e.NetworkAddress,
					IEEEAddress:    e.IEEEAddress,
				}
			case ZdoLeaveInd:
				z.removeDevice(e.IEEEAddress)
				z.events <- DeviceLeaveEvent{
					NetworkAddress: e.SourceAddress,
					IEEEAddress:    e.IEEEAddress,
				}
			default:
				fmt.Printf("received unknown %+v", reflect.TypeOf(ue))
			}
		}
	}
}

func (z *ZStack) pollRoutersForNetworkStatus() {
	for _, device := range z.devices {
		if device.Role == RoleCoordinator || device.Role == RoleRouter {
			go z.pollDeviceForNetworkStatus(*device)
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

	currentTime := time.Now()

	sourceDevice, sourceFound := z.getDevice(lqi.SourceAddress)

	for _, neighbour := range lqi.Neighbors {
		if neighbour.ExtendedPANID != z.NetworkProperties.ExtendedPANID || neighbour.IEEEAddress == zigbee.EmptyIEEEAddress {
			continue
		}

		role := RoleUnknown

		switch neighbour.Status & 0x03 {
		case 0x00:
			role = RoleCoordinator
		case 0x01:
			role = RoleRouter
		case 0x02:
			role = RoleEndDevice
		}

		z.addOrUpdateDevice(neighbour.IEEEAddress, z.NetAddr(neighbour.NetworkAddress), z.Role(role))

		if sourceFound {
			dn, found := sourceDevice.Neighbours[neighbour.IEEEAddress]

			if !found {
				dn = &DeviceNeighbour{}
				sourceDevice.Neighbours[neighbour.IEEEAddress] = dn
			}

			dn.LQI = neighbour.LQI
			dn.Relationship = DeviceRelationship((neighbour.Status & 0x70) >> 4)
			dn.LastObserved = currentTime
		}
	}

	if sourceFound {
		var oldNeighbours []zigbee.IEEEAddress

		for ieee, neighbour := range sourceDevice.Neighbours {
			if neighbour.LastObserved.Before(currentTime) {
				oldNeighbours = append(oldNeighbours, ieee)
			}
		}

		for _, ieee := range oldNeighbours {
			delete(sourceDevice.Neighbours, ieee)
		}
	}
}

func (z *ZStack) receiveLQIUpdate(u func(interface{}) error) {
	msg := ZdoMGMTLQIRsp{}
	if err := u(&msg); err == nil {
		z.networkManagerIncoming <- msg
	}
}

func (z *ZStack) receiveEndDeviceAnnouncement(u func(interface{}) error) {
	msg := ZdoEndDeviceAnnceInd{}
	var err error
	if err = u(&msg); err == nil {
		z.networkManagerIncoming <- msg
	}
}

func (z *ZStack) receiveLeaveAnnouncement(u func(interface{}) error) {
	msg := ZdoLeaveInd{}
	var err error
	if err = u(&msg); err == nil {
		z.networkManagerIncoming <- msg
	}
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