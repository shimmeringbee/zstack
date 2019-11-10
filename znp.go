package zstack // import "github.com/shimmeringbee/zstack"

import (
	"errors"
	"github.com/shimmeringbee/unpi"
	"io"
	"log"
	"sync"
)

type ZNP struct {
	reader io.Reader
	writer io.Writer

	requestsChannel chan OutgoingFrame
	requestsEnd     chan bool

	syncReceivingMutex    *sync.Mutex
	syncReceivingChannel  chan unpi.Frame
	asyncReceivingChannel chan unpi.Frame
	receivingEnd          chan bool
}

const PermittedQueuedRequests int = 50

type OutgoingFrame struct {
	Frame        unpi.Frame
	ErrorChannel chan error
}

func New(reader io.Reader, writer io.Writer) *ZNP {
	z := &ZNP{
		reader: reader,
		writer: writer,

		requestsChannel: make(chan OutgoingFrame, PermittedQueuedRequests),
		requestsEnd:     make(chan bool),

		syncReceivingMutex:    &sync.Mutex{},
		syncReceivingChannel:  make(chan unpi.Frame),
		asyncReceivingChannel: make(chan unpi.Frame, PermittedQueuedRequests),
		receivingEnd:          make(chan bool, 1),
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
				z.asyncReceivingChannel <- frame
			} else {
				select {
				case z.syncReceivingChannel <- frame:
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

	z.syncReceivingMutex.Lock()
	defer z.syncReceivingMutex.Unlock()

	if err := z.writeFrame(frame); err != nil {
		return unpi.Frame{}, err
	}

	return <-z.syncReceivingChannel, nil
}

func (z *ZNP) Receive() (unpi.Frame, error) {
	return <-z.asyncReceivingChannel, nil
}
