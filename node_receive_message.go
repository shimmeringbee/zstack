package zstack

import "github.com/shimmeringbee/zigbee"

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