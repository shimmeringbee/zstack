package zstack

import (
	"context"
	"errors"
	"github.com/shimmeringbee/zigbee"
)

func (z *ZStack) initialiseEvents() error {
	err, _ := z.subscriber.Subscribe(ZdoEndDeviceAnnceInd{}, z.handleZdoEndDeviceAnnceInd)
	if err != nil {
		return err
	}

	err, _ = z.subscriber.Subscribe(ZdoLeaveInd{}, z.handleZdoLeaveInd)
	if err != nil {
		return err
	}

	return nil
}

func (z *ZStack) handleZdoEndDeviceAnnceInd(u func(interface{}) error) {
	msg := ZdoEndDeviceAnnceInd{}
	err := u(&msg)

	if err == nil {
		z.events <- DeviceJoinEvent{
			NetworkAddress: msg.NetworkAddress,
			IEEEAddress:    msg.IEEEAddress,
		}
	}
}

func (z *ZStack) handleZdoLeaveInd(u func(interface{}) error) {
	msg := ZdoLeaveInd{}
	err := u(&msg)

	if err == nil {
		z.events <- DeviceLeaveEvent{
			NetworkAddress: msg.SourceAddress,
			IEEEAddress:    msg.IEEEAddress,
		}
	}
}

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
