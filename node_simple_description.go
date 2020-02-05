package zstack

import "github.com/shimmeringbee/zigbee"

type ZdoSimpleDescReq struct {
	DestinationAddress zigbee.NetworkAddress
	OfInterestAddress  zigbee.NetworkAddress
	Endpoint           byte
}

const ZdoSimpleDescReqID uint8 = 0x04

type ZdoSimpleDescReqReply GenericZStackStatus

const ZdoSimpleDescReqReplyID uint8 = 0x04

type ZdoSimpleDescRsp struct {
	SourceAddress     zigbee.NetworkAddress
	Status            ZStackStatus
	OfInterestAddress zigbee.NetworkAddress
	Length            uint8
	Endpoint          byte
	ProfileID         uint16
	DeviceID          uint16
	DeviceVersion     uint8
	InClusterList     []zigbee.ZCLClusterID `bclength:"8"`
	OutClusterList    []zigbee.ZCLClusterID `bclength:"8"`
}

const ZdoSimpleDescRspID uint8 = 0x84
