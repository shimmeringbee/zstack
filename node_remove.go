package zstack

import (
	"context"
	"fmt"
	"github.com/shimmeringbee/zigbee"
)

func (z *ZStack) RequestNodeLeave(ctx context.Context, nodeAddress zigbee.IEEEAddress) error {
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

func (z *ZStack) ForceNodeLeave(ctx context.Context, nodeAddress zigbee.IEEEAddress) error {
	if z.removeNode(nodeAddress) {
		return nil
	}

	return fmt.Errorf("no node with address provided could be found: %v", nodeAddress)
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
