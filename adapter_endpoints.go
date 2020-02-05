package zstack

import (
	"context"
	"github.com/shimmeringbee/zigbee"
)

func (z *ZStack) RegisterAdapterEndpoint(ctx context.Context, endpoint uint8, appProfileId uint16, appDeviceId uint16, appDeviceVersion uint8, inClusters []zigbee.ZCLClusterID, outClusters []zigbee.ZCLClusterID) error {
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
	Endpoint         uint8
	AppProfileId     uint16
	AppDeviceId      uint16
	AppDeviceVersion uint8
	LatencyReq       uint8
	AppInClusters    []zigbee.ZCLClusterID `bclength:"8"`
	AppOutClusters   []zigbee.ZCLClusterID `bclength:"8"`
}

const AFRegisterID uint8 = 0x00

type AFRegisterReply GenericZStackStatus

const AFRegisterReplyID uint8 = 0x00
