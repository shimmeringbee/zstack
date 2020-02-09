package zstack

import "github.com/shimmeringbee/zigbee"

type ZdoUnbindReq struct {
	TargetAddress          zigbee.NetworkAddress
	SourceAddress          zigbee.IEEEAddress
	SourceEndpoint         byte
	ClusterID              zigbee.ZCLClusterID
	DestinationAddressMode uint8
	DestinationAddress     uint64
	DestinationEndpoint    byte
}

const ZdoUnbindReqID uint8 = 0x22

type ZdoUnbindReqReply GenericZStackStatus

func (r ZdoUnbindReqReply) WasSuccessful() bool {
	return r.Status == ZSuccess
}

const ZdoUnbindReqReplyID uint8 = 0x22

type ZdoUnbindRsp struct {
	SourceAddress zigbee.NetworkAddress
	Status        ZStackStatus
}

func (r ZdoUnbindRsp) WasSuccessful() bool {
	return r.Status == ZSuccess
}

const ZdoUnbindRspID uint8 = 0xa2
