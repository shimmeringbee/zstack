package zstack

import (
	"context"
	"fmt"
	"github.com/shimmeringbee/zigbee"
)

func (z *ZStack) PermitJoin(ctx context.Context, allRouters bool) error {
	if allRouters {
		return z.sendJoin(ctx, zigbee.BroadcastRoutersCoordinators, JoiningOn, OnAllRouters)
	} else {
		return z.sendJoin(ctx, z.NetworkProperties.NetworkAddress, JoiningOn, OnCoordinator)
	}
}

func (z *ZStack) DenyJoin(ctx context.Context) error {
	return z.sendJoin(ctx, zigbee.BroadcastRoutersCoordinators, JoiningOff, Off)
}

func (z *ZStack) sendJoin(ctx context.Context, address zigbee.NetworkAddress, timeout uint8, newState JoinState) error {
	response := ZDOMgmtPermitJoinRequestReply{}

	if err := z.requestResponder.RequestResponse(ctx, ZDOMgmtPermitJoinRequest{
		Destination:    address,
		Duration:       timeout,
		TCSignificance: 0x00,
	}, &response); err != nil {
		return err
	}

	if response.Status != ZSuccess {
		return fmt.Errorf("adapter rejected permit join state change: state=%v", response.Status)
	}

	z.NetworkProperties.JoinState = newState

	return nil
}

const (
	JoiningOff uint8 = 0x00
	JoiningOn  uint8 = 0xff
)

type ZDOMgmtPermitJoinRequest struct {
	Destination    zigbee.NetworkAddress
	Duration       uint8
	TCSignificance uint8
}

const ZDOMgmtPermitJoinRequestID = 0x36

type ZDOMgmtPermitJoinRequestReply GenericZStackStatus

const ZDOMgmtPermitJoinRequestReplyID uint8 = 0x36
