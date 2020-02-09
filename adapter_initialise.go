package zstack

import (
	"context"
	"errors"
	"fmt"
	"github.com/shimmeringbee/zigbee"
)

func (z *ZStack) Initialise(ctx context.Context, nc zigbee.NetworkConfiguration) error {
	z.NetworkProperties.PANID = nc.PANID
	z.NetworkProperties.ExtendedPANID = nc.ExtendedPANID
	z.NetworkProperties.NetworkKey = nc.NetworkKey
	z.NetworkProperties.Channel = nc.Channel

	initFunctions := []func(context.Context) error{
		func(invokeCtx context.Context) error {
			return z.resetAdapter(invokeCtx, Soft)
		},
		func(invokeCtx context.Context) error {
			return z.writeNVRAM(invokeCtx, ZCDNVStartUpOption{StartOption: 0x03})
		},
		func(invokeCtx context.Context) error {
			return z.resetAdapter(invokeCtx, Soft)
		},
		func(invokeCtx context.Context) error {
			return z.writeNVRAM(invokeCtx, ZCDNVLogicalType{LogicalType: Coordinator})
		},
		func(invokeCtx context.Context) error {
			return z.resetAdapter(invokeCtx, Soft)
		},
		func(invokeCtx context.Context) error {
			return z.writeNVRAM(invokeCtx, ZCDNVSecurityMode{Enabled: 1})
		},
		func(invokeCtx context.Context) error {
			return z.writeNVRAM(invokeCtx, ZCDNVPreCfgKeysEnable{Enabled: 1})
		},
		func(invokeCtx context.Context) error {
			return z.writeNVRAM(invokeCtx, ZCDNVPreCfgKey{NetworkKey: nc.NetworkKey})
		},
		func(invokeCtx context.Context) error {
			return z.writeNVRAM(invokeCtx, ZCDNVZDODirectCB{Enabled: 1})
		},
		func(invokeCtx context.Context) error {
			channelBits := 1 << nc.Channel

			channelBytes := [4]byte{}
			channelBytes[0] = byte((channelBits >> 24) & 0xff)
			channelBytes[1] = byte((channelBits >> 16) & 0xff)
			channelBytes[2] = byte((channelBits >> 8) & 0xff)
			channelBytes[3] = byte((channelBits >> 0) & 0xff)

			return z.writeNVRAM(invokeCtx, ZCDNVChanList{Channels: channelBytes})
		},
		func(invokeCtx context.Context) error {
			return z.writeNVRAM(invokeCtx, ZCDNVPANID{PANID: nc.PANID})
		},
		func(invokeCtx context.Context) error {
			return z.writeNVRAM(invokeCtx, ZCDNVExtPANID{ExtendedPANID: nc.ExtendedPANID})
		},
		func(invokeCtx context.Context) error {
			return z.writeNVRAM(invokeCtx, ZCDNVUseDefaultTCLK{Enabled: 1})
		},
		func(invokeCtx context.Context) error {
			return z.writeNVRAM(invokeCtx, ZCDNVTCLKTableStart{
				Address:        zigbee.IEEEAddress(0xffffffffffffffff),
				NetworkKey:     zigbee.TCLinkKey,
				TXFrameCounter: 0,
				RXFrameCounter: 0,
			})
		},
		func(invokeCtx context.Context) error {
			address, err := z.GetAdapterIEEEAddress(ctx)

			if err != nil {
				return err
			}

			z.NetworkProperties.IEEEAddress = address

			return nil
		},
	}

	for _, f := range initFunctions {
		if err := Retry(ctx, DefaultZStackTimeout, DefaultZStackRetries, f); err != nil {
			return fmt.Errorf("failed during configuration and initialisation: %w", err)
		}
	}

	if err := z.startZigbeeStack(ctx); err != nil {
		return err
	}

	if err := Retry(ctx, DefaultZStackTimeout, DefaultZStackRetries, func(invokeCtx context.Context) error {
		address, err := z.GetAddressNetworkAddress(ctx)

		if err != nil {
			return err
		}

		z.NetworkProperties.NetworkAddress = address

		return nil
	}); err != nil {
		return err
	}

	if err := Retry(ctx, DefaultZStackTimeout, DefaultZStackRetries, func(invokeCtx context.Context) error {
		return z.DenyJoin(invokeCtx)
	}); err != nil {
		return err
	}

	z.startNetworkManager()
	z.startMessageReceiver()

	return nil
}

func (z *ZStack) startZigbeeStack(ctx context.Context) error {
	if err := Retry(ctx, DefaultZStackTimeout, DefaultZStackRetries, func(invokeCtx context.Context) error {
		return z.requestResponder.RequestResponse(invokeCtx, SAPIZBStartRequest{}, &SAPIZBStartRequestReply{})
	}); err != nil {
		return err
	}

	ch := make(chan bool, 1)
	defer close(ch)

	err, cancel := z.subscriber.Subscribe(&ZDOStateChangeInd{}, func(v interface{}) {
		stateChange := v.(*ZDOStateChangeInd)

		if stateChange.State == DeviceZBCoordinator {
			ch <- true
		}
	})
	defer cancel()

	if err != nil {
		return err
	}

	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return errors.New("context expired while waiting for adapter start up")
	}
}

type SAPIZBStartRequest struct{}

const SAPIZBStartRequestID uint8 = 0x00

type SAPIZBStartRequestReply struct{}

const SAPIZBStartRequestReplyID uint8 = 0x00

type ZBStartStatus uint8

const (
	ZBSuccess ZBStartStatus = 0x00
	ZBInit    ZBStartStatus = 0x22
)

type ZDOState uint8

const (
	DeviceCoordinatorStarting ZDOState = 0x08
	DeviceZBCoordinator       ZDOState = 0x09
)

type ZDOStateChangeInd struct {
	State ZDOState
}

const ZDOStateChangeIndID uint8 = 0xc0
