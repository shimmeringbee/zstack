package zstack

import (
	"context"
	"github.com/shimmeringbee/zigbee"
)

func (z *ZStack) GetAdapterIEEEAddress(ctx context.Context) (zigbee.IEEEAddress, error) {
	data, err := z.getAddressInfo(ctx, IEEEAddress)

	ieeeAddress := zigbee.IEEEAddress(data)

	return ieeeAddress, err
}

func (z *ZStack) GetAddressNetworkAddress(ctx context.Context) (zigbee.NetworkAddress, error) {
	data, err := z.getAddressInfo(ctx, NetworkAddress)

	networkAddress := zigbee.NetworkAddress(data & 0xffff)

	return networkAddress, err
}

func (z *ZStack) getAddressInfo(ctx context.Context, parameter DeviceInfoParameter) (uint64, error) {
	resp := SAPIZBGetDeviceInfoResp{}

	if err := z.requestResponder.RequestResponse(ctx, SAPIZBGetDeviceInfoReq{Parameter: parameter}, &resp); err != nil {
		return 0, err
	}

	return resp.Value, nil
}

type DeviceInfoParameter uint8

const (
	State                  DeviceInfoParameter = 0x00
	IEEEAddress            DeviceInfoParameter = 0x01
	NetworkAddress         DeviceInfoParameter = 0x02
	ParentNetworkAddress   DeviceInfoParameter = 0x03
	ParentIEEEAddress      DeviceInfoParameter = 0x04
	OperatingChannel       DeviceInfoParameter = 0x05
	OperatingPANID         DeviceInfoParameter = 0x06
	OperatingExtendedPANID DeviceInfoParameter = 0x07
)

type SAPIZBGetDeviceInfoReq struct {
	Parameter DeviceInfoParameter
}

const SAPIZBGetDeviceInfoReqID uint8 = 0x06

type SAPIZBGetDeviceInfoResp struct {
	Parameter DeviceInfoParameter
	Value     uint64
}

const SAPIZBGetDeviceInfoRespID uint8 = 0x06
