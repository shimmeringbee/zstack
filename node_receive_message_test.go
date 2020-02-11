package zstack

import (
	"context"
	"github.com/shimmeringbee/bytecodec"
	. "github.com/shimmeringbee/unpi"
	unpiTest "github.com/shimmeringbee/unpi/testing"
	"github.com/shimmeringbee/zigbee"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_ReceiveMessage(t *testing.T) {
	t.Run("messages which are received with a known network to ieee mapping are sent to event stream", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		defer unpiMock.AssertCalls(t)
		zstack := New(unpiMock)
		defer unpiMock.Stop()

		zstack.deviceTable.AddOrUpdate(zigbee.IEEEAddress(0x1122334455667788), zigbee.NetworkAddress(0x1000))

		zstack.startMessageReceiver()

		go func() {
			time.Sleep(10 * time.Millisecond)

			msg := AfIncomingMsg{
				GroupID:             0x01,
				ClusterID:           0x02,
				SourceAddress:       0x1000,
				SourceEndpoint:      3,
				DestinationEndpoint: 4,
				WasBroadcast:        1,
				LinkQuality:         55,
				SecurityUse:         1,
				TimeStamp:           123412,
				Sequence:            63,
				Data:                []byte{0x01, 0x02},
			}

			data, _ := bytecodec.Marshall(&msg)

			unpiMock.InjectOutgoing(Frame{
				MessageType: AREQ,
				Subsystem:   AF,
				CommandID:   AfIncomingMsgID,
				Payload:     data,
			})
		}()

		event, err := zstack.ReadEvent(ctx)
		assert.NoError(t, err)

		incommingMsg, ok := event.(zigbee.DeviceIncomingMessageEvent)
		assert.True(t, ok)

		expectedMsg := zigbee.DeviceIncomingMessageEvent{
			Device: zigbee.Device{
				IEEEAddress:    0x1122334455667788,
				NetworkAddress: 0x1000,
				LogicalType:    0xff,
				LQI:            0,
				Depth:          0,
				LastDiscovered: time.Time{},
				LastReceived:   time.Time{},
			},
			IncomingMessage: zigbee.IncomingMessage{
				GroupID:              0x01,
				ClusterID:            0x02,
				SourceNetworkAddress: 0x1000,
				SourceIEEEAddress:    0x1122334455667788,
				SourceEndpoint:       3,
				DestinationEndpoint:  4,
				Broadcast:            true,
				Secure:               true,
				LinkQuality:          55,
				Sequence:             63,
				Data:                 []byte{0x01, 0x02},
			},
		}

		assert.Equal(t, expectedMsg, incommingMsg)
	})
}

func Test_IncomingMessage(t *testing.T) {
	t.Run("verify AfIncomingMsg marshals", func(t *testing.T) {
		req := AfIncomingMsg{
			GroupID:             0x0102,
			ClusterID:           0x0304,
			SourceAddress:       0x0506,
			SourceEndpoint:      0x07,
			DestinationEndpoint: 0x08,
			WasBroadcast:        0x09,
			LinkQuality:         0x0a,
			SecurityUse:         0x0b,
			TimeStamp:           0x10111213,
			Sequence:            0x0d,
			Data:                []byte{0x00, 0x01, 0x02},
		}

		data, err := bytecodec.Marshall(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x02, 0x01, 0x04, 0x03, 0x06, 0x05, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x13, 0x12, 0x11, 0x10, 0x0d, 0x03, 0x00, 0x01, 0x02}, data)
	})
}
