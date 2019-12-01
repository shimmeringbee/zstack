package zstack

import "github.com/shimmeringbee/zigbee"

func (z *ZStack) addressUpdate(ieee zigbee.IEEEAddress, na zigbee.NetworkAddress) {
	z.addressCache[ieee] = na
}

func (z *ZStack) addressRemove(ieee zigbee.IEEEAddress) {
	delete(z.addressCache, ieee)
}

func (z *ZStack) addressLookup(ieee zigbee.IEEEAddress) (zigbee.NetworkAddress, bool) {
	address, found := z.addressCache[ieee]
	return address, found
}
