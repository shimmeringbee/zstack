package zstack

import (
	"context"
	"github.com/shimmeringbee/zigbee"
)

func (z *ZStack) QueryNodeEndpoints(ctx context.Context, ieeeAddress zigbee.IEEEAddress) ([]byte, error) {
	networkAddress, err := z.ResolveNodeNWKAddress(ctx, ieeeAddress)

	if err != nil {
		return []byte{}, err
	}

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
	ActiveEndpoints   []byte `bclength:"8"`
}

func (r ZdoActiveEpRsp) WasSuccessful() bool {
	return r.Status == ZSuccess
}

const ZdoActiveEpRspID uint8 = 0x85
