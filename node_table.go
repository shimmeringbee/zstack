package zstack

import (
	"github.com/shimmeringbee/persistence"
	"github.com/shimmeringbee/persistence/converter"
	"github.com/shimmeringbee/zigbee"
	"strconv"
	"sync"
	"time"
)

type nodeTable struct {
	callbacks     []func(zigbee.Node)
	ieeeToNode    map[zigbee.IEEEAddress]*zigbee.Node
	networkToIEEE map[zigbee.NetworkAddress]zigbee.IEEEAddress
	lock          *sync.RWMutex

	p       persistence.Section
	loading bool
}

func newNodeTable(p persistence.Section) *nodeTable {
	n := &nodeTable{
		callbacks:     []func(zigbee.Node){},
		ieeeToNode:    make(map[zigbee.IEEEAddress]*zigbee.Node),
		networkToIEEE: make(map[zigbee.NetworkAddress]zigbee.IEEEAddress),
		lock:          &sync.RWMutex{},
		p:             p,
	}

	n.load()

	return n
}

func (t *nodeTable) registerCallback(cb func(zigbee.Node)) {
	t.callbacks = append(t.callbacks, cb)
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

	s := t.p.Section(ieeeAddress.String())

	if found {
		if node.NetworkAddress != networkAddress {
			delete(t.networkToIEEE, node.NetworkAddress)
			node.NetworkAddress = networkAddress

			converter.Store(s, "NetworkAddress", node.NetworkAddress, converter.NetworkAddressEncoder)
		}
	} else {
		node = &zigbee.Node{
			IEEEAddress:    ieeeAddress,
			NetworkAddress: networkAddress,
			LogicalType:    zigbee.Unknown,
		}

		t.ieeeToNode[ieeeAddress] = node

		converter.Store(s, "NetworkAddress", node.NetworkAddress, converter.NetworkAddressEncoder)
		converter.Store(s, "LogicalType", node.LogicalType, converter.LogicalTypeEncoder)
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
		var s persistence.Section

		if !t.loading {
			s = t.p.Section(ieeeAddress.String())
		}

		for _, du := range updates {
			du(node, s)
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

	t.p.SectionDelete(ieeeAddress.String())
}

func (t *nodeTable) load() {
	t.lock.Lock()
	t.loading = true
	t.lock.Unlock()

	defer func() {
		t.lock.Lock()
		t.loading = false
		t.lock.Unlock()
	}()

	for _, key := range t.p.SectionKeys() {
		if value, err := strconv.ParseUint(key, 16, 64); err == nil {
			s := t.p.Section(key)
			ieee := zigbee.IEEEAddress(value)

			na, naFound := converter.Retrieve(s, "NetworkAddress", converter.NetworkAddressDecoder)
			if !naFound {
				continue
			}

			t.addOrUpdate(ieee, na)

			if lt, found := converter.Retrieve(s, "LogicalType", converter.LogicalTypeDecoder); found {
				t.update(ieee, logicalType(lt))
			}

			if l, found := s.UInt("LQI"); found {
				t.update(ieee, lqi(uint8(l)))
			}

			if d, found := s.UInt("Depth"); found {
				t.update(ieee, depth(uint8(d)))
			}

			if received, found := converter.Retrieve(s, "LastReceived", converter.TimeDecoder); found {
				t.update(ieee, setReceived(received))

			}

			if discovered, found := converter.Retrieve(s, "LastDiscovered", converter.TimeDecoder); found {
				t.update(ieee, setDiscovered(discovered))
			}
		}
	}
}

type nodeUpdate func(node *zigbee.Node, p persistence.Section)

func logicalType(logicalType zigbee.LogicalType) nodeUpdate {
	return func(node *zigbee.Node, p persistence.Section) {
		node.LogicalType = logicalType

		if p != nil {
			converter.Store(p, "LogicalType", node.LogicalType, converter.LogicalTypeEncoder)
		}
	}
}

func lqi(lqi uint8) nodeUpdate {
	return func(node *zigbee.Node, p persistence.Section) {
		node.LQI = lqi

		if p != nil {
			p.Set("LQI", uint64(node.LQI))
		}
	}
}

func depth(depth uint8) nodeUpdate {
	return func(node *zigbee.Node, p persistence.Section) {
		node.Depth = depth

		if p != nil {
			p.Set("Depth", uint64(node.Depth))
		}
	}
}

func updateReceived() nodeUpdate {
	return setReceived(time.Now())
}

func updateDiscovered() nodeUpdate {
	return setDiscovered(time.Now())
}

func setReceived(t time.Time) nodeUpdate {
	return func(node *zigbee.Node, p persistence.Section) {
		node.LastReceived = t

		if p != nil {
			converter.Store(p, "LastReceived", node.LastReceived, converter.TimeEncoder)
		}
	}
}

func setDiscovered(t time.Time) nodeUpdate {
	return func(node *zigbee.Node, p persistence.Section) {
		node.LastDiscovered = t

		if p != nil {
			converter.Store(p, "LastDiscovered", node.LastDiscovered, converter.TimeEncoder)
		}
	}
}
