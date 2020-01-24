package zstack

import (
	"context"
	"errors"
	"github.com/shimmeringbee/zigbee"
)

func (z *ZStack) QueryNodeEndpoints(ctx context.Context, networkAddress zigbee.NetworkAddress) ([]byte, error) {
	ch := make(chan ZdoActiveEpRsp)

	err, stop := z.subscriber.Subscribe(ZdoActiveEpRsp{}, func(unmarshal func(v interface{}) error) {
		msg := ZdoActiveEpRsp{}
		err := unmarshal(&msg)

		if err == nil && msg.OfInterestAddress == networkAddress {
			select {
			case ch <- msg:
			case <-ctx.Done():
			}
		}
	})

	defer stop()

	if err != nil {
		return []byte{}, err
	}

	request := ZdoActiveEpReq{
		DestinationAddress: networkAddress,
		OfInterestAddress:  networkAddress,
	}

	resp := ZdoActiveEpReqReply{}

	if err := z.requestResponder.RequestResponse(ctx, request, &resp); err != nil {
		return []byte{}, err
	}

	if resp.Status != ZSuccess {
		return []byte{}, ErrorZFailure
	}

	select {
	case response := <-ch:
		if response.Status == ZSuccess {
			return response.ActiveEndpoints, nil
		} else {
			return []byte{}, errors.New("error response received from node")
		}
	case <-ctx.Done():
		return []byte{}, errors.New("context expired while waiting for response from node")
	}
}

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
