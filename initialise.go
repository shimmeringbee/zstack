package zstack

import (
	"context"
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
			return err
		}
	}

	return nil
}
