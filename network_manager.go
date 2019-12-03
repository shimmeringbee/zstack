package zstack

import (
	"context"
	"encoding/json"
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
	z.addOrUpdateDevice(z.NetworkProperties.IEEEAddress, z.NetworkProperties.NetworkAddress).Role = RoleCoordinator

	immediateStart := make(chan bool, 1)
	defer close(immediateStart)
	immediateStart <- true

	_, cancel := z.subscriber.Subscribe(&ZdoMGMTLQIResp{}, z.receiveLQIUpdate)
	defer cancel()

	_, cancel = z.subscriber.Subscribe(&ZdoEndDeviceAnnceInd{}, z.handleEndDeviceAnnouncement)
	defer cancel()

	_, cancel = z.subscriber.Subscribe(&ZdoLeaveInd{}, z.handleLeaveAnnouncement)
	defer cancel()

	for {
		select {
		case <-immediateStart:
			z.pollForNetworkStatus()
		case <-time.After(defaultPollingInterval * time.Second):
			z.pollForNetworkStatus()
		case <-z.networkManagerStop:
			return
		case ue := <-z.networkManagerIncoming:
			switch e := ue.(type) {
			case ZdoMGMTLQIResp:
				d, _ := json.MarshalIndent(e, "", "\t")
				fmt.Println(string(d))
			case ZdoEndDeviceAnnceInd:
				z.addOrUpdateDevice(e.IEEEAddress, e.NetworkAddress)
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

func (z *ZStack) pollForNetworkStatus() {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultZStackTimeout)
	defer cancel()

	resp := ZdoMGMTLQIReqResp{}

	if err := z.requestResponder.RequestResponse(ctx, ZdoMGMTLQIReq{DestinationAddress: 0, StartIndex: 0}, &resp); err != nil {
		log.Printf("failed to request lqi tables: %v\n", err)
	}

	if resp.Status != ZSuccess {
		log.Printf("failed to request lqi tables: from adapter\n")
	}
}

func (z *ZStack) receiveLQIUpdate(u func(interface{}) error) {
	msg := ZdoMGMTLQIResp{}
	if err := u(&msg); err == nil {
		z.networkManagerIncoming <- msg
	}
}

func (z *ZStack) handleEndDeviceAnnouncement(u func(interface{}) error) {
	msg := ZdoEndDeviceAnnceInd{}
	var err error
	if err = u(&msg); err == nil {
		z.networkManagerIncoming <- msg
	}
}

func (z *ZStack) handleLeaveAnnouncement(u func(interface{}) error) {
	msg := ZdoLeaveInd{}
	var err error
	if err = u(&msg); err == nil {
		z.networkManagerIncoming <- msg
	}
}

func (z *ZStack) addOrUpdateDevice(ieee zigbee.IEEEAddress, network zigbee.NetworkAddress) *Device {
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
