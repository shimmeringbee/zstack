package zstack

import (
	"context"
	"fmt"
	"github.com/shimmeringbee/zigbee"
)

func (z *ZStack) Initialise(ctx context.Context, nc zigbee.NetworkConfiguration) error {
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
				Address:        zigbee.IEEEAddress{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
				NetworkKey:     zigbee.TCLinkKey,
				TXFrameCounter: 0,
				RXFrameCounter: 0,
			})
		},
	}

	for _, f := range initFunctions {
		if err := Retry(ctx, DefaultZStackTimeout, DefaultZStackRetries, f); err != nil {
			return fmt.Errorf("failed during configuration and initialisation: %w", err)
		}
	}

	return nil
}

func (z *ZStack) startZigbeeStack(ctx context.Context) error {
	if err := Retry(ctx, DefaultZStackTimeout, DefaultZStackRetries, func(invokeCtx context.Context) error {
		return z.RequestResponder.RequestResponse(invokeCtx, SAPIZBStartRequest{}, &SAPIZBStartResponse{})
	}); err != nil {
		return err
	}

	confirmation := SAPIZBStartConfirm{}
	if err := z.Awaiter.Await(ctx, &confirmation); err != nil {
		return err
	}

	if confirmation.Status != ZBSuccess {
		return fmt.Errorf("failed to start application on zigbee adapter: status %d", confirmation.Status)
	}

	return nil
}

type SAPIZBStartRequest struct{}

const SAPIZBStartRequestID uint8 = 0x00

type SAPIZBStartResponse struct{}

const SAPIZBStartResponseID uint8 = 0x00

type ZBStartStatus uint8

const (
	ZBSuccess ZBStartStatus = 0x00
	ZBInit    ZBStartStatus = 0x22
)

type SAPIZBStartConfirm struct {
	Status ZBStartStatus
}

const SAPIZBStartConfirmID uint8 = 0x80
