package zstack

import (
	"bytes"
	"errors"
	"github.com/shimmeringbee/unpi"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestZnp(t *testing.T) {
	t.Run("async outgoing request writes bytes", func(t *testing.T) {
		device := bytes.Buffer{}

		z := ZNP{
			writer:          &device,
			requestsChannel: make(chan OutgoingFrame, PermittedQueuedRequests),
			requestsEnd:     make(chan bool),
		}

		z.start()
		defer z.Stop()

		f := unpi.Frame{
			MessageType: unpi.AREQ,
			Subsystem:   unpi.ZDO,
			CommandID:   1,
			Payload:     []byte{0x78},
		}

		err := z.AsyncRequest(f)
		assert.NoError(t, err)

		expectedFrame := f.Marshall()
		actualFrame := device.Bytes()

		assert.Equal(t, expectedFrame, actualFrame)
	})

	t.Run("async outgoing request with non async request errors", func(t *testing.T) {
		z := ZNP{
			requestsChannel: make(chan OutgoingFrame, PermittedQueuedRequests),
			requestsEnd:     make(chan bool),
		}

		z.start()
		defer z.Stop()

		f := unpi.Frame{
			MessageType: unpi.SREQ,
			Subsystem:   unpi.ZDO,
			CommandID:   1,
			Payload:     []byte{0x78},
		}

		err := z.AsyncRequest(f)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, FrameNotAsynchronous))
	})

	t.Run("async outgoing request passes error back to caller", func(t *testing.T) {
		expectedError := errors.New("error")

		device := ControllableReaderWriter{
			Writer: func(p []byte) (n int, err error) {
				return 0, expectedError
			},
		}

		z := ZNP{
			writer:          &device,
			requestsChannel: make(chan OutgoingFrame, PermittedQueuedRequests),
			requestsEnd:     make(chan bool),
		}

		z.start()
		defer z.Stop()

		f := unpi.Frame{
			MessageType: unpi.AREQ,
			Subsystem:   unpi.ZDO,
			CommandID:   1,
			Payload:     []byte{0x78},
		}

		actualError := z.AsyncRequest(f)
		assert.Error(t, actualError)
		assert.Equal(t, expectedError, actualError)
	})

	t.Run("receive frames from unpi", func(t *testing.T) {
		device := bytes.Buffer{}

		expectedFrameOne := unpi.Frame{
			MessageType: 0,
			Subsystem:   unpi.ZDO,
			CommandID:   1,
			Payload:     []byte{0x78},
		}

		expectedFrameTwo := unpi.Frame{
			MessageType: 0,
			Subsystem:   unpi.SYS,
			CommandID:   2,
			Payload:     []byte{},
		}

		device.Write(expectedFrameOne.Marshall())
		device.Write(expectedFrameTwo.Marshall())

		z := ZNP{
			reader:          &device,
			requestsChannel: make(chan OutgoingFrame, PermittedQueuedRequests),
			requestsEnd:     make(chan bool),
		}

		z.start()
		defer z.Stop()

		frame, err := z.Receive()
		assert.NoError(t, err)
		assert.Equal(t, expectedFrameOne, frame)

		frame, err = z.Receive()
		assert.NoError(t, err)
		assert.Equal(t, expectedFrameTwo, frame)
	})

	t.Run("receive passes error back to caller", func(t *testing.T) {
		expectedError := errors.New("error")

		device := ControllableReaderWriter{
			Reader: func(p []byte) (n int, err error) {
				return 0, expectedError
			},
		}

		z := ZNP{
			reader:          &device,
			requestsChannel: make(chan OutgoingFrame, PermittedQueuedRequests),
			requestsEnd:     make(chan bool),
		}

		z.start()
		defer z.Stop()

		_, actualError := z.Receive()
		assert.Error(t, actualError)
		assert.Equal(t, expectedError, actualError)
	})

	t.Run("requesting a sync send with a non sync frame errors", func(t *testing.T) {
		z := ZNP{
			requestsChannel: make(chan OutgoingFrame, PermittedQueuedRequests),
			requestsEnd:     make(chan bool),
		}

		z.start()
		defer z.Stop()

		f := unpi.Frame{
			MessageType: unpi.AREQ,
			Subsystem:   unpi.ZDO,
			CommandID:   1,
			Payload:     []byte{0x78},
		}

		_, err := z.SyncRequest(f)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, FrameNotSynchronous))
	})

	t.Run("sync requests are sent to unpi and reply is read", func(t *testing.T) {
		responseFrame := unpi.Frame{
			MessageType: unpi.SRSP,
			Subsystem:   unpi.ZDO,
			CommandID:   1,
			Payload:     []byte{},
		}
		responseBytes := responseFrame.Marshall()

		beenWrittenBuffer := bytes.Buffer{}
		toBeReadBuffer := bytes.Buffer{}

		device := ControllableReaderWriter{
			Writer: func(p []byte) (n int, err error) {
				beenWrittenBuffer.Write(p)
				toBeReadBuffer.Write(responseBytes)
				return len(p), nil
			},
			Reader: func(p []byte) (n int, err error) {
				return toBeReadBuffer.Read(p)
			},
		}

		z := ZNP{
			writer:          &device,
			reader:          &device,
			requestsChannel: make(chan OutgoingFrame, PermittedQueuedRequests),
			requestsEnd:     make(chan bool),
		}

		z.start()
		defer z.Stop()

		f := unpi.Frame{
			MessageType: unpi.SREQ,
			Subsystem:   unpi.ZDO,
			CommandID:   1,
			Payload:     []byte{0x78},
		}

		actualResponseFrame, err := z.SyncRequest(f)
		assert.NoError(t, err)

		expectedFrame := f.Marshall()
		actualFrame := beenWrittenBuffer.Bytes()

		assert.Equal(t, expectedFrame, actualFrame)
		assert.Equal(t, responseFrame, actualResponseFrame)
	})

	t.Run("sync requests are sent to unpi and reply errors are handled", func(t *testing.T) {
		expectedError := errors.New("error")

		device := ControllableReaderWriter{
			Writer: func(p []byte) (n int, err error) {
				return len(p), nil
			},
			Reader: func(p []byte) (n int, err error) {
				return 0, expectedError
			},
		}

		z := ZNP{
			writer:          &device,
			reader:          &device,
			requestsChannel: make(chan OutgoingFrame, PermittedQueuedRequests),
			requestsEnd:     make(chan bool),
		}

		z.start()
		defer z.Stop()

		f := unpi.Frame{
			MessageType: unpi.SREQ,
			Subsystem:   unpi.ZDO,
			CommandID:   1,
			Payload:     []byte{0x78},
		}

		_, err := z.SyncRequest(f)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("sync outgoing request passes error during write back to caller", func(t *testing.T) {
		expectedError := errors.New("error")

		device := ControllableReaderWriter{
			Writer: func(p []byte) (n int, err error) {
				return 0, expectedError
			},
		}

		z := ZNP{
			writer:          &device,
			requestsChannel: make(chan OutgoingFrame, PermittedQueuedRequests),
			requestsEnd:     make(chan bool),
		}

		z.start()
		defer z.Stop()

		f := unpi.Frame{
			MessageType: unpi.SREQ,
			Subsystem:   unpi.ZDO,
			CommandID:   1,
			Payload:     []byte{0x78},
		}

		_, actualError := z.SyncRequest(f)
		assert.Error(t, actualError)
		assert.Equal(t, expectedError, actualError)
	})
}

type ControllableReaderWriter struct {
	Reader func(p []byte) (n int, err error)
	Writer func(p []byte) (n int, err error)
}

func (c *ControllableReaderWriter) Read(p []byte) (n int, err error) {
	return c.Reader(p)
}

func (c *ControllableReaderWriter) Write(p []byte) (n int, err error) {
	return c.Writer(p)
}
