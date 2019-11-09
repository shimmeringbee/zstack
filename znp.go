package zstack // import "github.com/shimmeringbee/zstack"

import (
	"errors"
	"github.com/shimmeringbee/unpi"
	"io"
	"log"
)

type ZNP struct {
	reader               io.Reader
	writer               io.Writer
	requestsChannel      chan OutgoingFrame
	requestsEnd          chan bool
	syncReceivingChannel chan chan unpi.Frame
	receivingChannel     chan unpi.Frame
	receivingEnd         chan bool
}

const PermittedQueuedRequests int = 50

type OutgoingFrame struct {
	Frame        unpi.Frame
	ErrorChannel chan error
}

func New(reader io.Reader, writer io.Writer) *ZNP {
	z := &ZNP{
		reader:               reader,
		writer:               writer,
		requestsChannel:      make(chan OutgoingFrame, PermittedQueuedRequests),
		requestsEnd:          make(chan bool),
		syncReceivingChannel: make(chan chan unpi.Frame, PermittedQueuedRequests),
		receivingChannel:     make(chan unpi.Frame, PermittedQueuedRequests),
		receivingEnd:         make(chan bool, 1),
	}

	z.start()

	return z
}

func (z *ZNP) start() {
	go z.handleRequests()
	go z.handleReceiving()
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

func (z *ZNP) handleReceiving() {
	for {
		frame, err := unpi.Read(z.reader)

		if err != nil {
			log.Printf("unpi read failed: %v\n", err)

			if errors.Is(err, io.EOF) {
				return
			}
		} else {
			if frame.MessageType != unpi.SRSP {
				z.receivingChannel <- frame
			} else {
				select {
				case syncChannel := <-z.syncReceivingChannel:
					syncChannel <- frame
				default:
					log.Println("received synchronous response, but no receivers in channel")
				}
			}
		}

		select {
		case <-z.receivingEnd:
			return
		default:
		}
	}
}

func (z *ZNP) Stop() {
	z.requestsEnd <- true
	z.receivingEnd <- true
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

	syncResponse := make(chan unpi.Frame)

	if err := z.writeFrame(frame); err != nil {
		return unpi.Frame{}, err
	}

	z.syncReceivingChannel <- syncResponse

	return <-syncResponse, nil
}

func (z *ZNP) Receive() (unpi.Frame, error) {
	return <-z.receivingChannel, nil
}
