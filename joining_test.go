package zstack

import (
	"context"
	"github.com/shimmeringbee/bytecodec"
	. "github.com/shimmeringbee/unpi"
	unpiTest "github.com/shimmeringbee/unpi/testing"
	"github.com/shimmeringbee/zigbee"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/semaphore"
	"testing"
	"time"
)

func Test_PermitJoin(t *testing.T) {
	t.Run("permit join for all routers sends message to all routers permitting joining", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock, NewNodeTable())
		zstack.sem = semaphore.NewWeighted(8)
		defer unpiMock.Stop()

		c := unpiMock.On(SREQ, ZDO, ZDOMgmtPermitJoinRequestID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   ZDO,
			CommandID:   ZDOMgmtPermitJoinRequestReplyID,
			Payload:     []byte{0x00},
		})

		err := zstack.PermitJoin(ctx, true)
		assert.NoError(t, err)

		unpiMock.AssertCalls(t)

		assert.Equal(t, []byte{0xfc, 0xff, 0xff, 0x00}, c.CapturedCalls[0].Frame.Payload)
		assert.Equal(t, OnAllRouters, zstack.NetworkProperties.JoinState)
	})

	t.Run("permit join for the coordinator sends message to coordinator permitting joining", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock, NewNodeTable())
		zstack.sem = semaphore.NewWeighted(8)
		defer unpiMock.Stop()
		zstack.NetworkProperties.NetworkAddress = zigbee.NetworkAddress(0x0102)

		c := unpiMock.On(SREQ, ZDO, ZDOMgmtPermitJoinRequestID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   ZDO,
			CommandID:   ZDOMgmtPermitJoinRequestReplyID,
			Payload:     []byte{0x00},
		})

		err := zstack.PermitJoin(ctx, false)
		assert.NoError(t, err)

		unpiMock.AssertCalls(t)

		assert.Equal(t, []byte{0x02, 0x01, 0xff, 0x00}, c.CapturedCalls[0].Frame.Payload)
		assert.Equal(t, OnCoordinator, zstack.NetworkProperties.JoinState)
	})

	t.Run("permit join rejection by adapter errors", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock, NewNodeTable())
		zstack.sem = semaphore.NewWeighted(8)
		defer unpiMock.Stop()

		unpiMock.On(SREQ, ZDO, ZDOMgmtPermitJoinRequestID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   ZDO,
			CommandID:   ZDOMgmtPermitJoinRequestReplyID,
			Payload:     []byte{0x01},
		})

		err := zstack.PermitJoin(ctx, true)
		assert.Error(t, err)

		unpiMock.AssertCalls(t)
	})
}

func Test_DenyJoin(t *testing.T) {
	t.Run("denying join sends message to all routers disabling joining", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock, NewNodeTable())
		zstack.sem = semaphore.NewWeighted(8)
		defer unpiMock.Stop()

		c := unpiMock.On(SREQ, ZDO, ZDOMgmtPermitJoinRequestID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   ZDO,
			CommandID:   ZDOMgmtPermitJoinRequestReplyID,
			Payload:     []byte{0x00},
		})

		zstack.NetworkProperties.JoinState = OnCoordinator
		err := zstack.DenyJoin(ctx)
		assert.NoError(t, err)

		unpiMock.AssertCalls(t)

		assert.Equal(t, []byte{0xfc, 0xff, 0x00, 0x00}, c.CapturedCalls[0].Frame.Payload)
		assert.Equal(t, Off, zstack.NetworkProperties.JoinState)
	})

	t.Run("denying join rejection by adapter errors", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		unpiMock := unpiTest.NewMockAdapter()
		zstack := New(unpiMock, NewNodeTable())
		zstack.sem = semaphore.NewWeighted(8)
		defer unpiMock.Stop()

		unpiMock.On(SREQ, ZDO, ZDOMgmtPermitJoinRequestID).Return(Frame{
			MessageType: SRSP,
			Subsystem:   ZDO,
			CommandID:   ZDOMgmtPermitJoinRequestReplyID,
			Payload:     []byte{0x01},
		})

		err := zstack.DenyJoin(ctx)
		assert.Error(t, err)

		unpiMock.AssertCalls(t)
	})
}

func Test_ZDOMgmtPermitJoin(t *testing.T) {
	t.Run("ZDOMgmtPermitJoinRequest", func(t *testing.T) {
		s := ZDOMgmtPermitJoinRequest{
			Destination:    0x0102,
			Duration:       0x03,
			TCSignificance: 0x04,
		}

		actualBytes, err := bytecodec.Marshal(s)

		expectedBytes := []byte{0x02, 0x01, 0x03, 0x04}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("ZDOMgmtPermitJoinRequestReply", func(t *testing.T) {
		s := ZDOMgmtPermitJoinRequestReply{
			Status: 0x01,
		}

		actualBytes, err := bytecodec.Marshal(s)

		expectedBytes := []byte{0x01}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})
}
