package zstack

import (
	"context"
	"fmt"
	"github.com/shimmeringbee/zigbee"
)

func (z *ZStack) UnbindNodeFromController(ctx context.Context, nodeAddress zigbee.IEEEAddress, sourceEndpoint zigbee.Endpoint, destinationEndpoint zigbee.Endpoint, cluster zigbee.ClusterID) error {
	networkAddress, err := z.ResolveNodeNWKAddress(ctx, nodeAddress)
	if err != nil {
		return nil
	}

	if err := z.sem.Acquire(ctx, 1); err != nil {
		return fmt.Errorf("failed to acquire semaphore: %w", err)
	}
	defer z.sem.Release(1)

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
	SourceEndpoint         zigbee.Endpoint
	ClusterID              zigbee.ClusterID
	DestinationAddressMode uint8
	DestinationAddress     uint64
	DestinationEndpoint    zigbee.Endpoint
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
