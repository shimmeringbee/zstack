package zstack

import "github.com/shimmeringbee/zigbee"

type ZdoIEEEAddrReq struct {
	NetworkAddress zigbee.NetworkAddress
	ReqType        uint8
	StartIndex     uint8
}

const ZdoIEEEAddrReqID uint8 = 0x01

type ZdoIEEEAddrReqReply GenericZStackStatus

func (s ZdoIEEEAddrReqReply) WasSuccessful() bool {
	return s.Status == ZSuccess
}

const ZdoIEEEAddrReqReplyID uint8 = 0x01

type ZdoIEEEAddrRsp struct {
	Status            ZStackStatus
	IEEEAddress       zigbee.IEEEAddress
	NetworkAddress    zigbee.NetworkAddress
	StartIndex        uint8
	AssociatedDevices []zigbee.NetworkAddress `bclength:"8"`
}

func (s ZdoIEEEAddrRsp) WasSuccessful() bool {
	return s.Status == ZSuccess
}

const ZdoIEEEAddrRspID uint8 = 0x81
