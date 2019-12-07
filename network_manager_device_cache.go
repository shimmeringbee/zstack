package zstack

import "github.com/shimmeringbee/zigbee"

func (z *ZStack) addOrUpdateDevice(ieee zigbee.IEEEAddress, facts ...DeviceFact) (*Device, bool) {
	_, present := z.devices[ieee]
	newDevice := false

	if !present {
		z.devices[ieee] = &Device{
			IEEEAddress: ieee,
			Role:        RoleUnknown,
			Neighbours:  map[zigbee.IEEEAddress]*DeviceNeighbour{},
		}

		newDevice = true
	}

	for _, f := range facts {
		f(z.devices[ieee])
	}

	return z.devices[ieee], newDevice
}

func (z *ZStack) getDevice(netaddr zigbee.NetworkAddress) (*Device, bool) {
	ieee, found := z.devicesByNetAddr[netaddr]

	if found {
		device, found := z.devices[ieee]
		return device, found
	}

	return nil, found
}

type DeviceFact func(*Device)

func (z *ZStack) Role(role DeviceRole) DeviceFact {
	return func(device *Device) {
		if device.Role != role && role == RoleRouter {
			go z.pollDeviceForNetworkStatus(*device)
		}

		device.Role = role
	}
}

func (z *ZStack) NetAddr(networkAddress zigbee.NetworkAddress) DeviceFact {
	return func(device *Device) {
		device.NetworkAddress = networkAddress
		z.devicesByNetAddr[networkAddress] = device.IEEEAddress
	}
}

func (z *ZStack) removeDevice(ieee zigbee.IEEEAddress) {
	device, found := z.devices[ieee]

	if found {
		delete(z.devices, ieee)
		delete(z.devicesByNetAddr, device.NetworkAddress)
	}

	for _, device := range z.devices {
		delete(device.Neighbours, ieee)
	}
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

type DeviceRelationship uint8

const (
	RelationshipParent  DeviceRelationship = 0x00
	RelationshipChild   DeviceRelationship = 0x01
	RelationshipSibling DeviceRelationship = 0x02
	RelationshipUnknown DeviceRelationship = 0x03
)

type DeviceNeighbour struct {
	Relationship DeviceRelationship
	LQI          uint8
}
