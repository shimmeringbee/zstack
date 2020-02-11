package zstack

import (
	"context"
	"errors"
	"github.com/shimmeringbee/zigbee"
)

const DefaultRadius uint8 = 0x20

func (z *ZStack) SendNodeMessage(ctx context.Context, destinationAddress zigbee.IEEEAddress, sourceEndpoint byte, destinationEndpoint byte, cluster zigbee.ZCLClusterID, data []byte) error {
	network, err := z.ResolveNodeNWKAddress(ctx, destinationAddress)

	if err != nil {
		return err
	}

	var transactionId uint8

	select {
	case transactionId = <-z.transactionIdStore:
		defer func() { z.transactionIdStore <- transactionId }()
	case <-ctx.Done():
		return errors.New("context expired while obtaining a free transaction ID")
	}

	request := AfDataRequest{
		DestinationAddress:  network,
		DestinationEndpoint: destinationEndpoint,
		SourceEndpoint:      sourceEndpoint,
		ClusterID:           cluster,
		TransactionID:       transactionId,
		Options:             0x10, //AF_ACK_REQUEST
		Radius:              DefaultRadius,
		Data:                data,
	}

	_, err = z.nodeRequest(ctx, &request, &AfDataRequestReply{}, &AfDataConfirm{}, func(i interface{}) bool {
		msg := i.(*AfDataConfirm)
		return msg.TransactionID == transactionId && msg.Endpoint == destinationEndpoint
	})

	return err
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
