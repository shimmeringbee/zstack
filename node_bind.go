package zstack

import (
	"context"
	"github.com/shimmeringbee/zigbee"
)

func (z *ZStack) BindNodeToController(ctx context.Context, nodeAddress zigbee.IEEEAddress, sourceEndpoint byte, destinationEndpoint byte, cluster zigbee.ZCLClusterID) error {
	networkAddress, err := z.ResolveNodeNWKAddress(ctx, nodeAddress)

	if err != nil {
		return nil
	}

	request := ZdoBindReq{
		TargetAddress:          networkAddress,
		SourceAddress:          nodeAddress,
		SourceEndpoint:         sourceEndpoint,
		ClusterID:              cluster,
		DestinationAddressMode: 0x02, // Network Address (16 bits)
		DestinationAddress:     uint64(0),
		DestinationEndpoint:    destinationEndpoint,
	}

	_, err = z.nodeRequest(ctx, &request, &ZdoBindReqReply{}, &ZdoBindRsp{}, func(i interface{}) bool {
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
