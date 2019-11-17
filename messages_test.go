package zstack

import (
	"github.com/shimmeringbee/unpi"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestMessageLibrary(t *testing.T) {
	t.Run("verifies that the message library returns false if message not found", func(t *testing.T) {
		ml := PopulateMessageLibrary()

		_, found := ml.GetByIdentifier(unpi.AREQ, unpi.SYS, 0xff)
		assert.False(t, found)

		type UnknownStruct struct{}

		_, found = ml.GetByObject(UnknownStruct{})
		assert.False(t, found)
	})

	t.Run("verifies that SYS_RESET_REQ is present", func(t *testing.T) {
		ml := PopulateMessageLibrary()

		expectedType := reflect.TypeOf(SysResetReq{})
		actualType, found := ml.GetByIdentifier(unpi.AREQ, unpi.SYS, SysResetReqCommand)

		assert.True(t, found)
		assert.Equal(t, expectedType, actualType)

		expectedIdentity := MessageIdentity{MessageType: unpi.AREQ, Subsystem: unpi.SYS, CommandID: SysResetReqCommand}
		actualIdentity, found := ml.GetByObject(SysResetReq{})

		assert.True(t, found)
		assert.Equal(t, expectedIdentity, actualIdentity)
	})

	t.Run("verifies that SYS_RESET_IND is present", func(t *testing.T) {
		ml := PopulateMessageLibrary()

		expectedType := reflect.TypeOf(SysResetInd{})
		actualType, found := ml.GetByIdentifier(unpi.AREQ, unpi.SYS, SysResetIndCommand)

		assert.True(t, found)
		assert.Equal(t, expectedType, actualType)

		expectedIdentity := MessageIdentity{MessageType: unpi.AREQ, Subsystem: unpi.SYS, CommandID: SysResetIndCommand}
		actualIdentity, found := ml.GetByObject(SysResetInd{})

		assert.True(t, found)
		assert.Equal(t, expectedIdentity, actualIdentity)
	})
}
