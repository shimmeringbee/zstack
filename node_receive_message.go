package zstack

import (
	"context"
	"github.com/shimmeringbee/logwrap"
	"github.com/shimmeringbee/zigbee"
)

func (z *ZStack) startMessageReceiver() {
	_, z.messageReceiverStop = z.subscriber.Subscribe(&AfIncomingMsg{}, func(v interface{}) {
		msg := v.(*AfIncomingMsg)

		ctx, cancel := context.WithTimeout(context.Background(), DefaultResolveIEEETimeout)
		defer cancel()

		ieee, err := z.ResolveNodeIEEEAddress(ctx, msg.SourceAddress)
		if err != nil {
			z.logger.LogError(ctx, "Received AfIncomingMsg (application message), however unable to resolve IEEE address to issue event.", logwrap.Err(err), logwrap.Datum("NetworkAddress", msg.SourceAddress))
			return
		}

		node, _ := z.nodeTable.getByIEEE(ieee)

		z.sendEvent(zigbee.NodeIncomingMessageEvent{
			Node: node,
			IncomingMessage: zigbee.IncomingMessage{
				GroupID: msg.GroupID,
				SourceAddress: zigbee.SourceAddress{
					IEEEAddress:    ieee,
					NetworkAddress: msg.SourceAddress,
				},
				Broadcast:   msg.WasBroadcast,
				Secure:      msg.SecurityUse,
				LinkQuality: msg.LinkQuality,
				Sequence:    msg.Sequence,
				ApplicationMessage: zigbee.ApplicationMessage{
					ClusterID:           msg.ClusterID,
					SourceEndpoint:      msg.SourceEndpoint,
					DestinationEndpoint: msg.DestinationEndpoint,
					Data:                msg.Data,
				},
			},
		})

		z.nodeTable.update(ieee, updateReceived, lqi(msg.LinkQuality))
	})
}

func (z *ZStack) stopMessageReceiver() {
	if z.messageReceiverStop != nil {
		z.messageReceiverStop()
	}
}

type AfIncomingMsg struct {
	GroupID             zigbee.GroupID
	ClusterID           zigbee.ClusterID
	SourceAddress       zigbee.NetworkAddress
	SourceEndpoint      zigbee.Endpoint
	DestinationEndpoint zigbee.Endpoint
	WasBroadcast        bool
	LinkQuality         uint8
	SecurityUse         bool
	TimeStamp           uint32
	Sequence            uint8
	Data                []byte `bcsliceprefix:"8"`
}

const AfIncomingMsgID uint8 = 0x81
