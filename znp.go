package zstack // import "github.com/shimmeringbee/zstack"

import (
	"context"
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

	waitForRequestsMutex *sync.Mutex
	waitForRequests      map[WaitFrameRequest]bool

	messageLibrary MessageLibrary
}

const PermittedQueuedRequests int = 50

type OutgoingFrame struct {
	Frame        unpi.Frame
	ErrorChannel chan error
}

type WaitFrameRequest struct {
	MessageType unpi.MessageType
	SubSystem   unpi.Subsystem
	CommandID   byte
	Response    chan unpi.Frame
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

		waitForRequestsMutex: &sync.Mutex{},
		waitForRequests:      map[WaitFrameRequest]bool{},

		messageLibrary: PopulateMessageLibrary(),
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
				z.serviceWaitForRequests(frame)
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

func (z *ZNP) serviceWaitForRequests(frame unpi.Frame) {
	z.waitForRequestsMutex.Lock()
	defer z.waitForRequestsMutex.Unlock()

	for req, _ := range z.waitForRequests {
		if req.MessageType == frame.MessageType &&
			req.SubSystem == frame.Subsystem &&
			req.CommandID == frame.CommandID {

			select {
			case req.Response <- frame:
			default:
				log.Println("wait for matched, but no receivers in channel, probably timed out")
			}

			delete(z.waitForRequests, req)
			close(req.Response)
		}
	}
}

func (z *ZNP) addWaitForRequest(request WaitFrameRequest) {
	z.waitForRequestsMutex.Lock()
	defer z.waitForRequestsMutex.Unlock()

	z.waitForRequests[request] = true
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

var WaitForFrameContextCancelled = errors.New("wait for frame context cancelled")

func (z *ZNP) WaitForFrame(ctx context.Context, messageType unpi.MessageType, subsystem unpi.Subsystem, commandID byte) (unpi.Frame, error) {
	wfr := WaitFrameRequest{
		MessageType: messageType,
		SubSystem:   subsystem,
		CommandID:   commandID,
		Response:    make(chan unpi.Frame),
	}

	z.addWaitForRequest(wfr)

	select {
	case frame := <-wfr.Response:
		return frame, nil
	case <-ctx.Done():
		return unpi.Frame{}, WaitForFrameContextCancelled
	}
}

var SyncRequestContextCancelled = errors.New("synchronous request context cancelled")

func (z *ZNP) SyncRequest(ctx context.Context, frame unpi.Frame) (unpi.Frame, error) {
	if frame.MessageType != unpi.SREQ {
		return unpi.Frame{}, FrameNotSynchronous
	}

	z.syncReceivingMutex.Lock()
	defer z.syncReceivingMutex.Unlock()

	if err := z.writeFrame(frame); err != nil {
		return unpi.Frame{}, err
	}

	select {
	case frame := <-z.syncReceivingChannel:
		return frame, nil
	case <-ctx.Done():
		return unpi.Frame{}, SyncRequestContextCancelled
	}
}

func (z *ZNP) Receive() (unpi.Frame, error) {
	return <-z.asyncReceivingChannel, nil
}
