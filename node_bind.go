package zstack

import "github.com/shimmeringbee/zigbee"

type ZdoBindReq struct {
	TargetAddress          zigbee.NetworkAddress
	SourceAddress          zigbee.IEEEAddress
	SourceEndpoint         byte
	ClusterID              zigbee.ZCLClusterID
	DestinationAddressMode uint8
	DestinationAddress     uint64
	DestinationEndpoint    byte
}

const ZdoBindReqID uint8 = 0x21

type ZdoBindReqReply GenericZStackStatus

func (r ZdoBindReqReply) WasSuccessful() bool {
	return r.Status == ZSuccess
}

const ZdoBindReqReplyID uint8 = 0x21

type ZdoBindRsp struct {
	SourceAddress zigbee.NetworkAddress
	Status        ZStackStatus
}

func (r ZdoBindRsp) WasSuccessful() bool {
	return r.Status == ZSuccess
}

const ZdoBindRspID uint8 = 0xa1
