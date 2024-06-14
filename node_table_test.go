package zstack

import (
	"github.com/shimmeringbee/persistence/converter"
	"github.com/shimmeringbee/persistence/impl/memory"
	"github.com/shimmeringbee/zigbee"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNodeTable(t *testing.T) {
	ieee := zigbee.GenerateLocalAdministeredIEEEAddress()
	network := zigbee.NetworkAddress(0xaabb)

	t.Run("an added node can be retrieved by its IEEE address, and has minimum information", func(t *testing.T) {
		s := memory.New()
		nt := newNodeTable(s)

		nt.addOrUpdate(ieee, network)

		node, found := nt.getByIEEE(ieee)

		assert.True(t, found)
		assert.Equal(t, ieee, node.IEEEAddress)
		assert.Equal(t, network, node.NetworkAddress)
		assert.Equal(t, zigbee.Unknown, node.LogicalType)

		assert.Contains(t, s.Keys(), ieee.String())

		ns := s.Section(ieee.String())
		na, ok := converter.Retrieve(ns, "NetworkAddress", converter.NetworkAddressDecoder)

		assert.True(t, ok)
		assert.Equal(t, network, na)
	})

	t.Run("an added node with updates can be retrieved and has updated information", func(t *testing.T) {
		s := memory.New()
		nt := newNodeTable(s)

		nt.addOrUpdate(ieee, network, logicalType(zigbee.EndDevice))

		node, found := nt.getByIEEE(ieee)

		assert.True(t, found)
		assert.Equal(t, zigbee.EndDevice, node.LogicalType)

		ns := s.Section(ieee.String())
		lt, ok := converter.Retrieve(ns, "LogicalType", converter.LogicalTypeDecoder)

		assert.True(t, ok)
		assert.Equal(t, zigbee.EndDevice, lt)
	})

	t.Run("an added node can be retrieved by its network address", func(t *testing.T) {
		nt := newNodeTable(memory.New())

		nt.addOrUpdate(ieee, network)

		_, found := nt.getByNetwork(network)
		assert.True(t, found)
	})

	t.Run("a missing node queried by its ieee address returns not found", func(t *testing.T) {
		nt := newNodeTable(memory.New())

		_, found := nt.getByIEEE(ieee)
		assert.False(t, found)
	})

	t.Run("a missing node queried by its network address returns not found", func(t *testing.T) {
		nt := newNodeTable(memory.New())

		_, found := nt.getByNetwork(network)
		assert.False(t, found)
	})

	t.Run("removing a node results in it not being found by ieee address", func(t *testing.T) {
		nt := newNodeTable(memory.New())

		nt.addOrUpdate(ieee, network)
		nt.remove(ieee)

		_, found := nt.getByIEEE(ieee)
		assert.False(t, found)
	})

	t.Run("removing a node results in it not being found by network address", func(t *testing.T) {
		s := memory.New()
		nt := newNodeTable(s)

		nt.addOrUpdate(ieee, network)
		assert.Contains(t, s.Keys(), ieee.String())

		nt.remove(ieee)
		assert.NotContains(t, s.Keys(), ieee.String())

		_, found := nt.getByNetwork(network)
		assert.False(t, found)
	})

	t.Run("an update using add makes the node available under the new network only, and updates the network address", func(t *testing.T) {
		s := memory.New()
		nt := newNodeTable(s)

		newNetwork := zigbee.NetworkAddress(0x1234)

		nt.addOrUpdate(ieee, network)
		nt.addOrUpdate(ieee, newNetwork)

		_, found := nt.getByNetwork(network)
		assert.False(t, found)

		node, found := nt.getByNetwork(newNetwork)
		assert.True(t, found)

		assert.Equal(t, newNetwork, node.NetworkAddress)

		ns := s.Section(ieee.String())
		na, ok := converter.Retrieve(ns, "NetworkAddress", converter.NetworkAddressDecoder)

		assert.True(t, ok)
		assert.Equal(t, newNetwork, na)
	})

	t.Run("an update makes all changes as requested by node updates", func(t *testing.T) {
		nt := newNodeTable(memory.New())

		nt.addOrUpdate(ieee, network)

		nt.update(ieee, logicalType(zigbee.EndDevice))

		d, _ := nt.getByIEEE(ieee)

		assert.Equal(t, zigbee.EndDevice, d.LogicalType)
	})

	t.Run("returns all nodes when queried", func(t *testing.T) {
		nt := newNodeTable(memory.New())

		nt.addOrUpdate(ieee, network)

		nodes := nt.nodes()
		assert.Equal(t, 1, len(nodes))
	})

	t.Run("callbacks are called for additions", func(t *testing.T) {
		callbackCalled := false

		nt := newNodeTable(memory.New())
		nt.registerCallback(func(node zigbee.Node) {
			callbackCalled = true
		})

		nt.addOrUpdate(zigbee.IEEEAddress(0x00), zigbee.NetworkAddress(0x00))

		assert.True(t, callbackCalled)
	})

	t.Run("callbacks are called for additions", func(t *testing.T) {
		callbackCalled := false

		nt := newNodeTable(memory.New())

		nt.addOrUpdate(zigbee.IEEEAddress(0x00), zigbee.NetworkAddress(0x00))

		nt.registerCallback(func(node zigbee.Node) {
			callbackCalled = true
		})

		nt.update(zigbee.IEEEAddress(0x00), updateReceived())

		assert.True(t, callbackCalled)
	})
}

func TestNodeUpdate(t *testing.T) {
	t.Run("logicalType updates the logical type of node", func(t *testing.T) {
		node := &zigbee.Node{}

		s := memory.New()
		logicalType(zigbee.EndDevice)(node, s)

		assert.Equal(t, zigbee.EndDevice, node.LogicalType)

		lt, ok := converter.Retrieve(s, "LogicalType", converter.LogicalTypeDecoder)

		assert.True(t, ok)
		assert.Equal(t, zigbee.EndDevice, lt)
	})

	t.Run("lqi updates the lqi of node", func(t *testing.T) {
		node := &zigbee.Node{}

		s := memory.New()
		lqi(48)(node, s)

		assert.Equal(t, uint8(48), node.LQI)

		l, ok := s.UInt("LQI")
		assert.True(t, ok)
		assert.Equal(t, uint64(48), l)
	})

	t.Run("depth updates the depth of node", func(t *testing.T) {
		node := &zigbee.Node{}

		s := memory.New()
		depth(3)(node, s)

		assert.Equal(t, uint8(3), node.Depth)

		d, ok := s.UInt("Depth")
		assert.True(t, ok)
		assert.Equal(t, uint64(3), d)
	})

	t.Run("updateReceived updates the last received time of node", func(t *testing.T) {
		node := &zigbee.Node{}

		s := memory.New()
		updateReceived()(node, s)

		assert.NotEqual(t, time.Time{}, node.LastReceived)

		date, ok := converter.Retrieve(s, "LastReceived", converter.TimeDecoder)
		assert.True(t, ok)
		assert.True(t, time.Now().After(date))
	})

	t.Run("updateDiscovered updates the last received time of node", func(t *testing.T) {
		node := &zigbee.Node{}

		s := memory.New()
		updateDiscovered()(node, s)

		assert.NotEqual(t, time.Time{}, node.LastDiscovered)

		date, ok := converter.Retrieve(s, "LastDiscovered", converter.TimeDecoder)
		assert.True(t, ok)
		assert.True(t, time.Now().After(date))
	})
}

func TestNodeTable_Load(t *testing.T) {
	t.Run("loading a table from persistence contains expected nodes", func(t *testing.T) {
		s := memory.New()

		ieee := zigbee.GenerateLocalAdministeredIEEEAddress()

		time := time.UnixMilli(time.Now().UnixMilli())

		nS := s.Section(ieee.String())
		converter.Store(nS, "NetworkAddress", zigbee.NetworkAddress(0x1122), converter.NetworkAddressEncoder)
		converter.Store(nS, "LastReceived", time, converter.TimeEncoder)
		converter.Store(nS, "LastDiscovered", time, converter.TimeEncoder)
		converter.Store(nS, "LogicalType", zigbee.Router, converter.LogicalTypeEncoder)
		nS.Set("LQI", uint64(8))
		nS.Set("Depth", uint64(2))

		nt := newNodeTable(s)

		node, ok := nt.getByIEEE(ieee)
		assert.True(t, ok)
		assert.Equal(t, zigbee.NetworkAddress(0x1122), node.NetworkAddress)
		assert.Equal(t, time, node.LastDiscovered)
		assert.Equal(t, time, node.LastReceived)
		assert.Equal(t, uint8(8), node.LQI)
		assert.Equal(t, uint8(2), node.Depth)
	})
}
