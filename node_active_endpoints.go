package zstack

import "github.com/shimmeringbee/zigbee"

type ZdoActiveEpReq struct {
	DestinationAddress zigbee.NetworkAddress
	OfInterestAddress  zigbee.NetworkAddress
}

const ZdoActiveEpReqID uint8 = 0x05

type ZdoActiveEpReqReply GenericZStackStatus

const ZdoActiveEpReqReplyID uint8 = 0x05

type ZdoActiveEpRsp struct {
	SourceAddress     zigbee.NetworkAddress
	Status            ZStackStatus
	OfInterestAddress zigbee.NetworkAddress
	ActiveEndpoints   []byte `bclength:"8"`
}

const ZdoActiveEpRspID uint8 = 0x85
