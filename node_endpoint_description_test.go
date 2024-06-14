package zstack

import (
	"context"
	"github.com/shimmeringbee/bytecodec"
	"github.com/shimmeringbee/persistence/impl/memory"
	. "github.com/shimmeringbee/unpi"
	unpiTest "github.com/shimmeringbee/unpi/testing"
	"github.com/shimmeringbee/zigbee"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/semaphore"
	"testing"
	"time"
)

func Test_QueryNodeEndpointDescription(t *testing.T) {
	t.Run("returns an success on query, response for requested network address is received", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock, memory.New())
		zstack.sem = semaphore.NewWeighted(8)
		defer unpiMock.Stop()

		unpiMock.On(SREQ, ZDO, ZdoSimpleDescReqID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   ZDO,
			CommandID:   ZdoSimpleDescReqReplyID,
			Payload:     []byte{0x00},
		})

		go func() {
			time.Sleep(10 * time.Millisecond)
			unpiMock.InjectOutgoing(Frame{
				MessageType: AREQ,
				Subsystem:   ZDO,
				CommandID:   ZdoSimpleDescRspID,
				Payload:     []byte{0x00, 0x20, 0x00, 0x00, 0x40, 0xff, 0x01, 0x01, 0x01, 0x01, 0x00, 0x02, 0x01, 0x01, 0x00, 0x01, 0x02, 0x00},
			})
		}()

		zstack.nodeTable.addOrUpdate(zigbee.IEEEAddress(0x11223344556677), zigbee.NetworkAddress(0x4000))

		endpoints, err := zstack.QueryNodeEndpointDescription(ctx, zigbee.IEEEAddress(0x11223344556677), 0x01)
		assert.NoError(t, err)
		assert.Equal(t, zigbee.EndpointDescription{
			Endpoint:       0x01,
			ProfileID:      0x0101,
			DeviceID:       1,
			DeviceVersion:  2,
			InClusterList:  []zigbee.ClusterID{0x0001},
			OutClusterList: []zigbee.ClusterID{0x0002},
		}, endpoints)

		unpiMock.AssertCalls(t)
	})
}

func Test_EndpointDescriptionMessages(t *testing.T) {
	t.Run("verify ZdoSimpleDescReq marshals", func(t *testing.T) {
		req := ZdoSimpleDescReq{
			DestinationAddress: zigbee.NetworkAddress(0x2000),
			OfInterestAddress:  zigbee.NetworkAddress(0x4000),
			Endpoint:           0x08,
		}

		data, err := bytecodec.Marshal(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x00, 0x20, 0x00, 0x40, 0x08}, data)
	})

	t.Run("verify ZdoSimpleDescReqReply marshals", func(t *testing.T) {
		req := ZdoSimpleDescReqReply{
			Status: 1,
		}

		data, err := bytecodec.Marshal(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x01}, data)
	})

	t.Run("ZdoSimpleDescReqReply returns true if success", func(t *testing.T) {
		g := ZdoSimpleDescReqReply{Status: ZSuccess}
		assert.True(t, g.WasSuccessful())
	})

	t.Run("ZdoSimpleDescReqReply returns false if not success", func(t *testing.T) {
		g := ZdoSimpleDescReqReply{Status: ZFailure}
		assert.False(t, g.WasSuccessful())
	})

	t.Run("verify ZdoSimpleDescRsp marshals", func(t *testing.T) {
		req := ZdoSimpleDescRsp{
			SourceAddress:     zigbee.NetworkAddress(0x2000),
			Status:            1,
			OfInterestAddress: zigbee.NetworkAddress(0x4000),
			Length:            0x0a,
			Endpoint:          0x08,
			ProfileID:         0x1234,
			DeviceID:          0x5678,
			DeviceVersion:     0,
			InClusterList:     []zigbee.ClusterID{0x1234},
			OutClusterList:    []zigbee.ClusterID{0x5678},
		}

		data, err := bytecodec.Marshal(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x00, 0x20, 0x01, 0x00, 0x40, 0x0a, 0x08, 0x34, 0x12, 0x78, 0x56, 0x00, 0x01, 0x34, 0x12, 0x01, 0x78, 0x56}, data)
	})

	t.Run("ZdoSimpleDescRsp returns true if success", func(t *testing.T) {
		g := ZdoSimpleDescRsp{Status: ZSuccess}
		assert.True(t, g.WasSuccessful())
	})

	t.Run("ZdoSimpleDescRsp returns false if not success", func(t *testing.T) {
		g := ZdoSimpleDescRsp{Status: ZFailure}
		assert.False(t, g.WasSuccessful())
	})
}
