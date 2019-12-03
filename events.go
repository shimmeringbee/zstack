package zstack

import (
	"context"
	"errors"
	"github.com/shimmeringbee/zigbee"
)

func (z *ZStack) ReadEvent(ctx context.Context) (interface{}, error) {
	select {
	case event := <-z.events:
		return event, nil
	case <-ctx.Done():
		return nil, errors.New("context expired")
	}
}

type BasicDeviceEvent struct {
	NetworkAddress zigbee.NetworkAddress
	IEEEAddress    zigbee.IEEEAddress
}

type DeviceJoinEvent BasicDeviceEvent

type DeviceLeaveEvent BasicDeviceEvent
