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
			writer: &device,
		}

		z.start()
		defer z.Stop()

		f := unpi.Frame{
			MessageType: 0,
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

	t.Run("async outgoing request passes error back to caller", func(t *testing.T) {
		expectedError := errors.New("error")

		device := ControllableReaderWriter{
			Writer: func(p []byte) (n int, err error) {
				return 0, expectedError
			},
		}

		z := ZNP{
			writer: &device,
		}

		z.start()
		defer z.Stop()

		f := unpi.Frame{
			MessageType: 0,
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
			reader: &device,
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
			reader: &device,
		}

		z.start()
		defer z.Stop()

		_, actualError := z.Receive()
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
