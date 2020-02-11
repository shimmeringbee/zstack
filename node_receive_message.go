package zstack

import (
	"context"
	"github.com/shimmeringbee/zigbee"
	"log"
)

func (z *ZStack) startMessageReceiver() {
	_, z.messageReceiverStop = z.subscriber.Subscribe(&AfIncomingMsg{}, func(v interface{}) {
		msg := v.(*AfIncomingMsg)

		ctx, cancel := context.WithTimeout(context.Background(), DefaultResolveIEEETimeout)
		defer cancel()

		ieee, err := z.ResolveNodeIEEEAddress(ctx, msg.SourceAddress)

		if err != nil {
			log.Printf("could not resolve IEEE address while receiving message: network address = %d, err = %+v", msg.SourceAddress, err)
			return
		}

		device, _ := z.deviceTable.GetByIEEE(ieee)

		z.sendEvent(zigbee.DeviceIncomingMessageEvent{
			Device: device,
			IncomingMessage: zigbee.IncomingMessage{
				GroupID:              msg.GroupID,
				ClusterID:            msg.ClusterID,
				SourceIEEEAddress:    ieee,
				SourceNetworkAddress: msg.SourceAddress,
				SourceEndpoint:       msg.SourceEndpoint,
				DestinationEndpoint:  msg.DestinationEndpoint,
				Broadcast:            msg.WasBroadcast != 0,
				Secure:               msg.SecurityUse != 0,
				LinkQuality:          msg.LinkQuality,
				Sequence:             msg.Sequence,
				Data:                 msg.Data,
			},
		})

		z.deviceTable.Update(ieee, UpdateReceived)
	})
}

func (z *ZStack) stopMessageReceiver() {
	if z.messageReceiverStop != nil {
		z.messageReceiverStop()
	}
}

type AfIncomingMsg struct {
	GroupID             uint16
	ClusterID           zigbee.ZCLClusterID
	SourceAddress       zigbee.NetworkAddress
	SourceEndpoint      byte
	DestinationEndpoint byte
	WasBroadcast        uint8
	LinkQuality         uint8
	SecurityUse         uint8
	TimeStamp           uint32
	Sequence            uint8
	Data                []byte `bclength:"8"`
}

const AfIncomingMsgID uint8 = 0x81
