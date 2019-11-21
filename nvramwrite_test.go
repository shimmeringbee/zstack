package zstack

import (
	"context"
	"errors"
	"github.com/shimmeringbee/bytecodec"
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

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		err := z.writeNVRAM(ctx, ZCDNVLogicalType{LogicalType: EndDevice})

		mrr.AssertExpectations(t)
		assert.NoError(t, err)
	})

	t.Run("verifies that write requests that fail raise an error", func(t *testing.T) {
		mrr := &FailingMockRequestResponse{}

		z := ZStack{RequestResponder: mrr}

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		err := z.writeNVRAM(ctx, ZCDNVLogicalType{LogicalType: EndDevice})

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

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		err := z.writeNVRAM(ctx, ZCDNVLogicalType{LogicalType: EndDevice})

		mrr.AssertExpectations(t)
		assert.Error(t, err)
	})

	t.Run("verifies that unknown structure raises an error", func(t *testing.T) {
		mrr := new(MockRequestResponder)

		z := ZStack{RequestResponder: mrr}

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		err := z.writeNVRAM(ctx, struct{}{})

		mrr.AssertExpectations(t)
		assert.Error(t, err)
	})
}

type FailingMockRequestResponse struct{}

func (m *FailingMockRequestResponse) MessageRequestResponse(ctx context.Context, req interface{}, resp interface{}) error {
	response, ok := resp.(*SysOSALNVWriteResp)

	if !ok {
		panic("incorrect type passed to mock")
	}

	response.Status = 0x01

	return nil
}

func TestZStack_NVRAMStructs(t *testing.T) {
	t.Run("NCDNVStartUpOption", func(t *testing.T) {
		s := NCDNVStartUpOption{
			StartOption: 3,
		}

		actualBytes, err := bytecodec.Marshall(s)

		expectedBytes := []byte{0x03}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("ZCDNVLogicalType", func(t *testing.T) {
		s := ZCDNVLogicalType{
			LogicalType: Router,
		}

		actualBytes, err := bytecodec.Marshall(s)

		expectedBytes := []byte{0x01}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("ZCDNVSecurityMode", func(t *testing.T) {
		s := ZCDNVSecurityMode{
			Enabled: 1,
		}

		actualBytes, err := bytecodec.Marshall(s)

		expectedBytes := []byte{0x01}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("ZCDNVPreCfgKeysEnable", func(t *testing.T) {
		s := ZCDNVPreCfgKeysEnable{
			Enabled: 1,
		}

		actualBytes, err := bytecodec.Marshall(s)

		expectedBytes := []byte{0x01}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("ZCDNVPreCfgKey", func(t *testing.T) {
		s := ZCDNVPreCfgKey{
			NetworkKey: [16]byte{0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03},
		}

		actualBytes, err := bytecodec.Marshall(s)

		expectedBytes := []byte{0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("ZCDNVZDODirectCB", func(t *testing.T) {
		s := ZCDNVZDODirectCB{
			Enabled: 1,
		}

		actualBytes, err := bytecodec.Marshall(s)

		expectedBytes := []byte{0x01}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("ZCDNVChanList", func(t *testing.T) {
		s := ZCDNVChanList{
			Channels: [4]byte{0x00, 0x00, 0x00, 0x03},
		}

		actualBytes, err := bytecodec.Marshall(s)

		expectedBytes := []byte{0x00, 0x00, 0x00, 0x03}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("ZCDNVPANID", func(t *testing.T) {
		s := ZCDNVPANID{
			PANID: [2]byte{0x01, 0x02},
		}

		actualBytes, err := bytecodec.Marshall(s)

		expectedBytes := []byte{0x01, 0x02}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("ZCDNVExtPANID", func(t *testing.T) {
		s := ZCDNVExtPANID{
			ExtendedPANID: [8]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
		}

		actualBytes, err := bytecodec.Marshall(s)

		expectedBytes := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("ZCDNVUseDefaultTCLK", func(t *testing.T) {
		s := ZCDNVUseDefaultTCLK{
			Enabled: 1,
		}

		actualBytes, err := bytecodec.Marshall(s)

		expectedBytes := []byte{0x01}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("ZCDNVTCLKTableStart", func(t *testing.T) {
		s := ZCDNVTCLKTableStart{
			Address:        [8]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			NetworkKey:     [16]byte{0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03},
			TXFrameCounter: 123456,
			RXFrameCounter: 654321,
		}

		actualBytes, err := bytecodec.Marshall(s)

		expectedBytes := []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x0, 0x1, 0x2, 0x3, 0x0, 0x1, 0x2, 0x3, 0x0, 0x1, 0x2, 0x3, 0x0, 0x1, 0x2, 0x3, 0x40, 0xe2, 0x1, 0x0, 0xf1, 0xfb, 0x9, 0x0}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})
}
