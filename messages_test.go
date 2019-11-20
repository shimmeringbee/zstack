package zstack

import (
	. "github.com/shimmeringbee/unpi"
	. "github.com/shimmeringbee/unpi/library"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func Test_registerMessages(t *testing.T) {
	ml := NewLibrary()
	registerMessages(ml)

	t.Run("SysResetReq", func(t *testing.T) {
		identity, found := ml.GetByObject(&SysResetReq{})

		assert.True(t, found)
		assert.Equal(t, AREQ, identity.MessageType)
		assert.Equal(t, SYS, identity.Subsystem)
		assert.Equal(t, uint8(0x00), identity.CommandID)

		ty, found := ml.GetByIdentifier(AREQ, SYS, 0x00)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(SysResetReq{}), ty)
	})

	t.Run("SysResetInd", func(t *testing.T) {
		identity, found := ml.GetByObject(&SysResetInd{})

		assert.True(t, found)
		assert.Equal(t, AREQ, identity.MessageType)
		assert.Equal(t, SYS, identity.Subsystem)
		assert.Equal(t, uint8(0x80), identity.CommandID)

		ty, found := ml.GetByIdentifier(AREQ, SYS, 0x80)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(SysResetInd{}), ty)
	})

	t.Run("SysOSALNVWriteReq", func(t *testing.T) {
		identity, found := ml.GetByObject(&SysOSALNVWriteReq{})

		assert.True(t, found)
		assert.Equal(t, SREQ, identity.MessageType)
		assert.Equal(t, SYS, identity.Subsystem)
		assert.Equal(t, uint8(0x09), identity.CommandID)

		ty, found := ml.GetByIdentifier(SREQ, SYS, 0x09)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(SysOSALNVWriteReq{}), ty)
	})

	t.Run("SysOSALNVWriteResp", func(t *testing.T) {
		identity, found := ml.GetByObject(&SysOSALNVWriteResp{})

		assert.True(t, found)
		assert.Equal(t, SRSP, identity.MessageType)
		assert.Equal(t, SYS, identity.Subsystem)
		assert.Equal(t, uint8(0x09), identity.CommandID)

		ty, found := ml.GetByIdentifier(SRSP, SYS, 0x09)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(SysOSALNVWriteResp{}), ty)
	})
}