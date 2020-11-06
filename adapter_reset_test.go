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

func Test_resetMessages(t *testing.T) {
	t.Run("verify SysResetReq marshals", func(t *testing.T) {
		req := SysResetReq{ResetType: Soft}

		data, err := bytecodec.Marshal(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x01}, data)
	})

	t.Run("verify SysResetInd marshals", func(t *testing.T) {
		req := SysResetInd{
			Reason: External,
			Version: Version{
				TransportRevision: 2,
				ProductID:         1,
				MajorRelease:      2,
				MinorRelease:      4,
				HardwareRevision:  1,
			},
		}

		data, err := bytecodec.Marshal(req)

		assert.NoError(t, err)
		assert.Equal(t, []byte{0x01, 0x02, 0x01, 0x02, 0x04, 0x01}, data)
	})
}

type MockRequestResponder struct {
	mock.Mock
}

func (m *MockRequestResponder) RequestResponse(ctx context.Context, req interface{}, resp interface{}) error {
	args := m.Called(ctx, req, resp)
	return args.Error(0)
}

func Test_resetAdapter(t *testing.T) {
	t.Run("verifies that a request response is made to unpi", func(t *testing.T) {
		mrr := new(MockRequestResponder)
		defer mrr.AssertExpectations(t)

		expectedVersion := Version{
			TransportRevision: 1,
			ProductID:         2,
			MajorRelease:      3,
			MinorRelease:      4,
			HardwareRevision:  5,
		}

		mrr.On("RequestResponse", mock.Anything, SysResetReq{ResetType: Soft}, &SysResetInd{}).Return(nil).Run(func(args mock.Arguments) {
			sysResetInd := args.Get(2).(*SysResetInd)
			sysResetInd.Version = expectedVersion
		})

		z := ZStack{requestResponder: mrr}

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		actualVersion, err := z.resetAdapter(ctx, Soft)
		assert.NoError(t, err)
		assert.Equal(t, expectedVersion, actualVersion)
	})

	t.Run("verifies that a request response with errors is raised", func(t *testing.T) {
		mrr := new(MockRequestResponder)
		defer mrr.AssertExpectations(t)

		mrr.On("RequestResponse", mock.Anything, SysResetReq{ResetType: Soft}, &SysResetInd{}).Return(errors.New("context expired"))

		z := ZStack{requestResponder: mrr}

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		_, err := z.resetAdapter(ctx, Soft)
		assert.Error(t, err)
	})
}
