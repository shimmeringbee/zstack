package zstack

import (
	"context"
	"fmt"
	"github.com/shimmeringbee/zigbee"
)

func (z *ZStack) QueryNodeEndpointDescription(ctx context.Context, ieeeAddress zigbee.IEEEAddress, endpoint zigbee.Endpoint) (zigbee.EndpointDescription, error) {
	networkAddress, err := z.ResolveNodeNWKAddress(ctx, ieeeAddress)
	if err != nil {
		return zigbee.EndpointDescription{}, err
	}

	if err := z.sem.Acquire(ctx, 1); err != nil {
		return zigbee.EndpointDescription{}, fmt.Errorf("failed to acquire semaphore: %w", err)
	}
	defer z.sem.Release(1)

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
	Endpoint           zigbee.Endpoint
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
	Endpoint          zigbee.Endpoint
	ProfileID         zigbee.ProfileID
	DeviceID          uint16
	DeviceVersion     uint8
	InClusterList     []zigbee.ClusterID `bcsliceprefix:"8"`
	OutClusterList    []zigbee.ClusterID `bcsliceprefix:"8"`
}

func (r ZdoSimpleDescRsp) WasSuccessful() bool {
	return r.Status == ZSuccess
}

const ZdoSimpleDescRspID uint8 = 0x84
