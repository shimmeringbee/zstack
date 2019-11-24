package zstack

import (
	"context"
	"github.com/shimmeringbee/zigbee"
)

func (z *ZStack) GetDeviceIEEEAddress(ctx context.Context) (zigbee.IEEEAddress, error) {
	return z.getDeviceInfo(ctx, IEEEAddress)
}

func (z *ZStack) GetDeviceNetworkAddress(ctx context.Context) (zigbee.NetworkAddress, error) {
	data, err := z.getDeviceInfo(ctx, NetworkAddress)

	networkAddress := zigbee.NetworkAddress{}
	copy(networkAddress[0:2], data[0:2])

	return networkAddress, err
}

func (z *ZStack) getDeviceInfo(ctx context.Context, parameter DeviceInfoParameter) ([8]byte, error) {
	resp := SAPIZBGetDeviceInfoResp{}

	if err := z.RequestResponder.RequestResponse(ctx, SAPIZBGetDeviceInfoReq{Parameter:parameter}, &resp); err != nil {
		return [8]byte{}, err
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
	Value     [8]byte
}

const SAPIZBGetDeviceInfoRespID uint8 = 0x06
