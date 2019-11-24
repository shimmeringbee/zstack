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
	response := SAPIZBPermitJoiningResponse{}

	if err := z.requestResponder.RequestResponse(ctx, SAPIZBPermitJoiningRequest{
		Destination: address,
		Timeout:     timeout,
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

type SAPIZBPermitJoiningRequest struct {
	Destination zigbee.NetworkAddress
	Timeout     uint8
}

const SAPIZBPermitJoiningRequestID uint8 = 0x08

type SAPIZBPermitJoiningResponse GenericZStackStatus

const SAPIZBPermitJoiningResponseID uint8 = 0x08
