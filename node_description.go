package zstack

import (
	"context"
	"github.com/shimmeringbee/zigbee"
)

func (z *ZStack) QueryNodeDescription(ctx context.Context, ieeeAddress zigbee.IEEEAddress) (zigbee.NodeDescription, error) {
	nwkAddress, err := z.ResolveNodeNWKAddress(ctx, ieeeAddress)

	if err != nil {
		return zigbee.NodeDescription{}, err
	}

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
			LogicalType:      zigbee.LogicalType(castResp.LogicalTypeDescriptor >> 5),
			ManufacturerCode: castResp.ManufacturerCode,
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

type ZdoNodeDescRsp struct {
	SourceAddress          zigbee.NetworkAddress
	Status                 ZStackStatus
	OfInterestAddress      zigbee.NetworkAddress
	LogicalTypeDescriptor  uint8
	APSFlagsFrequency      uint8
	MacCapabilitiesFlags   uint8
	ManufacturerCode       uint16
	MaxBufferSize          uint8
	MaxInTransferSize      uint16
	ServerMask             uint16
	MaxOutTransferSize     uint16
	DescriptorCapabilities uint8
}

func (r ZdoNodeDescRsp) WasSuccessful() bool {
	return r.Status == ZSuccess
}

const ZdoNodeDescRspID uint8 = 0x82
