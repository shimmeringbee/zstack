package zstack

import (
	"context"
	"errors"
	"github.com/shimmeringbee/zigbee"
)

const DefaultRadius uint8 = 0x20

func (z *ZStack) SendApplicationMessageToNode(ctx context.Context, destinationAddress zigbee.IEEEAddress, message zigbee.ApplicationMessage, requireAck bool) error {
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
		DestinationEndpoint: message.DestinationEndpoint,
		SourceEndpoint:      message.SourceEndpoint,
		ClusterID:           message.ClusterID,
		TransactionID:       transactionId,
		Options:             AfDataRequestOptions{ACKRequest: true},
		Radius:              DefaultRadius,
		Data:                message.Data,
	}

	_, err = z.nodeRequest(ctx, &request, &AfDataRequestReply{}, &AfDataConfirm{}, func(i interface{}) bool {
		msg := i.(*AfDataConfirm)
		return msg.TransactionID == transactionId && msg.Endpoint == message.DestinationEndpoint
	})

	return err
}

type AfDataRequestOptions struct {
	Reserved0      uint8 `bcfieldwidth:"1"`
	EnableSecurity bool  `bcfieldwidth:"1"`
	DiscoveryRoute bool  `bcfieldwidth:"1"`
	ACKRequest     bool  `bcfieldwidth:"1"`
	Reserved1      uint8 `bcfieldwidth:"4"`
}

type AfDataRequest struct {
	DestinationAddress  zigbee.NetworkAddress
	DestinationEndpoint zigbee.Endpoint
	SourceEndpoint      zigbee.Endpoint
	ClusterID           zigbee.ClusterID
	TransactionID       uint8
	Options             AfDataRequestOptions
	Radius              uint8
	Data                []byte `bcsliceprefix:"8"`
}

const AfDataRequestID uint8 = 0x01

type AfDataRequestReply GenericZStackStatus

func (s AfDataRequestReply) WasSuccessful() bool {
	return s.Status == ZSuccess
}

const AfDataRequestReplyID uint8 = 0x01

type AfDataConfirm struct {
	Status        ZStackStatus
	Endpoint      zigbee.Endpoint
	TransactionID uint8
}

func (s AfDataConfirm) WasSuccessful() bool {
	return s.Status == ZSuccess
}

const AfDataConfirmID uint8 = 0x80
