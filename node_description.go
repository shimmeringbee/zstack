package zstack

import (
	"context"
	"fmt"
	"github.com/shimmeringbee/zigbee"
)

func (z *ZStack) QueryNodeDescription(ctx context.Context, ieeeAddress zigbee.IEEEAddress) (zigbee.NodeDescription, error) {
	nwkAddress, err := z.ResolveNodeNWKAddress(ctx, ieeeAddress)
	if err != nil {
		return zigbee.NodeDescription{}, err
	}

	if err := z.sem.Acquire(ctx, 1); err != nil {
		return zigbee.NodeDescription{}, fmt.Errorf("failed to acquire semaphore: %w", err)
	}
	defer z.sem.Release(1)

	request := ZdoNodeDescReq{
		DestinationAddress: nwkAddress,
		OfInterestAddress:  nwkAddress,
	}

	resp, err := z.nodeRequest(ctx, &request, &ZdoNodeDescReqReply{}, &ZdoNodeDescRsp{}, func(i interface{}) bool {
		msg := i.(*ZdoNodeDescRsp)
		return msg.OfInterestAddress == nwkAddress
	})

	castResp, ok := resp.(*ZdoNodeDescRsp)

	if ok {
		return zigbee.NodeDescription{
			LogicalType:      castResp.Capabilities.LogicalType,
			ManufacturerCode: zigbee.ManufacturerCode(castResp.ManufacturerCode),
		}, nil
	} else {
		return zigbee.NodeDescription{}, err
	}
}

type ZdoNodeDescReq struct {
	DestinationAddress zigbee.NetworkAddress
	OfInterestAddress  zigbee.NetworkAddress
}

const ZdoNodeDescReqID uint8 = 0x02

type ZdoNodeDescReqReply GenericZStackStatus

func (r ZdoNodeDescReqReply) WasSuccessful() bool {
	return r.Status == ZSuccess
}

const ZdoNodeDescReqReplyID uint8 = 0x02

type ZdoNodeDescRspCapabilities struct {
	Reserved                   uint8              `bcfieldwidth:"3"`
	UserDescriptorAvailable    bool               `bcfieldwidth:"1"`
	ComplexDescriptorAvailable bool               `bcfieldwidth:"1"`
	LogicalType                zigbee.LogicalType `bcfieldwidth:"3"`
}

type ZdoNodeDescRspServerMask struct {
	Reserved0                uint8 `bcfieldwidth:"8"`
	Reserved1                uint8 `bcfieldwidth:"2"`
	BackupDiscoveryCache     bool  `bcfieldwidth:"1"`
	PrimaryDiscoveryCache    bool  `bcfieldwidth:"1"`
	BackupBindingTableCache  bool  `bcfieldwidth:"1"`
	PrimaryBindingTableCache bool  `bcfieldwidth:"1"`
	BackupTrustCenter        bool  `bcfieldwidth:"1"`
	PrimaryTrustCenter       bool  `bcfieldwidth:"1"`
}

type ZdoNodeDescRsp struct {
	SourceAddress          zigbee.NetworkAddress
	Status                 ZStackStatus
	OfInterestAddress      zigbee.NetworkAddress
	Capabilities           ZdoNodeDescRspCapabilities
	NodeFrequencyBand      uint8 `bcfieldwidth:"3"`
	APSFlags               uint8 `bcfieldwidth:"5"`
	MacCapabilitiesFlags   uint8
	ManufacturerCode       uint16
	MaxBufferSize          uint8
	MaxInTransferSize      uint16
	ServerMask             ZdoNodeDescRspServerMask
	MaxOutTransferSize     uint16
	DescriptorCapabilities uint8
}

func (r ZdoNodeDescRsp) WasSuccessful() bool {
	return r.Status == ZSuccess
}

const ZdoNodeDescRspID uint8 = 0x82
