package zstack

import (
	"github.com/shimmeringbee/zigbee"
	"time"
)

type DeviceTable struct {
	callbacks     []func(zigbee.Device)
	ieeeToDevice  map[zigbee.IEEEAddress]*zigbee.Device
	networkToIEEE map[zigbee.NetworkAddress]zigbee.IEEEAddress
}

func NewDeviceTable() *DeviceTable {
	return &DeviceTable{
		callbacks:     []func(zigbee.Device){},
		ieeeToDevice:  make(map[zigbee.IEEEAddress]*zigbee.Device),
		networkToIEEE: make(map[zigbee.NetworkAddress]zigbee.IEEEAddress),
	}
}

func (t *DeviceTable) RegisterCallback(cb func(zigbee.Device)) {
	t.callbacks = append(t.callbacks, cb)
}

func (t *DeviceTable) GetAllDevices() []zigbee.Device {
	var devices []zigbee.Device

	for _, device := range t.ieeeToDevice {
		devices = append(devices, *device)
	}

	return devices
}

func (t *DeviceTable) GetByIEEE(ieeeAddress zigbee.IEEEAddress) (zigbee.Device, bool) {
	device, found := t.ieeeToDevice[ieeeAddress]

	if found {
		return *device, found
	} else {
		return zigbee.Device{}, false
	}
}

func (t *DeviceTable) GetByNetwork(networkAddress zigbee.NetworkAddress) (zigbee.Device, bool) {
	ieee, found := t.networkToIEEE[networkAddress]

	if !found {
		return zigbee.Device{}, false
	} else {
		return t.GetByIEEE(ieee)
	}
}

func (t *DeviceTable) AddOrUpdate(ieeeAddress zigbee.IEEEAddress, networkAddress zigbee.NetworkAddress, updates ...DeviceUpdate) {
	device, found := t.ieeeToDevice[ieeeAddress]

	if found {
		if device.NetworkAddress != networkAddress {
			delete(t.networkToIEEE, device.NetworkAddress)
			device.NetworkAddress = networkAddress
		}
	} else {
		t.ieeeToDevice[ieeeAddress] = &zigbee.Device{
			IEEEAddress:    ieeeAddress,
			NetworkAddress: networkAddress,
			LogicalType:    zigbee.Unknown,
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

	for _, cb := range t.callbacks {
		cb(*device)
	}
}

func (t *DeviceTable) Remove(ieeeAddress zigbee.IEEEAddress) {
	device, found := t.GetByIEEE(ieeeAddress)

	if found {
		delete(t.networkToIEEE, device.NetworkAddress)
		delete(t.ieeeToDevice, device.IEEEAddress)
	}
}

type DeviceUpdate func(device *zigbee.Device)

func LogicalType(logicalType zigbee.LogicalType) DeviceUpdate {
	return func(device *zigbee.Device) {
		device.LogicalType = logicalType
	}
}

func LQI(lqi uint8) DeviceUpdate {
	return func(device *zigbee.Device) {
		device.LQI = lqi
	}
}

func Depth(depth uint8) DeviceUpdate {
	return func(device *zigbee.Device) {
		device.Depth = depth
	}
}

func UpdateReceived(device *zigbee.Device) {
	device.LastReceived = time.Now()
}

func UpdateDiscovered(device *zigbee.Device) {
	device.LastDiscovered = time.Now()
}
