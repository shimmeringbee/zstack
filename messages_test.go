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
		actualType, found := ml.GetByIdentifier(unpi.AREQ, unpi.SYS, SysResetRequestID)

		assert.True(t, found)
		assert.Equal(t, expectedType, actualType)

		expectedIdentity := MessageIdentity{MessageType: unpi.AREQ, Subsystem: unpi.SYS, CommandID: SysResetRequestID}
		actualIdentity, found := ml.GetByObject(SysResetReq{})

		assert.True(t, found)
		assert.Equal(t, expectedIdentity, actualIdentity)
	})

	t.Run("verifies that SYS_RESET_IND is present", func(t *testing.T) {
		ml := PopulateMessageLibrary()

		expectedType := reflect.TypeOf(SysResetInd{})
		actualType, found := ml.GetByIdentifier(unpi.AREQ, unpi.SYS, SysResetIndidcationCommandID)

		assert.True(t, found)
		assert.Equal(t, expectedType, actualType)

		expectedIdentity := MessageIdentity{MessageType: unpi.AREQ, Subsystem: unpi.SYS, CommandID: SysResetIndidcationCommandID}
		actualIdentity, found := ml.GetByObject(SysResetInd{})

		assert.True(t, found)
		assert.Equal(t, expectedIdentity, actualIdentity)
	})

	t.Run("verifies that SYS_OSAL_NV_WRITE is present", func(t *testing.T) {
		ml := PopulateMessageLibrary()

		expectedType := reflect.TypeOf(SysOSALNVWrite{})
		actualType, found := ml.GetByIdentifier(unpi.SREQ, unpi.SYS, SysOSALNVWriteRequestID)

		assert.True(t, found)
		assert.Equal(t, expectedType, actualType)

		expectedIdentity := MessageIdentity{MessageType: unpi.SREQ, Subsystem: unpi.SYS, CommandID: SysOSALNVWriteRequestID}
		actualIdentity, found := ml.GetByObject(SysOSALNVWrite{})

		assert.True(t, found)
		assert.Equal(t, expectedIdentity, actualIdentity)
	})

	t.Run("verifies that SYS_OSAL_NV_WRITE response is present", func(t *testing.T) {
		ml := PopulateMessageLibrary()

		expectedType := reflect.TypeOf(SysOSALNVWriteResponse{})
		actualType, found := ml.GetByIdentifier(unpi.SRSP, unpi.SYS, SysOSALNVWriteResponseID)

		assert.True(t, found)
		assert.Equal(t, expectedType, actualType)

		expectedIdentity := MessageIdentity{MessageType: unpi.SRSP, Subsystem: unpi.SYS, CommandID: SysOSALNVWriteResponseID}
		actualIdentity, found := ml.GetByObject(SysOSALNVWriteResponse{})

		assert.True(t, found)
		assert.Equal(t, expectedIdentity, actualIdentity)
	})
}
