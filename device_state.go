package zstack

import "github.com/shimmeringbee/zigbee"

type ZdoEndDeviceAnnceIndCapabilities struct {
	AddressAllocated   bool  `bcfieldwidth:"1"`
	SecurityCapability bool  `bcfieldwidth:"1"`
	Reserved           uint8 `bcfieldwidth:"2"`
	ReceiveOnIdle      bool  `bcfieldwidth:"1"`
	PowerSource        bool  `bcfieldwidth:"1"`
	Router             bool  `bcfieldwidth:"1"`
	AltPANController   bool  `bcfieldwidth:"1"`
}

type ZdoEndDeviceAnnceInd struct {
	SourceAddress  zigbee.NetworkAddress
	NetworkAddress zigbee.NetworkAddress
	IEEEAddress    zigbee.IEEEAddress
	Capabilities   ZdoEndDeviceAnnceIndCapabilities
}

const ZdoEndDeviceAnnceIndID uint8 = 0xc1

type ZdoLeaveInd struct {
	SourceAddress zigbee.NetworkAddress
	IEEEAddress   zigbee.IEEEAddress
	Request       bool
	Remove        bool
	Rejoin        bool
}

const ZdoLeaveIndID uint8 = 0xc9

type ZdoTcDevInd struct {
	NetworkAddress zigbee.NetworkAddress
	IEEEAddress    zigbee.IEEEAddress
	ParentAddress  zigbee.NetworkAddress
}

const ZdoTcDevIndID uint8 = 0xca
