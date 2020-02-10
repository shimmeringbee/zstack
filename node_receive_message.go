package zstack

import (
	"github.com/shimmeringbee/zigbee"
	"log"
)

func (z *ZStack) startMessageReceiver() {
	_, z.messageReceiverStop = z.subscriber.Subscribe(&AfIncomingMsg{}, func(v interface{}) {
		msg := v.(*AfIncomingMsg)

		device, found := z.deviceTable.GetByNetwork(msg.SourceAddress)

		if !found {
			log.Printf("could not resolve IEEE address while receiving message: network address = %d", msg.SourceAddress)
			return
		}

		z.sendEvent(zigbee.DeviceIncomingMessageEvent{
			GroupID:              msg.GroupID,
			ClusterID:            msg.ClusterID,
			SourceIEEEAddress:    device.IEEEAddress,
			SourceNetworkAddress: msg.SourceAddress,
			SourceEndpoint:       msg.SourceEndpoint,
			DestinationEndpoint:  msg.DestinationEndpoint,
			Broadcast:            msg.WasBroadcast != 0,
			Secure:               msg.SecurityUse != 0,
			LinkQuality:          msg.LinkQuality,
			Sequence:             msg.Sequence,
			Data:                 msg.Data,
		})
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
