package zstack

import (
	"context"
	"github.com/shimmeringbee/zigbee"
)

func (z *ZStack) SendNodeMessage(ctx context.Context, destinationAddress zigbee.IEEEAddress, sourceEndpoint byte, destinationEndpoint byte, cluster zigbee.ZCLClusterID, data []byte) error {
	return nil
}

type AfDataRequest struct {
	DestinationAddress  zigbee.NetworkAddress
	DestinationEndpoint byte
	SourceEndpoint      byte
	ClusterID           zigbee.ZCLClusterID
	TransactionID       uint8
	Options             uint8
	Radius              uint8
	Data                []byte `bclength:"8"`
}

const AfDataRequestID uint8 = 0x01

type AfDataRequestReply GenericZStackStatus

func (s AfDataRequestReply) WasSuccessful() bool {
	return s.Status == ZSuccess
}

const AfDataRequestReplyID uint8 = 0x01

type AfDataConfirm struct {
	Status        ZStackStatus
	Endpoint      byte
	TransactionID uint8
}

func (s AfDataConfirm) WasSuccessful() bool {
	return s.Status == ZSuccess
}

const AfDataConfirmID uint8 = 0x80