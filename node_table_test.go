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
		nt := newNodeTable()

		nt.addOrUpdate(ieee, network)

		node, found := nt.getByIEEE(ieee)

		assert.True(t, found)
		assert.Equal(t, ieee, node.IEEEAddress)
		assert.Equal(t, network, node.NetworkAddress)
		assert.Equal(t, zigbee.Unknown, node.LogicalType)
	})

	t.Run("an added node with updates can be retrieved and has updated information", func(t *testing.T) {
		nt := newNodeTable()

		nt.addOrUpdate(ieee, network, logicalType(zigbee.EndDevice))

		node, found := nt.getByIEEE(ieee)

		assert.True(t, found)
		assert.Equal(t, zigbee.EndDevice, node.LogicalType)
	})

	t.Run("an added node can be retrieved by its network address", func(t *testing.T) {
		nt := newNodeTable()

		nt.addOrUpdate(ieee, network)

		_, found := nt.getByNetwork(network)
		assert.True(t, found)
	})

	t.Run("a missing node queried by its ieee address returns not found", func(t *testing.T) {
		nt := newNodeTable()

		_, found := nt.getByIEEE(ieee)
		assert.False(t, found)
	})

	t.Run("a missing node queried by its network address returns not found", func(t *testing.T) {
		nt := newNodeTable()

		_, found := nt.getByNetwork(network)
		assert.False(t, found)
	})

	t.Run("removing a node results in it not being found by ieee address", func(t *testing.T) {
		nt := newNodeTable()

		nt.addOrUpdate(ieee, network)
		nt.remove(ieee)

		_, found := nt.getByIEEE(ieee)
		assert.False(t, found)
	})

	t.Run("removing a node results in it not being found by network address", func(t *testing.T) {
		nt := newNodeTable()

		nt.addOrUpdate(ieee, network)
		nt.remove(ieee)

		_, found := nt.getByNetwork(network)
		assert.False(t, found)
	})

	t.Run("an update using add makes the node available under the new network only, and updates the network address", func(t *testing.T) {
		nt := newNodeTable()

		newNetwork := zigbee.NetworkAddress(0x1234)

		nt.addOrUpdate(ieee, network)
		nt.addOrUpdate(ieee, newNetwork)

		_, found := nt.getByNetwork(network)
		assert.False(t, found)

		node, found := nt.getByNetwork(newNetwork)
		assert.True(t, found)

		assert.Equal(t, newNetwork, node.NetworkAddress)
	})

	t.Run("an update makes all changes as requested by node updates", func(t *testing.T) {
		nt := newNodeTable()

		nt.addOrUpdate(ieee, network)

		nt.update(ieee, logicalType(zigbee.EndDevice))

		d, _ := nt.getByIEEE(ieee)

		assert.Equal(t, zigbee.EndDevice, d.LogicalType)
	})

	t.Run("returns all nodes when queried", func(t *testing.T) {
		nt := newNodeTable()

		nt.addOrUpdate(ieee, network)

		nodes := nt.nodes()
		assert.Equal(t, 1, len(nodes))
	})

	t.Run("callbacks are called for additions", func(t *testing.T) {
		callbackCalled := false

		nt := newNodeTable()
		nt.registerCallback(func(node zigbee.Node) {
			callbackCalled = true
		})

		nt.addOrUpdate(zigbee.IEEEAddress(0x00), zigbee.NetworkAddress(0x00))

		assert.True(t, callbackCalled)
	})

	t.Run("callbacks are called for additions", func(t *testing.T) {
		callbackCalled := false

		nt := newNodeTable()

		nt.addOrUpdate(zigbee.IEEEAddress(0x00), zigbee.NetworkAddress(0x00))

		nt.registerCallback(func(node zigbee.Node) {
			callbackCalled = true
		})

		nt.update(zigbee.IEEEAddress(0x00), updateReceived)

		assert.True(t, callbackCalled)
	})

	t.Run("dumping and loading result in the same nodes being present in the table", func(t *testing.T) {
		ntOne := newNodeTable()
		ntOne.addOrUpdate(zigbee.IEEEAddress(0x01), zigbee.NetworkAddress(0x01))
		ntOneDump := ntOne.nodes()

		ntTwo := newNodeTable()
		ntTwo.Load(ntOneDump)
		ntTwoDump := ntTwo.nodes()

		assert.Equal(t, ntOneDump, ntTwoDump)
	})
}

func TestNodeUpdate(t *testing.T) {
	t.Run("logicalType updates the logical type of node", func(t *testing.T) {
		node := &zigbee.Node{}

		logicalType(zigbee.EndDevice)(node)

		assert.Equal(t, zigbee.EndDevice, node.LogicalType)
	})

	t.Run("lqi updates the lqi of node", func(t *testing.T) {
		node := &zigbee.Node{}

		lqi(48)(node)

		assert.Equal(t, uint8(48), node.LQI)
	})

	t.Run("depth updates the depth of node", func(t *testing.T) {
		node := &zigbee.Node{}

		depth(3)(node)

		assert.Equal(t, uint8(3), node.Depth)
	})

	t.Run("updateReceived updates the last received time of node", func(t *testing.T) {
		node := &zigbee.Node{}

		updateReceived(node)

		assert.NotEqual(t, time.Time{}, node.LastReceived)
	})

	t.Run("updateDiscovered updates the last received time of node", func(t *testing.T) {
		node := &zigbee.Node{}

		updateDiscovered(node)

		assert.NotEqual(t, time.Time{}, node.LastDiscovered)
	})
}
