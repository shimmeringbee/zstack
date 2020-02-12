package zstack

import (
	"context"
	"github.com/shimmeringbee/zigbee"
)

func (z *ZStack) RegisterAdapterEndpoint(ctx context.Context, endpoint zigbee.Endpoint, appProfileId zigbee.ProfileID, appDeviceId uint16, appDeviceVersion uint8, inClusters []zigbee.ClusterID, outClusters []zigbee.ClusterID) error {
	request := AFRegister{
		Endpoint:         endpoint,
		AppProfileId:     appProfileId,
		AppDeviceId:      appDeviceId,
		AppDeviceVersion: appDeviceVersion,
		LatencyReq:       0x00, // No latency, no other valid option for Zigbee
		AppInClusters:    inClusters,
		AppOutClusters:   outClusters,
	}

	resp := AFRegisterReply{}

	if err := z.requestResponder.RequestResponse(ctx, request, &resp); err != nil {
		return err
	}

	if resp.Status != ZSuccess {
		return ErrorZFailure
	}

	return nil
}

type AFRegister struct {
	Endpoint         zigbee.Endpoint
	AppProfileId     zigbee.ProfileID
	AppDeviceId      uint16
	AppDeviceVersion uint8
	LatencyReq       uint8
	AppInClusters    []zigbee.ClusterID `bclength:"8"`
	AppOutClusters   []zigbee.ClusterID `bclength:"8"`
}

const AFRegisterID uint8 = 0x00

type AFRegisterReply GenericZStackStatus

const AFRegisterReplyID uint8 = 0x00
