package zstack

import (
	"github.com/shimmeringbee/zigbee"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNodeTable(t *testing.T) {
	ieee := zigbee.IEEEAddress(0x001122334455)
	network := zigbee.NetworkAddress(0xaabb)

	t.Run("an added node can be retrieved by its IEEE address, and has minimum information", func(t *testing.T) {
		nt := NewNodeTable()

		nt.AddOrUpdate(ieee, network)

		node, found := nt.GetByIEEE(ieee)

		assert.True(t, found)
		assert.Equal(t, ieee, node.IEEEAddress)
		assert.Equal(t, network, node.NetworkAddress)
		assert.Equal(t, zigbee.Unknown, node.LogicalType)
	})

	t.Run("an added node with updates can be retrieved and has updated information", func(t *testing.T) {
		nt := NewNodeTable()

		nt.AddOrUpdate(ieee, network, LogicalType(zigbee.EndDevice))

		node, found := nt.GetByIEEE(ieee)

		assert.True(t, found)
		assert.Equal(t, zigbee.EndDevice, node.LogicalType)
	})

	t.Run("an added node can be retrieved by its network address", func(t *testing.T) {
		nt := NewNodeTable()

		nt.AddOrUpdate(ieee, network)

		_, found := nt.GetByNetwork(network)
		assert.True(t, found)
	})

	t.Run("a missing node queried by its ieee address returns not found", func(t *testing.T) {
		nt := NewNodeTable()

		_, found := nt.GetByIEEE(ieee)
		assert.False(t, found)
	})

	t.Run("a missing node queried by its network address returns not found", func(t *testing.T) {
		nt := NewNodeTable()

		_, found := nt.GetByNetwork(network)
		assert.False(t, found)
	})

	t.Run("removing a node results in it not being found by ieee address", func(t *testing.T) {
		nt := NewNodeTable()

		nt.AddOrUpdate(ieee, network)
		nt.Remove(ieee)

		_, found := nt.GetByIEEE(ieee)
		assert.False(t, found)
	})

	t.Run("removing a node results in it not being found by network address", func(t *testing.T) {
		nt := NewNodeTable()

		nt.AddOrUpdate(ieee, network)
		nt.Remove(ieee)

		_, found := nt.GetByNetwork(network)
		assert.False(t, found)
	})

	t.Run("an update using add makes the node available under the new network only, and updates the network address", func(t *testing.T) {
		nt := NewNodeTable()

		newNetwork := zigbee.NetworkAddress(0x1234)

		nt.AddOrUpdate(ieee, network)
		nt.AddOrUpdate(ieee, newNetwork)

		_, found := nt.GetByNetwork(network)
		assert.False(t, found)

		node, found := nt.GetByNetwork(newNetwork)
		assert.True(t, found)

		assert.Equal(t, newNetwork, node.NetworkAddress)
	})

	t.Run("an update makes all changes as requested by node updates", func(t *testing.T) {
		nt := NewNodeTable()

		nt.AddOrUpdate(ieee, network)

		nt.Update(ieee, LogicalType(zigbee.EndDevice))

		d, _ := nt.GetByIEEE(ieee)

		assert.Equal(t, zigbee.EndDevice, d.LogicalType)
	})

	t.Run("returns all nodes when queried", func(t *testing.T) {
		nt := NewNodeTable()

		nt.AddOrUpdate(ieee, network)

		nodes := nt.GetAllNodes()
		assert.Equal(t, 1, len(nodes))
	})

	t.Run("callbacks are called for additions", func(t *testing.T) {
		callbackCalled := false

		nt := NewNodeTable()
		nt.RegisterCallback(func(node zigbee.Node) {
			callbackCalled = true
		})

		nt.AddOrUpdate(zigbee.IEEEAddress(0x00), zigbee.NetworkAddress(0x00))

		assert.True(t, callbackCalled)
	})

	t.Run("callbacks are called for additions", func(t *testing.T) {
		callbackCalled := false

		nt := NewNodeTable()

		nt.AddOrUpdate(zigbee.IEEEAddress(0x00), zigbee.NetworkAddress(0x00))

		nt.RegisterCallback(func(node zigbee.Node) {
			callbackCalled = true
		})

		nt.Update(zigbee.IEEEAddress(0x00), UpdateReceived)

		assert.True(t, callbackCalled)
	})
}

func TestNodeUpdate(t *testing.T) {
	t.Run("LogicalType updates the logical type of node", func(t *testing.T) {
		node := &zigbee.Node{}

		LogicalType(zigbee.EndDevice)(node)

		assert.Equal(t, zigbee.EndDevice, node.LogicalType)
	})

	t.Run("LQI updates the lqi of node", func(t *testing.T) {
		node := &zigbee.Node{}

		LQI(48)(node)

		assert.Equal(t, uint8(48), node.LQI)
	})

	t.Run("Depth updates the depth of node", func(t *testing.T) {
		node := &zigbee.Node{}

		Depth(3)(node)

		assert.Equal(t, uint8(3), node.Depth)
	})

	t.Run("UpdateReceived updates the last received time of node", func(t *testing.T) {
		node := &zigbee.Node{}

		UpdateReceived(node)

		assert.NotEqual(t, time.Time{}, node.LastReceived)
	})

	t.Run("UpdateDiscovered updates the last received time of node", func(t *testing.T) {
		node := &zigbee.Node{}

		UpdateDiscovered(node)

		assert.NotEqual(t, time.Time{}, node.LastDiscovered)
	})
}
