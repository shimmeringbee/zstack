package zstack

import "github.com/shimmeringbee/zigbee"

type ZdoEndDeviceAnnceInd struct {
	SourceAddress  zigbee.NetworkAddress
	NetworkAddress zigbee.NetworkAddress
	IEEEAddress    zigbee.IEEEAddress
	Capabilities   uint8
}

const ZdoEndDeviceAnnceIndID uint8 = 0xc1

type ZdoLeaveInd struct {
	SourceAddress zigbee.NetworkAddress
	IEEEAddress   zigbee.IEEEAddress
	Request       uint8
	Remove        uint8
	Rejoin        uint8
}

const ZdoLeaveIndID uint8 = 0xc9

type ZdoTcDevInd struct {
	NetworkAddress zigbee.NetworkAddress
	IEEEAddress    zigbee.IEEEAddress
	ParentAddress  zigbee.NetworkAddress
}

const ZdoTcDevIndID uint8 = 0xca
