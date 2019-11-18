package zstack

import (
	"bytes"
	"context"
	"github.com/shimmeringbee/bytecodec"
	"github.com/shimmeringbee/unpi"
	"github.com/stretchr/testify/assert"
	"io"
	"reflect"
	"testing"
	"time"
)

func TestMessageRequestResponse(t *testing.T) {
	t.Run("verifies send with receive on synchronous messages are handled", func(t *testing.T) {
		sentMessage := SysOSALNVWrite{
			NVItemID: 1,
			Offset:   0,
			Value:    []byte{0x02, 0x03},
		}

		expectedMessage := SysOSALNVWriteResponse{
			Status: 1,
		}

		payloadBytes, err := bytecodec.Marshall(expectedMessage)
		assert.NoError(t, err)

		frame := unpi.Frame{
			MessageType: unpi.SRSP,
			Subsystem:   unpi.SYS,
			CommandID:   0x09,
			Payload:     payloadBytes,
		}

		responseBytes := frame.Marshall()
		writtenBytes := bytes.Buffer{}

		r, w := io.Pipe()
		defer w.Close()

		device := ControllableReaderWriter{
			Writer: func(p []byte) (n int, err error) {
				go func() { w.Write(responseBytes) }()
				return writtenBytes.Write(p)
			},
			Reader: func(p []byte) (n int, err error) {
				return r.Read(p)
			},
		}

		z := New(&device, &device)
		defer z.Stop()

		actualReceivedMessage := SysOSALNVWriteResponse{}

		ctx, _ := context.WithTimeout(context.Background(), 10*time.Millisecond)
		err = z.MessageRequestResponse(ctx, sentMessage, &actualReceivedMessage)

		assert.NoError(t, err)
		assert.Equal(t, expectedMessage, actualReceivedMessage)

		frame, _ = unpi.UnmarshallFrame(writtenBytes.Bytes())

		assert.Equal(t, unpi.SREQ, frame.MessageType)
		assert.Equal(t, unpi.SYS, frame.Subsystem)
		assert.Equal(t, SysOSALNVWriteRequestID, frame.CommandID)

		actualSentMessage := SysOSALNVWrite{}
		err = bytecodec.Unmarshall(frame.Payload, &actualSentMessage)
		assert.NoError(t, err)

		assert.Equal(t, sentMessage, actualSentMessage)
	})

	t.Run("verifies send with receive on asynchronous messages are handled", func(t *testing.T) {
		sentMessage := SysResetReq{
			ResetType: Soft,
		}

		expectedMessage := SysResetInd{
			Reason:            External,
			TransportRevision: 2,
			ProductID:         1,
			MajorRelease:      2,
			MinorRelease:      3,
			HardwareRevision:  1,
		}

		payloadBytes, err := bytecodec.Marshall(expectedMessage)
		assert.NoError(t, err)

		frame := unpi.Frame{
			MessageType: unpi.AREQ,
			Subsystem:   unpi.SYS,
			CommandID:   SysResetIndicationCommandID,
			Payload:     payloadBytes,
		}

		responseBytes := frame.Marshall()
		writtenBytes := bytes.Buffer{}

		r, w := io.Pipe()
		defer w.Close()

		device := ControllableReaderWriter{
			Writer: func(p []byte) (n int, err error) {
				go func() { w.Write(responseBytes) }()
				return writtenBytes.Write(p)
			},
			Reader: func(p []byte) (n int, err error) {
				return r.Read(p)
			},
		}

		z := New(&device, &device)
		defer z.Stop()

		actualReceivedMessage := SysResetInd{}

		ctx, _ := context.WithTimeout(context.Background(), 10*time.Millisecond)
		err = z.MessageRequestResponse(ctx, sentMessage, &actualReceivedMessage)

		assert.NoError(t, err)
		assert.Equal(t, expectedMessage, actualReceivedMessage)

		frame, _ = unpi.UnmarshallFrame(writtenBytes.Bytes())

		assert.Equal(t, unpi.AREQ, frame.MessageType)
		assert.Equal(t, unpi.SYS, frame.Subsystem)
		assert.Equal(t, SysResetRequestID, frame.CommandID)

		actualSentMessage := SysResetReq{}
		err = bytecodec.Unmarshall(frame.Payload, &actualSentMessage)
		assert.NoError(t, err)

		assert.Equal(t, sentMessage, actualSentMessage)
	})
}

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
		actualType, found := ml.GetByIdentifier(unpi.AREQ, unpi.SYS, SysResetIndicationCommandID)

		assert.True(t, found)
		assert.Equal(t, expectedType, actualType)

		expectedIdentity := MessageIdentity{MessageType: unpi.AREQ, Subsystem: unpi.SYS, CommandID: SysResetIndicationCommandID}
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
