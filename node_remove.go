package zstack

import (
	"context"
	"github.com/shimmeringbee/zigbee"
)

func (z *ZStack) RemoveNode(ctx context.Context, nodeAddress zigbee.IEEEAddress) error {
	networkAddress, err := z.ResolveNodeNWKAddress(ctx, nodeAddress)

	if err != nil {
		return nil
	}

	request := ZdoMgmtLeaveReq{
		NetworkAddress: networkAddress,
		IEEEAddress:    nodeAddress,
		RemoveChildren: false,
	}

	_, err = z.nodeRequest(ctx, &request, &ZdoMgmtLeaveReqReply{}, &ZdoMgmtLeaveRsp{}, func(i interface{}) bool {
		msg := i.(*ZdoMgmtLeaveRsp)
		return msg.SourceAddress == networkAddress
	})

	return err
}

type ZdoMgmtLeaveReq struct {
	NetworkAddress zigbee.NetworkAddress
	IEEEAddress    zigbee.IEEEAddress
	RemoveChildren bool
}

const ZdoMgmtLeaveReqID uint8 = 0x34

type ZdoMgmtLeaveReqReply GenericZStackStatus

func (r ZdoMgmtLeaveReqReply) WasSuccessful() bool {
	return r.Status == ZSuccess
}

const ZdoMgmtLeaveReqReplyID uint8 = 0x34

type ZdoMgmtLeaveRsp struct {
	SourceAddress zigbee.NetworkAddress
	Status        ZStackStatus
}

func (r ZdoMgmtLeaveRsp) WasSuccessful() bool {
	return r.Status == ZSuccess
}

const ZdoMgmtLeaveRspID uint8 = 0xb4
