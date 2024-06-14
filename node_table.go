package zstack

import (
	"github.com/shimmeringbee/zigbee"
	"sync"
	"time"
)

type nodeTable struct {
	callbacks     []func(zigbee.Node)
	ieeeToNode    map[zigbee.IEEEAddress]*zigbee.Node
	networkToIEEE map[zigbee.NetworkAddress]zigbee.IEEEAddress
	lock          *sync.RWMutex
}

func newNodeTable() *nodeTable {
	return &nodeTable{
		callbacks:     []func(zigbee.Node){},
		ieeeToNode:    make(map[zigbee.IEEEAddress]*zigbee.Node),
		networkToIEEE: make(map[zigbee.NetworkAddress]zigbee.IEEEAddress),
		lock:          &sync.RWMutex{},
	}
}

func (t *nodeTable) registerCallback(cb func(zigbee.Node)) {
	t.callbacks = append(t.callbacks, cb)
}

func (t *nodeTable) Load(nodes []zigbee.Node) {
	for _, node := range nodes {
		t.addOrUpdate(node.IEEEAddress, node.NetworkAddress, logicalType(node.LogicalType), lqi(node.LQI), depth(node.Depth), setReceived(node.LastReceived), setDiscovered(node.LastDiscovered))
	}
}

func (t *nodeTable) nodes() []zigbee.Node {
	t.lock.RLock()
	defer t.lock.RUnlock()

	var nodes []zigbee.Node

	for _, node := range t.ieeeToNode {
		nodes = append(nodes, *node)
	}

	return nodes
}

func (t *nodeTable) getByIEEE(ieeeAddress zigbee.IEEEAddress) (zigbee.Node, bool) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	node, found := t.ieeeToNode[ieeeAddress]

	if found {
		return *node, found
	} else {
		return zigbee.Node{}, false
	}
}

func (t *nodeTable) getByNetwork(networkAddress zigbee.NetworkAddress) (zigbee.Node, bool) {
	t.lock.RLock()
	ieee, found := t.networkToIEEE[networkAddress]
	t.lock.RUnlock()

	if !found {
		return zigbee.Node{}, false
	} else {
		return t.getByIEEE(ieee)
	}
}

func (t *nodeTable) addOrUpdate(ieeeAddress zigbee.IEEEAddress, networkAddress zigbee.NetworkAddress, updates ...nodeUpdate) {
	t.lock.Lock()
	node, found := t.ieeeToNode[ieeeAddress]

	if found {
		if node.NetworkAddress != networkAddress {
			delete(t.networkToIEEE, node.NetworkAddress)
			node.NetworkAddress = networkAddress
		}
	} else {
		t.ieeeToNode[ieeeAddress] = &zigbee.Node{
			IEEEAddress:    ieeeAddress,
			NetworkAddress: networkAddress,
			LogicalType:    zigbee.Unknown,
		}
	}

	t.networkToIEEE[networkAddress] = ieeeAddress
	t.lock.Unlock()

	t.update(ieeeAddress, updates...)
}

func (t *nodeTable) update(ieeeAddress zigbee.IEEEAddress, updates ...nodeUpdate) {
	t.lock.Lock()
	defer t.lock.Unlock()

	node, found := t.ieeeToNode[ieeeAddress]

	if found {
		for _, du := range updates {
			du(node)
		}

		for _, cb := range t.callbacks {
			cb(*node)
		}
	}
}

func (t *nodeTable) remove(ieeeAddress zigbee.IEEEAddress) {
	node, found := t.getByIEEE(ieeeAddress)

	t.lock.Lock()
	defer t.lock.Unlock()

	if found {
		delete(t.networkToIEEE, node.NetworkAddress)
		delete(t.ieeeToNode, node.IEEEAddress)
	}
}

type nodeUpdate func(device *zigbee.Node)

func logicalType(logicalType zigbee.LogicalType) nodeUpdate {
	return func(node *zigbee.Node) {
		node.LogicalType = logicalType
	}
}

func lqi(lqi uint8) nodeUpdate {
	return func(node *zigbee.Node) {
		node.LQI = lqi
	}
}

func depth(depth uint8) nodeUpdate {
	return func(node *zigbee.Node) {
		node.Depth = depth
	}
}

func updateReceived(node *zigbee.Node) {
	node.LastReceived = time.Now()
}

func updateDiscovered(node *zigbee.Node) {
	node.LastDiscovered = time.Now()
}

func setReceived(t time.Time) nodeUpdate {
	return func(node *zigbee.Node) {
		node.LastReceived = t
	}
}

func setDiscovered(t time.Time) nodeUpdate {
	return func(node *zigbee.Node) {
		node.LastDiscovered = t
	}
}
