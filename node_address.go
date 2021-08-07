package zstack

import (
	"context"
	"fmt"
	"github.com/shimmeringbee/logwrap"
	"github.com/shimmeringbee/zigbee"
)

func (z *ZStack) ResolveNodeIEEEAddress(ctx context.Context, address zigbee.NetworkAddress) (zigbee.IEEEAddress, error) {
	if node, found := z.nodeTable.getByNetwork(address); found {
		return node.IEEEAddress, nil
	} else {
		z.logger.LogDebug(ctx, "Asked to resolve Network Address to IEEE Address, but not present in node table, querying adapter.", logwrap.Datum("NetworkAddress", address))
		return z.QueryNodeIEEEAddress(ctx, address)
	}
}

func (z *ZStack) ResolveNodeNWKAddress(ctx context.Context, address zigbee.IEEEAddress) (zigbee.NetworkAddress, error) {
	if node, found := z.nodeTable.getByIEEE(address); found {
		return node.NetworkAddress, nil
	} else {
		z.logger.LogDebug(ctx, "Asked to resolve IEEE Address to Network Address, but not present in node table, querying adapter.", logwrap.Datum("IEEEAddress", address.String()))
		return z.QueryNodeNWKAddress(ctx, address)
	}
}

func (z *ZStack) QueryNodeIEEEAddress(ctx context.Context, address zigbee.NetworkAddress) (zigbee.IEEEAddress, error) {
	if err := z.sem.Acquire(ctx, 1); err != nil {
		return zigbee.EmptyIEEEAddress, fmt.Errorf("failed to acquire semaphore: %w", err)
	}
	defer z.sem.Release(1)

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
		z.logger.LogError(ctx, "Failed to query adapter for IEEE Address.", logwrap.Datum("NetworkAddress", address), logwrap.Err(err))
		return zigbee.EmptyIEEEAddress, err
	}
}

func (z *ZStack) QueryNodeNWKAddress(ctx context.Context, address zigbee.IEEEAddress) (zigbee.NetworkAddress, error) {
	if err := z.sem.Acquire(ctx, 1); err != nil {
		return zigbee.NetworkAddress(0x0), fmt.Errorf("failed to acquire semaphore: %w", err)
	}
	defer z.sem.Release(1)

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
		z.logger.LogError(ctx, "Failed to query adapter for Network Address.", logwrap.Datum("IEEEAddress", address.String()), logwrap.Err(err))
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
