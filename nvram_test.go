package zstack

import (
	"context"
	"errors"
	"github.com/shimmeringbee/bytecodec"
	"github.com/shimmeringbee/zigbee"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func Test_readNVRAM(t *testing.T) {
	t.Run("verifies that a request response is made to unpi", func(t *testing.T) {
		mrr := &ReturningMockRequestResponse{}

		z := ZStack{requestResponder: mrr}

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		val := &ZCDNVLogicalType{}
		err := z.readNVRAM(ctx, val)

		assert.NoError(t, err)
		assert.Equal(t, zigbee.EndDevice, val.LogicalType)
	})

	t.Run("verifies that read requests that fail raise an error", func(t *testing.T) {
		mrr := &ReadFailingMockRequestResponse{}

		z := ZStack{requestResponder: mrr}

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		err := z.readNVRAM(ctx, &ZCDNVLogicalType{})

		assert.Error(t, err)
	})

	t.Run("verifies that a request response with errors is raised", func(t *testing.T) {
		mrr := new(MockRequestResponder)

		mrr.On("RequestResponse", mock.Anything, SysOSALNVRead{
			NVItemID: ZCDNVLogicalTypeID,
			Offset:   0,
		}, &SysOSALNVReadReply{}).Return(errors.New("context expired"))

		z := ZStack{requestResponder: mrr}

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		err := z.readNVRAM(ctx, &ZCDNVLogicalType{})

		mrr.AssertExpectations(t)
		assert.Error(t, err)
	})

	t.Run("verifies that unknown structure raises an error", func(t *testing.T) {
		mrr := new(MockRequestResponder)

		z := ZStack{requestResponder: mrr}

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		err := z.readNVRAM(ctx, struct{}{})

		mrr.AssertExpectations(t)
		assert.Error(t, err)
	})
}

func Test_writeNVRAM(t *testing.T) {
	t.Run("verifies that a request response is made to unpi", func(t *testing.T) {
		mrr := new(MockRequestResponder)

		mrr.On("RequestResponse", mock.Anything, SysOSALNVWrite{
			NVItemID: ZCDNVLogicalTypeID,
			Offset:   0,
			Value:    []byte{0x02},
		}, &SysOSALNVWriteReply{}).Return(nil)

		z := ZStack{requestResponder: mrr}

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		err := z.writeNVRAM(ctx, ZCDNVLogicalType{LogicalType: zigbee.EndDevice})

		mrr.AssertExpectations(t)
		assert.NoError(t, err)
	})

	t.Run("verifies that write requests that fail raise an error", func(t *testing.T) {
		mrr := &WriteFailingMockRequestResponse{}

		z := ZStack{requestResponder: mrr}

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		err := z.writeNVRAM(ctx, ZCDNVLogicalType{LogicalType: zigbee.EndDevice})

		assert.Error(t, err)
	})

	t.Run("verifies that a request response with errors is raised", func(t *testing.T) {
		mrr := new(MockRequestResponder)

		mrr.On("RequestResponse", mock.Anything, SysOSALNVWrite{
			NVItemID: ZCDNVLogicalTypeID,
			Offset:   0,
			Value:    []byte{0x02},
		}, &SysOSALNVWriteReply{}).Return(errors.New("context expired"))

		z := ZStack{requestResponder: mrr}

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		err := z.writeNVRAM(ctx, ZCDNVLogicalType{LogicalType: zigbee.EndDevice})

		mrr.AssertExpectations(t)
		assert.Error(t, err)
	})

	t.Run("verifies that unknown structure raises an error", func(t *testing.T) {
		mrr := new(MockRequestResponder)

		z := ZStack{requestResponder: mrr}

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		err := z.writeNVRAM(ctx, struct{}{})

		mrr.AssertExpectations(t)
		assert.Error(t, err)
	})
}

type WriteFailingMockRequestResponse struct{}

func (m *WriteFailingMockRequestResponse) RequestResponse(ctx context.Context, req interface{}, resp interface{}) error {
	response, ok := resp.(*SysOSALNVWriteReply)

	if !ok {
		panic("incorrect type passed to mock")
	}

	response.Status = 0x01

	return nil
}

type ReadFailingMockRequestResponse struct{}

func (m *ReadFailingMockRequestResponse) RequestResponse(ctx context.Context, req interface{}, resp interface{}) error {
	response, ok := resp.(*SysOSALNVReadReply)

	if !ok {
		panic("incorrect type passed to mock")
	}

	response.Status = 0x01

	return nil
}

type ReturningMockRequestResponse struct{}

func (m *ReturningMockRequestResponse) RequestResponse(ctx context.Context, req interface{}, resp interface{}) error {
	response, ok := resp.(*SysOSALNVReadReply)

	if !ok {
		panic("incorrect type passed to mock")
	}

	response.Status = ZSuccess
	response.Value = []byte{0x02}

	return nil
}

func Test_NVRAMStructs(t *testing.T) {
	t.Run("ZCDNVStartUpOption", func(t *testing.T) {
		s := ZCDNVStartUpOption{
			StartOption: 3,
		}

		actualBytes, err := bytecodec.Marshal(s)

		expectedBytes := []byte{0x03}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("ZCDNVLogicalType", func(t *testing.T) {
		s := ZCDNVLogicalType{
			LogicalType: zigbee.Router,
		}

		actualBytes, err := bytecodec.Marshal(s)

		expectedBytes := []byte{0x01}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("ZCDNVSecurityMode", func(t *testing.T) {
		s := ZCDNVSecurityMode{
			Enabled: 1,
		}

		actualBytes, err := bytecodec.Marshal(s)

		expectedBytes := []byte{0x01}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("ZCDNVPreCfgKeysEnable", func(t *testing.T) {
		s := ZCDNVPreCfgKeysEnable{
			Enabled: 1,
		}

		actualBytes, err := bytecodec.Marshal(s)

		expectedBytes := []byte{0x01}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("ZCDNVPreCfgKey", func(t *testing.T) {
		s := ZCDNVPreCfgKey{
			NetworkKey: [16]byte{0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03},
		}

		actualBytes, err := bytecodec.Marshal(s)

		expectedBytes := []byte{0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("ZCDNVZDODirectCB", func(t *testing.T) {
		s := ZCDNVZDODirectCB{
			Enabled: 1,
		}

		actualBytes, err := bytecodec.Marshal(s)

		expectedBytes := []byte{0x01}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("ZCDNVChanList", func(t *testing.T) {
		s := ZCDNVChanList{
			Channels: [4]byte{0x00, 0x00, 0x00, 0x03},
		}

		actualBytes, err := bytecodec.Marshal(s)

		expectedBytes := []byte{0x00, 0x00, 0x00, 0x03}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("ZCDNVPANID", func(t *testing.T) {
		s := ZCDNVPANID{
			PANID: zigbee.PANID(0x0102),
		}

		actualBytes, err := bytecodec.Marshal(s)

		expectedBytes := []byte{0x02, 0x01}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("ZCDNVExtPANID", func(t *testing.T) {
		s := ZCDNVExtPANID{
			ExtendedPANID: zigbee.ExtendedPANID(0x0102030405060708),
		}

		actualBytes, err := bytecodec.Marshal(s)

		expectedBytes := []byte{0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("ZCDNVUseDefaultTCLK", func(t *testing.T) {
		s := ZCDNVUseDefaultTCLK{
			Enabled: 1,
		}

		actualBytes, err := bytecodec.Marshal(s)

		expectedBytes := []byte{0x01}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("ZCDNVTCLKTableStart", func(t *testing.T) {
		s := ZCDNVTCLKTableStart{
			Address:        zigbee.IEEEAddress(0x0807060504030201),
			NetworkKey:     [16]byte{0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03},
			TXFrameCounter: 123456,
			RXFrameCounter: 654321,
		}

		actualBytes, err := bytecodec.Marshal(s)

		expectedBytes := []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x0, 0x1, 0x2, 0x3, 0x0, 0x1, 0x2, 0x3, 0x0, 0x1, 0x2, 0x3, 0x0, 0x1, 0x2, 0x3, 0x40, 0xe2, 0x1, 0x0, 0xf1, 0xfb, 0x9, 0x0}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})
}
