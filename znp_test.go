package zstack

import (
	"bytes"
	"context"
	"errors"
	"github.com/shimmeringbee/unpi"
	"github.com/stretchr/testify/assert"
	"io"
	"sync"
	"testing"
	"time"
)

func TestZnp(t *testing.T) {
	t.Run("async outgoing request writes bytes", func(t *testing.T) {
		writer := bytes.Buffer{}
		reader := EmptyReader{
			End: make(chan bool),
		}
		defer reader.Done()

		z := New(&reader, &writer)
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
		actualFrame := writer.Bytes()

		assert.Equal(t, expectedFrame, actualFrame)
	})

	t.Run("async outgoing request with non async request errors", func(t *testing.T) {
		writer := bytes.Buffer{}
		reader := EmptyReader{
			End: make(chan bool),
		}
		defer reader.Done()

		z := New(&reader, &writer)
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

		writer := ControllableReaderWriter{
			Writer: func(p []byte) (n int, err error) {
				return 0, expectedError
			},
		}
		reader := EmptyReader{
			End: make(chan bool),
		}
		defer reader.Done()

		z := New(&reader, &writer)
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
		reader := bytes.Buffer{}
		writer := bytes.Buffer{}

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

		reader.Write(expectedFrameOne.Marshall())
		reader.Write(expectedFrameTwo.Marshall())

		z := New(&reader, &writer)
		defer z.Stop()

		frame, err := z.Receive()
		assert.NoError(t, err)
		assert.Equal(t, expectedFrameOne, frame)

		frame, err = z.Receive()
		assert.NoError(t, err)
		assert.Equal(t, expectedFrameTwo, frame)
	})

	t.Run("requesting a sync send with a non sync frame errors", func(t *testing.T) {
		reader := bytes.Buffer{}
		writer := bytes.Buffer{}

		z := New(&reader, &writer)
		defer z.Stop()

		f := unpi.Frame{
			MessageType: unpi.AREQ,
			Subsystem:   unpi.ZDO,
			CommandID:   1,
			Payload:     []byte{0x78},
		}

		_, err := z.SyncRequest(context.Background(), f)
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
		r, w := io.Pipe()

		device := ControllableReaderWriter{
			Writer: func(p []byte) (n int, err error) {
				beenWrittenBuffer.Write(p)
				go func() { w.Write(responseBytes) }()
				return len(p), nil
			},
			Reader: func(p []byte) (n int, err error) {
				return r.Read(p)
			},
		}

		z := New(&device, &device)
		defer z.Stop()

		f := unpi.Frame{
			MessageType: unpi.SREQ,
			Subsystem:   unpi.ZDO,
			CommandID:   1,
			Payload:     []byte{0x78},
		}

		actualResponseFrame, err := z.SyncRequest(context.Background(), f)
		assert.NoError(t, err)

		expectedFrame := f.Marshall()
		actualFrame := beenWrittenBuffer.Bytes()

		assert.Equal(t, expectedFrame, actualFrame)
		assert.Equal(t, responseFrame, actualResponseFrame)
	})

	t.Run("sync outgoing request passes error during write back to caller", func(t *testing.T) {
		expectedError := errors.New("error")

		reader := bytes.Buffer{}
		writer := ControllableReaderWriter{
			Writer: func(p []byte) (n int, err error) {
				return 0, expectedError
			},
		}

		z := New(&reader, &writer)
		defer z.Stop()

		f := unpi.Frame{
			MessageType: unpi.SREQ,
			Subsystem:   unpi.ZDO,
			CommandID:   1,
			Payload:     []byte{0x78},
		}

		_, actualError := z.SyncRequest(context.Background(), f)
		assert.Error(t, actualError)
		assert.Equal(t, expectedError, actualError)
	})

	t.Run("sync outgoing context cancellation causes function to error", func(t *testing.T) {
		reader := EmptyReader{End: make(chan bool)}
		writer := bytes.Buffer{}

		z := New(&reader, &writer)
		defer z.Stop()

		f := unpi.Frame{
			MessageType: unpi.SREQ,
			Subsystem:   unpi.ZDO,
			CommandID:   1,
			Payload:     []byte{0x78},
		}

		ctx, _ := context.WithTimeout(context.Background(), 1*time.Microsecond)

		_, actualError := z.SyncRequest(ctx, f)
		assert.Error(t, actualError)
		assert.Equal(t, SyncRequestContextCancelled, actualError)
	})

	t.Run("wait for frame responds to multiple listeners when frame matches", func(t *testing.T) {
		r, w := io.Pipe()
		defer w.Close()

		device := ControllableReaderWriter{
			Writer: func(p []byte) (n int, err error) {
				return len(p), nil
			},
			Reader: func(p []byte) (n int, err error) {
				return r.Read(p)
			},
		}

		z := New(&device, &device)
		defer z.Stop()

		expectedFrame := unpi.Frame{
			MessageType: unpi.AREQ,
			Subsystem:   unpi.SYS,
			CommandID:   0x20,
			Payload:     []byte{},
		}

		wg := &sync.WaitGroup{}

		for i := 0; i < 2; i++ {
			wg.Add(1)

			go func() {
				ctxWithTimeout, _ := context.WithTimeout(context.Background(), 20*time.Millisecond)
				actualFrame, err := z.WaitForFrame(ctxWithTimeout, unpi.AREQ, unpi.SYS, 0x20)

				assert.NoError(t, err)
				assert.Equal(t, expectedFrame, actualFrame)

				wg.Done()
			}()
		}

		time.Sleep(10 * time.Millisecond)
		data := expectedFrame.Marshall()
		w.Write(data)

		wg.Wait()
	})

	t.Run("wait for frame respects context timeout", func(t *testing.T) {
		reader := EmptyReader{End: make(chan bool)}
		defer reader.Done()

		writer := bytes.Buffer{}

		z := New(&reader, &writer)
		defer z.Stop()

		ctxWithTimeout, _ := context.WithTimeout(context.Background(), 25*time.Millisecond)
		_, err := z.WaitForFrame(ctxWithTimeout, unpi.AREQ, unpi.SYS, 0x20)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, WaitForFrameContextCancelled))
	})

	t.Run("wait for frame ignores unrelated frames", func(t *testing.T) {
		r, w := io.Pipe()

		device := ControllableReaderWriter{
			Writer: func(p []byte) (n int, err error) {
				return len(p), nil
			},
			Reader: func(p []byte) (n int, err error) {
				return r.Read(p)
			},
		}

		z := New(&device, &device)
		defer z.Stop()

		expectedFrame := unpi.Frame{
			MessageType: unpi.AREQ,
			Subsystem:   unpi.SYS,
			CommandID:   0x21,
			Payload:     nil,
		}

		go func() {
			time.Sleep(10 * time.Millisecond)
			data := expectedFrame.Marshall()
			w.Write(data)
		}()

		ctxWithTimeout, _ := context.WithTimeout(context.Background(), 20*time.Millisecond)
		_, err := z.WaitForFrame(ctxWithTimeout, unpi.AREQ, unpi.SYS, 0x20)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, WaitForFrameContextCancelled))
	})
}

type EmptyReader struct {
	End chan bool
}

func (e *EmptyReader) Done() {
	e.End <- true
}

func (e *EmptyReader) Read(p []byte) (n int, err error) {
	<-e.End

	return 0, io.EOF
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
