package zstack

import (
	"context"
	"errors"
	"fmt"
	"github.com/shimmeringbee/logwrap"
	"github.com/shimmeringbee/retry"
	"github.com/shimmeringbee/zigbee"
	"golang.org/x/sync/semaphore"
	"reflect"
)

func (z *ZStack) Initialise(pctx context.Context, nc zigbee.NetworkConfiguration) error {
	z.NetworkProperties.PANID = nc.PANID
	z.NetworkProperties.ExtendedPANID = nc.ExtendedPANID
	z.NetworkProperties.NetworkKey = nc.NetworkKey
	z.NetworkProperties.Channel = nc.Channel

	ctx, segmentEnd := z.logger.Segment(pctx, "Adapter Initialise.")
	defer segmentEnd()

	z.logger.LogInfo(ctx, "Restarting adapter.")
	version, err := z.waitForAdapterReset(ctx)
	if err != nil {
		return err
	}

	if version.IsV3() {
		z.sem = semaphore.NewWeighted(16)
	} else {
		z.sem = semaphore.NewWeighted(2)
	}

	z.logger.LogInfo(ctx, "Verifying existing network configuration.")
	if valid, err := z.verifyAdapterNetworkConfig(ctx, version); err != nil {
		return err
	} else if !valid {
		z.logger.LogWarn(ctx, "Adapter network configuration is invalid, resetting adapter.")
		if err := z.wipeAdapter(ctx); err != nil {
			return err
		}

		z.logger.LogInfo(ctx, "Setting adapter to coordinator.")
		if err := z.makeCoordinator(ctx); err != nil {
			return err
		}

		z.logger.LogInfo(ctx, "Configuring adapter.")
		if err := z.configureNetwork(ctx, version); err != nil {
			return err
		}
	}

	z.logger.LogInfo(ctx, "Starting Zigbee stack.")
	if err := z.startZigbeeStack(ctx, version); err != nil {
		return err
	}

	z.logger.LogInfo(ctx, "Fetching adapter IEEE and Network addresses.")
	if err := z.retrieveAdapterAddresses(ctx); err != nil {
		return err
	}

	z.logger.LogInfo(ctx, "Enforcing denial of network joins.")
	if err := z.DenyJoin(ctx); err != nil {
		return err
	}

	z.startNetworkManager()
	z.startMessageReceiver()

	return nil
}

func (z *ZStack) waitForAdapterReset(ctx context.Context) (Version, error) {
	retVersion := Version{}

	err := retry.Retry(ctx, DefaultZStackTimeout, 18, func(invokeCtx context.Context) error {
		version, err := z.resetAdapter(invokeCtx, Soft)
		retVersion = version
		return err
	})

	return retVersion, err
}

func (z *ZStack) verifyAdapterNetworkConfig(ctx context.Context, version Version) (bool, error) {
	configToVerify := []interface{}{
		&ZCDNVLogicalType{LogicalType: zigbee.Coordinator},
		&ZCDNVPANID{PANID: z.NetworkProperties.PANID},
		&ZCDNVExtPANID{ExtendedPANID: z.NetworkProperties.ExtendedPANID},
		&ZCDNVChanList{Channels: channelToBits(z.NetworkProperties.Channel)},
	}

	for _, expectedConfig := range configToVerify {
		configType := reflect.TypeOf(expectedConfig).Elem()
		actualConfig := reflect.New(configType).Interface()

		if err := z.readNVRAM(ctx, actualConfig); err != nil {
			return false, err
		}

		if !reflect.DeepEqual(expectedConfig, actualConfig) {
			return false, nil
		}
	}

	return true, nil
}

func (z *ZStack) wipeAdapter(ctx context.Context) error {
	return retryFunctions(ctx, []func(context.Context) error{
		func(invokeCtx context.Context) error {
			return z.writeNVRAM(invokeCtx, ZCDNVStartUpOption{StartOption: 0x03})
		},
		func(invokeCtx context.Context) error {
			_, err := z.resetAdapter(invokeCtx, Soft)
			return err
		},
	})
}

func (z *ZStack) makeCoordinator(ctx context.Context) error {
	return retryFunctions(ctx, []func(context.Context) error{
		func(invokeCtx context.Context) error {
			return z.writeNVRAM(invokeCtx, ZCDNVLogicalType{LogicalType: zigbee.Coordinator})
		},
		func(invokeCtx context.Context) error {
			_, err := z.resetAdapter(invokeCtx, Soft)
			return err
		},
	})
}

func (z *ZStack) configureNetwork(ctx context.Context, version Version) error {
	channelBits := channelToBits(z.NetworkProperties.Channel)

	if err := retryFunctions(ctx, []func(context.Context) error{
		func(invokeCtx context.Context) error {
			z.logger.LogDebug(ctx, "Adapater Initialisation: Enabling preconfigured keys.")
			return z.writeNVRAM(invokeCtx, ZCDNVPreCfgKeysEnable{Enabled: 1})
		},
		func(invokeCtx context.Context) error {
			z.logger.LogDebug(ctx, "Adapater Initialisation: Configuring network key.")
			return z.writeNVRAM(invokeCtx, ZCDNVPreCfgKey{NetworkKey: z.NetworkProperties.NetworkKey})
		},
		func(invokeCtx context.Context) error {
			z.logger.LogDebug(ctx, "Adapater Initialisation: Enable ZDO callbacks.")
			return z.writeNVRAM(invokeCtx, ZCDNVZDODirectCB{Enabled: 1})
		},
		func(invokeCtx context.Context) error {
			z.logger.LogDebug(ctx, "Adapater Initialisation: Configuring network channel.")
			return z.writeNVRAM(invokeCtx, ZCDNVChanList{Channels: channelBits})
		},
		func(invokeCtx context.Context) error {
			z.logger.LogDebug(ctx, "Adapater Initialisation: Configuring network PANID.")
			return z.writeNVRAM(invokeCtx, ZCDNVPANID{PANID: z.NetworkProperties.PANID})
		},
		func(invokeCtx context.Context) error {
			z.logger.LogDebug(ctx, "Adapater Initialisation: Configuring network extended PANID.")
			return z.writeNVRAM(invokeCtx, ZCDNVExtPANID{ExtendedPANID: z.NetworkProperties.ExtendedPANID})
		},
	}); err != nil {
		return err
	}

	if !version.IsV3() {
		z.logger.LogDebug(ctx, "Adapater Initialisation: Not Version 3.X.X.")
		/* Less than Z-Stack 3.X.X requires the Trust Centre key to be loaded. */
		return retryFunctions(ctx, []func(context.Context) error{
			func(invokeCtx context.Context) error {
				z.logger.LogDebug(ctx, "Adapater Initialisation: Enable default trust center.")
				return z.writeNVRAM(invokeCtx, ZCDNVUseDefaultTCLK{Enabled: 1})
			},
			func(invokeCtx context.Context) error {
				z.logger.LogDebug(ctx, "Adapater Initialisation: Configuring ZLL trust center key.")
				return z.writeNVRAM(invokeCtx, ZCDNVTCLKTableStart{
					Address:        zigbee.IEEEAddress(0xffffffffffffffff),
					NetworkKey:     zigbee.TCLinkKey,
					TXFrameCounter: 0,
					RXFrameCounter: 0,
				})
			},
		})
	} else {
		/* Z-Stack 3.X.X requires configuration of Base Device Behaviour. */
		z.logger.LogDebug(ctx, "Adapater Initialisation: Version 3.X.X.")
		if err := retryFunctions(ctx, []func(context.Context) error{
			func(invokeCtx context.Context) error {
				z.logger.LogDebug(ctx, "Adapater Initialisation: Configure primary channel.")
				return z.requestResponder.RequestResponse(ctx, APPCNFBDBSetChannelRequest{IsPrimary: true, Channel: channelBits}, &APPCNFBDBSetChannelRequestReply{})
			},
			func(invokeCtx context.Context) error {
				z.logger.LogDebug(ctx, "Adapater Initialisation: Configure secondary channels.")
				return z.requestResponder.RequestResponse(ctx, APPCNFBDBSetChannelRequest{IsPrimary: false, Channel: [4]byte{}}, &APPCNFBDBSetChannelRequestReply{})
			},
			func(invokeCtx context.Context) error {
				z.logger.LogDebug(ctx, "Adapater Initialisation: Request commissioning.")
				return z.requestResponder.RequestResponse(ctx, APPCNFBDBStartCommissioningRequest{Mode: 0x04}, &APPCNFBDBStartCommissioningRequestReply{})
			},
		}); err != nil {
			return err
		}

		z.logger.LogDebug(ctx, "Adapater Initialisation: Waiting for coordinator to start.")
		if err := z.waitForCoordinatorStart(ctx); err != nil {
			return err
		}

		return retryFunctions(ctx, []func(context.Context) error{
			func(invokeCtx context.Context) error {
				z.logger.LogDebug(ctx, "Adapater Initialisation: Waiting for commissioning to complete.")
				return z.requestResponder.RequestResponse(ctx, APPCNFBDBStartCommissioningRequest{Mode: 0x02}, &APPCNFBDBStartCommissioningRequestReply{})
			},
		})
	}
}

func (z *ZStack) retrieveAdapterAddresses(ctx context.Context) error {
	return retryFunctions(ctx, []func(context.Context) error{
		func(invokeCtx context.Context) error {
			if address, err := z.GetAdapterIEEEAddress(ctx); err != nil {
				return err
			} else {
				z.NetworkProperties.IEEEAddress = address
				return nil
			}
		},
		func(ctx context.Context) error {
			if address, err := z.GetAdapterNetworkAddress(ctx); err != nil {
				return err
			} else {
				z.NetworkProperties.NetworkAddress = address
				return nil
			}
		},
	})
}

func (z *ZStack) startZigbeeStack(ctx context.Context, version Version) error {
	if err := retry.Retry(ctx, DefaultZStackTimeout, DefaultZStackRetries, func(invokeCtx context.Context) error {
		return z.requestResponder.RequestResponse(invokeCtx, ZDOStartUpFromAppRequest{StartDelay: 100}, &ZDOStartUpFromAppRequestReply{})
	}); err != nil {
		return err
	}

	if version.IsV3() {
		return nil
	}

	ch := make(chan bool, 1)
	defer close(ch)

	err, cancel := z.subscriber.Subscribe(&ZDOStateChangeInd{}, func(v interface{}) {
		stateChange := v.(*ZDOStateChangeInd)
		z.logger.LogDebug(ctx, "Waiting for zigbee stack to start, state change.", logwrap.Datum("State", stateChange.State), logwrap.Datum("DesiredState", DeviceZBCoordinator))
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

func (z *ZStack) waitForCoordinatorStart(ctx context.Context) error {
	ch := make(chan bool, 1)
	defer close(ch)

	err, cancel := z.subscriber.Subscribe(&ZDOStateChangeInd{}, func(v interface{}) {
		stateChange := v.(*ZDOStateChangeInd)
		z.logger.LogDebug(ctx, "Waiting for coordinator start, state change.", logwrap.Datum("State", stateChange.State), logwrap.Datum("DesiredState", DeviceZBCoordinator))
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

func retryFunctions(ctx context.Context, funcs []func(context.Context) error) error {
	for _, f := range funcs {
		if err := retry.Retry(ctx, DefaultZStackTimeout, DefaultZStackRetries, f); err != nil {
			return fmt.Errorf("failed during configuration and initialisation: %w", err)
		}
	}

	return nil
}

func channelToBits(channel uint8) [4]byte {
	channelBits := 1 << channel

	channelBytes := [4]byte{}
	channelBytes[0] = byte((channelBits >> 0) & 0xff)
	channelBytes[1] = byte((channelBits >> 8) & 0xff)
	channelBytes[2] = byte((channelBits >> 16) & 0xff)
	channelBytes[3] = byte((channelBits >> 24) & 0xff)

	return channelBytes
}

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

type APPCNFBDBStartCommissioningRequest struct {
	Mode uint8
}

const APPCNFBDBStartCommissioningRequestID uint8 = 0x05

type APPCNFBDBStartCommissioningRequestReply GenericZStackStatus

const APPCNFBDBStartCommissioningRequestReplyID uint8 = 0x05

type APPCNFBDBSetChannelRequest struct {
	IsPrimary bool `bcwidth:"8"`
	Channel   [4]byte
}

const APPCNFBDBSetChannelRequestID uint8 = 0x08

type APPCNFBDBSetChannelRequestReply GenericZStackStatus

const APPCNFBDBSetChannelRequestReplyID uint8 = 0x08

type ZDOStartUpFromAppRequest struct {
	StartDelay uint16
}

const ZDOStartUpFromAppRequestId uint8 = 0x40

type ZDOStartUpFromAppRequestReply struct {
	Status uint8
}

const ZDOStartUpFromAppRequestReplyID uint8 = 0x40
