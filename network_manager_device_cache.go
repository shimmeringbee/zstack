package zstack

import "github.com/shimmeringbee/zigbee"

func (z *ZStack) addOrUpdateDevice(ieee zigbee.IEEEAddress, facts ...DeviceFact) (*Device, bool) {
	_, present := z.devices[ieee]
	newDevice := false

	if !present {
		z.devices[ieee] = &Device{
			IEEEAddress:    ieee,
			Role:           RoleUnknown,
			Neighbours:     map[zigbee.IEEEAddress]*DeviceNeighbour{},
		}

		newDevice = true
	}

	for _, f := range facts {
		f(z.devices[ieee])
	}

	return z.devices[ieee], newDevice
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
	LQI         uint8
}
