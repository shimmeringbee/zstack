package zstack

import (
	"context"
	"github.com/shimmeringbee/zigbee"
)

func (z *ZStack) QueryNodeEndpointDescription(ctx context.Context, ieeeAddress zigbee.IEEEAddress, endpoint byte) (zigbee.EndpointDescription, error) {
	networkAddress, err := z.ResolveNodeNWKAddress(ctx, ieeeAddress)

	if err != nil {
		return zigbee.EndpointDescription{}, err
	}

	request := ZdoSimpleDescReq{
		DestinationAddress: networkAddress,
		OfInterestAddress:  networkAddress,
		Endpoint:           endpoint,
	}

	resp, err := z.nodeRequest(ctx, &request, &ZdoSimpleDescReqReply{}, &ZdoSimpleDescRsp{}, func(i interface{}) bool {
		msg := i.(*ZdoSimpleDescRsp)
		return msg.OfInterestAddress == networkAddress && msg.Endpoint == endpoint
	})

	castResp, ok := resp.(*ZdoSimpleDescRsp)

	if ok {
		return zigbee.EndpointDescription{
			Endpoint:       castResp.Endpoint,
			ProfileID:      castResp.ProfileID,
			DeviceID:       castResp.DeviceID,
			DeviceVersion:  castResp.DeviceVersion,
			InClusterList:  castResp.InClusterList,
			OutClusterList: castResp.OutClusterList,
		}, nil
	} else {
		return zigbee.EndpointDescription{}, err
	}
}

type ZdoSimpleDescReq struct {
	DestinationAddress zigbee.NetworkAddress
	OfInterestAddress  zigbee.NetworkAddress
	Endpoint           byte
}

const ZdoSimpleDescReqID uint8 = 0x04

type ZdoSimpleDescReqReply GenericZStackStatus

func (r ZdoSimpleDescReqReply) WasSuccessful() bool {
	return r.Status == ZSuccess
}

const ZdoSimpleDescReqReplyID uint8 = 0x04

type ZdoSimpleDescRsp struct {
	SourceAddress     zigbee.NetworkAddress
	Status            ZStackStatus
	OfInterestAddress zigbee.NetworkAddress
	Length            uint8
	Endpoint          byte
	ProfileID         uint16
	DeviceID          uint16
	DeviceVersion     uint8
	InClusterList     []zigbee.ZCLClusterID `bclength:"8"`
	OutClusterList    []zigbee.ZCLClusterID `bclength:"8"`
}

func (r ZdoSimpleDescRsp) WasSuccessful() bool {
	return r.Status == ZSuccess
}

const ZdoSimpleDescRspID uint8 = 0x84
