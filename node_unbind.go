package zstack

import (
	"context"
	"github.com/shimmeringbee/zigbee"
)

func (z *ZStack) UnbindToNode(ctx context.Context, networkAddress zigbee.NetworkAddress, sourceAddress zigbee.IEEEAddress, sourceEndpoint byte, destinationAddress zigbee.IEEEAddress, destinationEndpoint byte, cluster zigbee.ZCLClusterID) error {
	request := ZdoUnbindReq{
		TargetAddress:          networkAddress,
		SourceAddress:          sourceAddress,
		SourceEndpoint:         sourceEndpoint,
		ClusterID:              cluster,
		DestinationAddressMode: 0x03,
		DestinationAddress:     uint64(destinationAddress),
		DestinationEndpoint:    destinationEndpoint,
	}

	_, err := z.nodeRequest(ctx, &request, &ZdoUnbindReqReply{}, &ZdoUnbindRsp{}, func(i interface{}) bool {
		msg := i.(*ZdoUnbindRsp)
		return msg.SourceAddress == networkAddress
	})

	return err
}

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
