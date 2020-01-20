package zstack

import (
	"context"
	"github.com/shimmeringbee/zigbee"
)

func (z *ZStack) RegisterAdapterEndpoint(ctx context.Context, endpoint uint8, appProfileId uint16, appDeviceId uint16, appDeviceVersion uint8, inClusters []zigbee.ZCLClusterID, outClusters []zigbee.ZCLClusterID) error {
	return nil
}

type AFRegisterReq struct {
	Endpoint         uint8
	AppProfileId     uint16
	AppDeviceId      uint16
	AppDeviceVersion uint8
	LatencyReq       uint8
	AppInClusters    []zigbee.ZCLClusterID `bclength:"8"`
	AppOutClusters   []zigbee.ZCLClusterID `bclength:"8"`
}

const AFRegisterReqID uint8 = 0x00

type AFRegisterResp GenericZStackStatus

const AFRegisterRespID uint8 = 0x00
