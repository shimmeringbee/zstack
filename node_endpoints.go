package zstack

import (
	"context"
	"fmt"
	"github.com/shimmeringbee/zigbee"
)

func (z *ZStack) QueryNodeEndpoints(ctx context.Context, ieeeAddress zigbee.IEEEAddress) ([]zigbee.Endpoint, error) {
	networkAddress, err := z.ResolveNodeNWKAddress(ctx, ieeeAddress)
	if err != nil {
		return []zigbee.Endpoint{}, err
	}

	if err := z.sem.Acquire(ctx, 1); err != nil {
		return nil, fmt.Errorf("failed to acquire semaphore: %w", err)
	}
	defer z.sem.Release(1)

	request := ZdoActiveEpReq{
		DestinationAddress: networkAddress,
		OfInterestAddress:  networkAddress,
	}

	resp, err := z.nodeRequest(ctx, &request, &ZdoActiveEpReqReply{}, &ZdoActiveEpRsp{}, func(i interface{}) bool {
		msg := i.(*ZdoActiveEpRsp)
		return msg.OfInterestAddress == networkAddress
	})

	castResp, ok := resp.(*ZdoActiveEpRsp)

	if ok {
		return castResp.ActiveEndpoints, err
	} else {
		return nil, err
	}
}

type ZdoActiveEpReq struct {
	DestinationAddress zigbee.NetworkAddress
	OfInterestAddress  zigbee.NetworkAddress
}

const ZdoActiveEpReqID uint8 = 0x05

type ZdoActiveEpReqReply GenericZStackStatus

func (r ZdoActiveEpReqReply) WasSuccessful() bool {
	return r.Status == ZSuccess
}

const ZdoActiveEpReqReplyID uint8 = 0x05

type ZdoActiveEpRsp struct {
	SourceAddress     zigbee.NetworkAddress
	Status            ZStackStatus
	OfInterestAddress zigbee.NetworkAddress
	ActiveEndpoints   []zigbee.Endpoint `bcsliceprefix:"8"`
}

func (r ZdoActiveEpRsp) WasSuccessful() bool {
	return r.Status == ZSuccess
}

const ZdoActiveEpRspID uint8 = 0x85
