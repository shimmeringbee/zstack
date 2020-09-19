package zstack

import (
	"context"
	"github.com/shimmeringbee/zigbee"
)

func (z *ZStack) ResolveNodeIEEEAddress(ctx context.Context, address zigbee.NetworkAddress) (zigbee.IEEEAddress, error) {
	node, found := z.nodeTable.getByNetwork(address)

	if found {
		return node.IEEEAddress, nil
	}

	return z.QueryNodeIEEEAddress(ctx, address)
}

func (z *ZStack) ResolveNodeNWKAddress(ctx context.Context, address zigbee.IEEEAddress) (zigbee.NetworkAddress, error) {
	node, found := z.nodeTable.getByIEEE(address)

	if found {
		return node.NetworkAddress, nil
	}

	return z.QueryNodeNWKAddress(ctx, address)
}

func (z *ZStack) QueryNodeIEEEAddress(ctx context.Context, address zigbee.NetworkAddress) (zigbee.IEEEAddress, error) {
	request := ZdoIEEEAddrReq{
		NetworkAddress: address,
		ReqType:        0x00,
		StartIndex:     0x00,
	}

	resp, err := z.nodeRequest(ctx, &request, &ZdoIEEEAddrReqReply{}, &ZdoIEEEAddrRsp{}, func(i interface{}) bool {
		msg := i.(*ZdoIEEEAddrRsp)
		return msg.NetworkAddress == address
	})

	castResp, ok := resp.(*ZdoIEEEAddrRsp)

	if ok {
		return castResp.IEEEAddress, nil
	} else {
		return zigbee.EmptyIEEEAddress, err
	}
}

func (z *ZStack) QueryNodeNWKAddress(ctx context.Context, address zigbee.IEEEAddress) (zigbee.NetworkAddress, error) {
	request := ZdoNWKAddrReq{
		IEEEAddress: address,
		ReqType:     0x00,
		StartIndex:  0x00,
	}

	resp, err := z.nodeRequest(ctx, &request, &ZdoNWKAddrReqReply{}, &ZdoNWKAddrRsp{}, func(i interface{}) bool {
		msg := i.(*ZdoNWKAddrRsp)
		return msg.IEEEAddress == address
	})

	castResp, ok := resp.(*ZdoNWKAddrRsp)

	if ok {
		return castResp.NetworkAddress, nil
	} else {
		return zigbee.NetworkAddress(0x0), err
	}
}

type ZdoIEEEAddrReq struct {
	NetworkAddress zigbee.NetworkAddress
	ReqType        uint8
	StartIndex     uint8
}

const ZdoIEEEAddrReqID uint8 = 0x01

type ZdoIEEEAddrReqReply GenericZStackStatus

func (s ZdoIEEEAddrReqReply) WasSuccessful() bool {
	return s.Status == ZSuccess
}

const ZdoIEEEAddrReqReplyID uint8 = 0x01

type ZdoIEEEAddrRsp struct {
	Status            ZStackStatus
	IEEEAddress       zigbee.IEEEAddress
	NetworkAddress    zigbee.NetworkAddress
	StartIndex        uint8
	AssociatedDevices []zigbee.NetworkAddress `bcsliceprefix:"8"`
}

func (s ZdoIEEEAddrRsp) WasSuccessful() bool {
	return s.Status == ZSuccess
}

const ZdoIEEEAddrRspID uint8 = 0x81

type ZdoNWKAddrReq struct {
	IEEEAddress zigbee.IEEEAddress
	ReqType     uint8
	StartIndex  uint8
}

const ZdoNWKAddrReqID uint8 = 0x00

type ZdoNWKAddrReqReply GenericZStackStatus

func (s ZdoNWKAddrReqReply) WasSuccessful() bool {
	return s.Status == ZSuccess
}

const ZdoNWKAddrReqReplyID uint8 = 0x00

type ZdoNWKAddrRsp struct {
	Status            ZStackStatus
	IEEEAddress       zigbee.IEEEAddress
	NetworkAddress    zigbee.NetworkAddress
	StartIndex        uint8
	AssociatedDevices []zigbee.NetworkAddress `bcsliceprefix:"8"`
}

func (s ZdoNWKAddrRsp) WasSuccessful() bool {
	return s.Status == ZSuccess
}

const ZdoNWKAddrRspID uint8 = 0x80
