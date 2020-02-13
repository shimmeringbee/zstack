package zstack

import (
	"github.com/shimmeringbee/zigbee"
	"time"
)

type NodeTable struct {
	callbacks     []func(zigbee.Node)
	ieeeToNode    map[zigbee.IEEEAddress]*zigbee.Node
	networkToIEEE map[zigbee.NetworkAddress]zigbee.IEEEAddress
}

func NewNodeTable() *NodeTable {
	return &NodeTable{
		callbacks:     []func(zigbee.Node){},
		ieeeToNode:    make(map[zigbee.IEEEAddress]*zigbee.Node),
		networkToIEEE: make(map[zigbee.NetworkAddress]zigbee.IEEEAddress),
	}
}

func (t *NodeTable) RegisterCallback(cb func(zigbee.Node)) {
	t.callbacks = append(t.callbacks, cb)
}

func (t *NodeTable) GetAllNodes() []zigbee.Node {
	var nodes []zigbee.Node

	for _, node := range t.ieeeToNode {
		nodes = append(nodes, *node)
	}

	return nodes
}

func (t *NodeTable) GetByIEEE(ieeeAddress zigbee.IEEEAddress) (zigbee.Node, bool) {
	node, found := t.ieeeToNode[ieeeAddress]

	if found {
		return *node, found
	} else {
		return zigbee.Node{}, false
	}
}

func (t *NodeTable) GetByNetwork(networkAddress zigbee.NetworkAddress) (zigbee.Node, bool) {
	ieee, found := t.networkToIEEE[networkAddress]

	if !found {
		return zigbee.Node{}, false
	} else {
		return t.GetByIEEE(ieee)
	}
}

func (t *NodeTable) AddOrUpdate(ieeeAddress zigbee.IEEEAddress, networkAddress zigbee.NetworkAddress, updates ...NodeUpdate) {
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
	t.Update(ieeeAddress, updates...)
}

func (t *NodeTable) Update(ieeeAddress zigbee.IEEEAddress, updates ...NodeUpdate) {
	node, found := t.ieeeToNode[ieeeAddress]

	if found {
		for _, du := range updates {
			du(node)
		}
	}

	for _, cb := range t.callbacks {
		cb(*node)
	}
}

func (t *NodeTable) Remove(ieeeAddress zigbee.IEEEAddress) {
	node, found := t.GetByIEEE(ieeeAddress)

	if found {
		delete(t.networkToIEEE, node.NetworkAddress)
		delete(t.ieeeToNode, node.IEEEAddress)
	}
}

type NodeUpdate func(device *zigbee.Node)

func LogicalType(logicalType zigbee.LogicalType) NodeUpdate {
	return func(node *zigbee.Node) {
		node.LogicalType = logicalType
	}
}

func LQI(lqi uint8) NodeUpdate {
	return func(node *zigbee.Node) {
		node.LQI = lqi
	}
}

func Depth(depth uint8) NodeUpdate {
	return func(node *zigbee.Node) {
		node.Depth = depth
	}
}

func UpdateReceived(node *zigbee.Node) {
	node.LastReceived = time.Now()
}

func UpdateDiscovered(node *zigbee.Node) {
	node.LastDiscovered = time.Now()
}
