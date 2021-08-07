package zstack

import (
	"context"
	"fmt"
	"github.com/shimmeringbee/zigbee"
)

func (z *ZStack) AdapterNode() zigbee.Node {
	return zigbee.Node{
		IEEEAddress:    z.NetworkProperties.IEEEAddress,
		NetworkAddress: z.NetworkProperties.NetworkAddress,
		LogicalType:    zigbee.Coordinator,
	}
}

func (z *ZStack) GetAdapterIEEEAddress(ctx context.Context) (zigbee.IEEEAddress, error) {
	data, err := z.getAddressInfo(ctx)
	ieeeAddress := data.IEEEAddress

	return ieeeAddress, err
}

func (z *ZStack) GetAdapterNetworkAddress(ctx context.Context) (zigbee.NetworkAddress, error) {
	data, err := z.getAddressInfo(ctx)

	networkAddress := data.NetworkAddress
	return networkAddress, err
}

func (z *ZStack) getAddressInfo(ctx context.Context) (UtilGetDeviceInfoRequestReply, error) {
	resp := UtilGetDeviceInfoRequestReply{}

	if err := z.sem.Acquire(ctx, 1); err != nil {
		return resp, fmt.Errorf("failed to acquire semaphore: %w", err)
	}
	defer z.sem.Release(1)

	err := z.requestResponder.RequestResponse(ctx, UtilGetDeviceInfoRequest{}, &resp)
	return resp, err
}

type UtilGetDeviceInfoRequest struct{}

const UtilGetDeviceInfoRequestID uint8 = 0x00

type UtilGetDeviceInfoRequestReply struct {
	Status         uint8
	IEEEAddress    zigbee.IEEEAddress
	NetworkAddress zigbee.NetworkAddress
}

const UtilGetDeviceInfoRequestReplyID uint8 = 0x00
