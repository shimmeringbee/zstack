package zstack // import "github.com/shimmeringbee/zstack"

import (
	"errors"
	"github.com/shimmeringbee/unpi"
	"io"
)

type ZNP struct {
	reader          io.Reader
	writer          io.Writer
	requestsChannel chan OutgoingFrame
	requestsEnd     chan bool
}

const PermittedQueuedRequests int = 50

type OutgoingFrame struct {
	Frame        unpi.Frame
	ErrorChannel chan error
}

func New(device io.ReadWriter) *ZNP {
	z := &ZNP{
		reader:          device,
		writer:          device,
		requestsChannel: make(chan OutgoingFrame, PermittedQueuedRequests),
		requestsEnd:     make(chan bool),
	}

	z.start()

	return z
}

func (z *ZNP) start() {
	go z.handleRequests()
}

func (z *ZNP) handleRequests() {
	for {
		select {
		case outgoing := <-z.requestsChannel:
			outgoing.ErrorChannel <- unpi.Write(z.writer, outgoing.Frame)
		case <-z.requestsEnd:
			return
		}
	}
}

func (z *ZNP) Stop() {
	z.requestsEnd <- true
}

func (z *ZNP) writeFrame(frame unpi.Frame) error {
	errCh := make(chan error)

	z.requestsChannel <- OutgoingFrame{
		Frame:        frame,
		ErrorChannel: errCh,
	}

	return <-errCh
}

var FrameNotAsynchronous = errors.New("frame not asynchronous")
var FrameNotSynchronous = errors.New("frame not synchronous")

func (z *ZNP) AsyncRequest(frame unpi.Frame) error {
	if frame.MessageType != unpi.AREQ {
		return FrameNotAsynchronous
	}

	return z.writeFrame(frame)
}

func (z *ZNP) SyncRequest(frame unpi.Frame) (unpi.Frame, error) {
	if frame.MessageType != unpi.SREQ {
		return unpi.Frame{}, FrameNotSynchronous
	}

	if err := z.writeFrame(frame); err != nil {
		return unpi.Frame{}, err
	}

	return unpi.Read(z.reader)
}

func (z *ZNP) Receive() (unpi.Frame, error) {
	return unpi.Read(z.reader)
}
