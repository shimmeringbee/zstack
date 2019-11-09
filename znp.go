package zstack // import "github.com/shimmeringbee/zstack"

import (
	"github.com/shimmeringbee/unpi"
	"io"
)

type ZNP struct {
	reader io.Reader
	writer io.Writer
}

func New(device io.ReadWriter) *ZNP {
	z := &ZNP{
		reader: device,
		writer: device,
	}

	z.start()

	return z
}

func (z *ZNP) start() {

}

func (z *ZNP) Stop() {

}

func (z *ZNP) AsyncRequest(frame unpi.Frame) error {
	return unpi.Write(z.writer, &frame)
}

func (z *ZNP) SyncRequest(frame unpi.Frame) (unpi.Frame, error) {
	return unpi.Frame{}, nil
}

func (z *ZNP) Receive() (unpi.Frame, error) {
	return unpi.Frame{}, nil
}
