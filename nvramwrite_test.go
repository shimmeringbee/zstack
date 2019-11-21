package zstack

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestZStack_writeNVRAM(t *testing.T) {
	t.Run("verifies that a request response is made to unpi", func(t *testing.T) {
		mrr := new(MockRequestResponder)

		mrr.On("MessageRequestResponse", mock.Anything, SysOSALNVWriteReq{
			NVItemID: ZCDNVLogicalTypeID,
			Offset:   0,
			Value:    []byte{0x02},
		}, &SysOSALNVWriteResp{}).Return(nil)

		z := ZStack{RequestResponder: mrr}

		ctx, cancel := context.WithTimeout(context.Background(), 50 * time.Millisecond)
		defer cancel()

		err := z.writeNVRAM(ctx, ZCDNVLogicalType{LogicalType:EndDevice})

		mrr.AssertExpectations(t)
		assert.NoError(t, err)
	})

	t.Run("verifies that write requests that fail raise an error", func(t *testing.T) {
		mrr := &FailingMockRequestResponse{}

		z := ZStack{RequestResponder: mrr}

		ctx, cancel := context.WithTimeout(context.Background(), 50 * time.Millisecond)
		defer cancel()

		err := z.writeNVRAM(ctx, ZCDNVLogicalType{LogicalType:EndDevice})

		assert.Error(t, err)
	})

	t.Run("verifies that a request response with errors is raised", func(t *testing.T) {
		mrr := new(MockRequestResponder)

		mrr.On("MessageRequestResponse", mock.Anything, SysOSALNVWriteReq{
			NVItemID: ZCDNVLogicalTypeID,
			Offset:   0,
			Value:    []byte{0x02},
		}, &SysOSALNVWriteResp{}).Return(errors.New("context expired"))

		z := ZStack{RequestResponder: mrr}

		ctx, cancel := context.WithTimeout(context.Background(), 50 * time.Millisecond)
		defer cancel()

		err := z.writeNVRAM(ctx, ZCDNVLogicalType{LogicalType:EndDevice})

		mrr.AssertExpectations(t)
		assert.Error(t, err)
	})

	t.Run("verifies that unknown structure raises an error", func(t *testing.T) {
		mrr := new(MockRequestResponder)

		z := ZStack{RequestResponder: mrr}

		ctx, cancel := context.WithTimeout(context.Background(), 50 * time.Millisecond)
		defer cancel()

		err := z.writeNVRAM(ctx, struct{}{})

		mrr.AssertExpectations(t)
		assert.Error(t, err)
	})
}

type FailingMockRequestResponse struct {}

func (m *FailingMockRequestResponse) MessageRequestResponse(ctx context.Context, req interface{}, resp interface{}) error {
	response, ok := resp.(*SysOSALNVWriteResp)

	if !ok {
		panic("incorrect type passed to mock")
	}

	response.Status = 0x01

	return nil
}