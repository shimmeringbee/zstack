package zstack

import (
	"context"
	"github.com/shimmeringbee/zigbee"
)

func (z *ZStack) UnbindNodeFromController(ctx context.Context, nodeAddress zigbee.IEEEAddress, sourceEndpoint byte, destinationEndpoint byte, cluster zigbee.ZCLClusterID) error {
	networkAddress, err := z.ResolveNodeNWKAddress(ctx, nodeAddress)

	if err != nil {
		return nil
	}

	request := ZdoUnbindReq{
		TargetAddress:          networkAddress,
		SourceAddress:          nodeAddress,
		SourceEndpoint:         sourceEndpoint,
		ClusterID:              cluster,
		DestinationAddressMode: 0x02, // Network Address (16 bits)
		DestinationAddress:     uint64(0),
		DestinationEndpoint:    destinationEndpoint,
	}

	_, err = z.nodeRequest(ctx, &request, &ZdoUnbindReqReply{}, &ZdoUnbindRsp{}, func(i interface{}) bool {
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
