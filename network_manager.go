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
	z.addOrUpdateDevice(z.NetworkProperties.IEEEAddress, z.NetworkProperties.NetworkAddress, Role(RoleCoordinator))

	immediateStart := make(chan bool, 1)
	defer close(immediateStart)
	immediateStart <- true

	_, cancel := z.subscriber.Subscribe(&ZdoMGMTLQIResp{}, z.receiveLQIUpdate)
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
			case ZdoMGMTLQIResp:
				z.processLQITable(e)
			case ZdoEndDeviceAnnceInd:
				role := RoleEndDevice

				if e.Capabilities&0x02 == 0x02 {
					role = RoleRouter
				}

				z.addOrUpdateDevice(e.IEEEAddress, e.NetworkAddress, Role(role))
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

	resp := ZdoMGMTLQIReqResp{}
	if err := z.requestResponder.RequestResponse(ctx, ZdoMGMTLQIReq{DestinationAddress: device.NetworkAddress, StartIndex: 0}, &resp); err != nil {
		log.Printf("failed to request lqi tables: %v\n", err)
	}

	if resp.Status != ZSuccess {
		log.Printf("failed to request lqi tables: from adapter\n")
	}
}

func (z *ZStack) processLQITable(lqi ZdoMGMTLQIResp) {

}

func (z *ZStack) receiveLQIUpdate(u func(interface{}) error) {
	msg := ZdoMGMTLQIResp{}
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

type DeviceFact func(*Device)

func Role(role DeviceRole) DeviceFact {
	return func(device *Device) {
		device.Role = role
	}
}

func (z *ZStack) addOrUpdateDevice(ieee zigbee.IEEEAddress, network zigbee.NetworkAddress, facts ...DeviceFact) *Device {
	_, present := z.devices[ieee]

	if present {
		z.devices[ieee].NetworkAddress = network
	} else {
		z.devices[ieee] = &Device{
			NetworkAddress: network,
			IEEEAddress:    ieee,
			Role:           RoleUnknown,
			Neighbours:     map[zigbee.IEEEAddress]*DeviceNeighbour{},
		}
	}

	for _, f := range facts {
		f(z.devices[ieee])
	}

	return z.devices[ieee]
}

func (z *ZStack) removeDevice(ieee zigbee.IEEEAddress) {
	delete(z.devices, ieee)
}

type DeviceRole uint8

const (
	RoleCoordinator DeviceRole = 0x00
	RoleRouter      DeviceRole = 0x01
	RoleEndDevice   DeviceRole = 0x02
	RoleUnknown     DeviceRole = 0xff
)

type Device struct {
	NetworkAddress zigbee.NetworkAddress
	IEEEAddress    zigbee.IEEEAddress
	Role           DeviceRole
	Neighbours     map[zigbee.IEEEAddress]*DeviceNeighbour
}

type DeviceNeighbour struct {
	IEEEAddress zigbee.IEEEAddress
	LQI         uint8
}

type ZdoMGMTLQIReq struct {
	DestinationAddress zigbee.NetworkAddress
	StartIndex         uint8
}

const ZdoMGMTLQIReqID uint8 = 0x31

type ZdoMGMTLQIReqResp GenericZStackStatus

const ZdoMGMTLQIReqRespID uint8 = 0x31

type ZdoMGMTLQINeighbour struct {
	ExtendedPANID  zigbee.ExtendedPANID
	IEEEAddress    zigbee.IEEEAddress
	NetworkAddress zigbee.NetworkAddress
	Status         uint8
	PermitJoining  uint8
	Depth          uint8
	LQI            uint8
}

type ZdoMGMTLQIResp struct {
	SourceAddress         zigbee.NetworkAddress
	Status                ZStackStatus
	NeighbourTableEntries uint8
	StartIndex            uint8
	Neighbors             []ZdoMGMTLQINeighbour `bclength:"8"`
}

const ZdoMGMTLQIRespID uint8 = 0xb1
