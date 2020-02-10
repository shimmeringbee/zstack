package zstack

import (
	"context"
	"github.com/shimmeringbee/zigbee"
)

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
	AssociatedDevices []zigbee.NetworkAddress `bclength:"8"`
}

func (s ZdoIEEEAddrRsp) WasSuccessful() bool {
	return s.Status == ZSuccess
}

const ZdoIEEEAddrRspID uint8 = 0x81
