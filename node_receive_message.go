package zstack

import "github.com/shimmeringbee/zigbee"

func (z *ZStack) startMessageReceiver() {
	_, z.messageReceiverStop = z.subscriber.Subscribe(&AfIncomingMsg{}, func(v interface{}) {
		msg := v.(*AfIncomingMsg)

		sourceAddress := z.devicesByNetAddr[msg.SourceAddress]

		z.sendEvent(zigbee.DeviceIncomingMessageEvent{
			GroupID:             msg.GroupID,
			ClusterID:           msg.ClusterID,
			SourceAddress:       sourceAddress,
			SourceEndpoint:      msg.SourceEndpoint,
			DestinationEndpoint: msg.DestinationEndpoint,
			Broadcast:           msg.WasBroadcast != 0,
			Secure:              msg.SecurityUse != 0,
			LinkQuality:         msg.LinkQuality,
			Sequence:            msg.Sequence,
			Data:                msg.Data,
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