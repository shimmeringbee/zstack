package zstack

import (
	"github.com/shimmeringbee/zigbee"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDeviceTable(t *testing.T) {
	ieee := zigbee.IEEEAddress(0x001122334455)
	network := zigbee.NetworkAddress(0xaabb)

	t.Run("an added device can be retrieved by its IEEE address, and has minimum information", func(t *testing.T) {
		dt := NewDeviceTable()

		dt.AddOrUpdate(ieee, network)

		device, found := dt.GetByIEEE(ieee)

		assert.True(t, found)
		assert.Equal(t, ieee, device.IEEEAddress)
		assert.Equal(t, network, device.NetworkAddress)
		assert.Equal(t, zigbee.Unknown, device.LogicalType)
	})

	t.Run("an added device with updates can be retrieved and has updated information", func(t *testing.T) {
		dt := NewDeviceTable()

		dt.AddOrUpdate(ieee, network, LogicalType(zigbee.EndDevice))

		device, found := dt.GetByIEEE(ieee)

		assert.True(t, found)
		assert.Equal(t, zigbee.EndDevice, device.LogicalType)
	})

	t.Run("an added device can be retrieved by its network address", func(t *testing.T) {
		dt := NewDeviceTable()

		dt.AddOrUpdate(ieee, network)

		_, found := dt.GetByNetwork(network)
		assert.True(t, found)
	})

	t.Run("a missing device queried by its ieee address returns not found", func(t *testing.T) {
		dt := NewDeviceTable()

		_, found := dt.GetByIEEE(ieee)
		assert.False(t, found)
	})

	t.Run("a missing device queried by its network address returns not found", func(t *testing.T) {
		dt := NewDeviceTable()

		_, found := dt.GetByNetwork(network)
		assert.False(t, found)
	})

	t.Run("removing a device results in it not being found by ieee address", func(t *testing.T) {
		dt := NewDeviceTable()

		dt.AddOrUpdate(ieee, network)
		dt.Remove(ieee)

		_, found := dt.GetByIEEE(ieee)
		assert.False(t, found)
	})

	t.Run("removing a device results in it not being found by network address", func(t *testing.T) {
		dt := NewDeviceTable()

		dt.AddOrUpdate(ieee, network)
		dt.Remove(ieee)

		_, found := dt.GetByNetwork(network)
		assert.False(t, found)
	})

	t.Run("an update using add makes the device available under the new network only, and updates the network address", func(t *testing.T) {
		dt := NewDeviceTable()

		newNetwork := zigbee.NetworkAddress(0x1234)

		dt.AddOrUpdate(ieee, network)
		dt.AddOrUpdate(ieee, newNetwork)

		_, found := dt.GetByNetwork(network)
		assert.False(t, found)

		device, found := dt.GetByNetwork(newNetwork)
		assert.True(t, found)

		assert.Equal(t, newNetwork, device.NetworkAddress)
	})

	t.Run("an update makes all changes as requested by device updates", func(t *testing.T) {
		dt := NewDeviceTable()

		dt.AddOrUpdate(ieee, network)

		dt.Update(ieee, LogicalType(zigbee.EndDevice))

		d, _ := dt.GetByIEEE(ieee)

		assert.Equal(t, zigbee.EndDevice, d.LogicalType)
	})

	t.Run("returns all devices when queried", func(t *testing.T) {
		dt := NewDeviceTable()

		dt.AddOrUpdate(ieee, network)

		devices := dt.GetAllDevices()
		assert.Equal(t, 1, len(devices))
	})

	t.Run("callbacks are called for additions", func(t *testing.T) {
		callbackCalled := false

		dt := NewDeviceTable()
		dt.RegisterCallback(func(device zigbee.Device) {
			callbackCalled = true
		})

		dt.AddOrUpdate(zigbee.IEEEAddress(0x00), zigbee.NetworkAddress(0x00))

		assert.True(t, callbackCalled)
	})

	t.Run("callbacks are called for additions", func(t *testing.T) {
		callbackCalled := false

		dt := NewDeviceTable()

		dt.AddOrUpdate(zigbee.IEEEAddress(0x00), zigbee.NetworkAddress(0x00))

		dt.RegisterCallback(func(device zigbee.Device) {
			callbackCalled = true
		})

		dt.Update(zigbee.IEEEAddress(0x00), UpdateReceived)

		assert.True(t, callbackCalled)
	})
}

func TestDeviceUpdate(t *testing.T) {
	t.Run("LogicalType updates the logical type of device", func(t *testing.T) {
		device := &zigbee.Device{}

		LogicalType(zigbee.EndDevice)(device)

		assert.Equal(t, zigbee.EndDevice, device.LogicalType)
	})

	t.Run("LQI updates the lqi of device", func(t *testing.T) {
		device := &zigbee.Device{}

		LQI(48)(device)

		assert.Equal(t, uint8(48), device.LQI)
	})

	t.Run("Depth updates the depth of device", func(t *testing.T) {
		device := &zigbee.Device{}

		Depth(3)(device)

		assert.Equal(t, uint8(3), device.Depth)
	})

	t.Run("UpdateReceived updates the last received time of device", func(t *testing.T) {
		device := &zigbee.Device{}

		UpdateReceived(device)

		assert.NotEqual(t, time.Time{}, device.LastReceived)
	})

	t.Run("UpdateDiscovered updates the last received time of device", func(t *testing.T) {
		device := &zigbee.Device{}

		UpdateDiscovered(device)

		assert.NotEqual(t, time.Time{}, device.LastDiscovered)
	})
}
