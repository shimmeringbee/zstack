package zstack

import (
	"github.com/shimmeringbee/zigbee"
	"time"
)

type DeviceTable struct {
	ieeeToDevice  map[zigbee.IEEEAddress]*Device
	networkToIEEE map[zigbee.NetworkAddress]zigbee.IEEEAddress
}

func NewDeviceTable() *DeviceTable {
	return &DeviceTable{
		ieeeToDevice:  make(map[zigbee.IEEEAddress]*Device),
		networkToIEEE: make(map[zigbee.NetworkAddress]zigbee.IEEEAddress),
	}
}

func (t *DeviceTable) GetAllDevices() []Device {
	var devices []Device

	for _, device := range t.ieeeToDevice {
		devices = append(devices, *device)
	}

	return devices
}

func (t *DeviceTable) GetByIEEE(ieeeAddress zigbee.IEEEAddress) (Device, bool) {
	device, found := t.ieeeToDevice[ieeeAddress]

	if found {
		return *device, found
	} else {
		return Device{}, false
	}
}

func (t *DeviceTable) GetByNetwork(networkAddress zigbee.NetworkAddress) (Device, bool) {
	ieee, found := t.networkToIEEE[networkAddress]

	if !found {
		return Device{}, false
	} else {
		return t.GetByIEEE(ieee)
	}
}

func (t *DeviceTable) AddOrUpdate(ieeeAddress zigbee.IEEEAddress, networkAddress zigbee.NetworkAddress, updates ...DeviceUpdate) {
	device, found := t.ieeeToDevice[ieeeAddress]

	if found {
		delete(t.networkToIEEE, device.NetworkAddress)
		device.NetworkAddress = networkAddress
	} else {
		t.ieeeToDevice[ieeeAddress] = &Device{
			IEEEAddress:    ieeeAddress,
			NetworkAddress: networkAddress,
		}
	}

	t.networkToIEEE[networkAddress] = ieeeAddress
	t.Update(ieeeAddress, updates...)
}

func (t *DeviceTable) Update(ieeeAddress zigbee.IEEEAddress, updates ...DeviceUpdate) {
	device, found := t.ieeeToDevice[ieeeAddress]

	if found {
		for _, du := range updates {
			du(device)
		}
	}
}

func (t *DeviceTable) Remove(ieeeAddress zigbee.IEEEAddress) {
	device, found := t.GetByIEEE(ieeeAddress)

	if found {
		delete(t.networkToIEEE, device.NetworkAddress)
		delete(t.ieeeToDevice, device.IEEEAddress)
	}
}

type DeviceUpdate func(device *Device)

func LogicalType(logicalType zigbee.LogicalType) DeviceUpdate {
	return func(device *Device) {
		device.LogicalType = logicalType
	}
}

func LQI(lqi uint8) DeviceUpdate {
	return func(device *Device) {
		device.LQI = lqi
	}
}

func Depth(depth uint8) DeviceUpdate {
	return func(device *Device) {
		device.Depth = depth
	}
}

func UpdateReceived(device *Device) {
	device.LastReceived = time.Now()
}

func UpdateDiscovered(device *Device) {
	device.LastDiscovered = time.Now()
}

type Device struct {
	IEEEAddress    zigbee.IEEEAddress
	NetworkAddress zigbee.NetworkAddress
	LogicalType    zigbee.LogicalType
	LQI            uint8
	Depth          uint8
	LastDiscovered time.Time
	LastReceived   time.Time
}
