package zstack

import (
	"context"
	"github.com/shimmeringbee/zigbee"
)

func (z *ZStack) BindToNode(ctx context.Context, networkAddress zigbee.NetworkAddress, sourceAddress zigbee.IEEEAddress, sourceEndpoint byte, destinationAddress zigbee.IEEEAddress, destinationEndpoint byte, cluster zigbee.ZCLClusterID) error {
	request := ZdoBindReq{
		TargetAddress:          networkAddress,
		SourceAddress:          sourceAddress,
		SourceEndpoint:         sourceEndpoint,
		ClusterID:              cluster,
		DestinationAddressMode: 0x03,
		DestinationAddress:     uint64(destinationAddress),
		DestinationEndpoint:    destinationEndpoint,
	}

	_, err := z.nodeRequest(ctx, &request, &ZdoBindReqReply{}, &ZdoBindRsp{}, func(i interface{}) bool {
		msg := i.(*ZdoBindRsp)
		return msg.SourceAddress == networkAddress
	})

	return err
}

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
