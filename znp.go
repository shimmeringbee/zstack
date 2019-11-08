package zstack // import "github.com/shimmeringbee/zstack"

import (
	"github.com/shimmeringbee/unpi"
	"io"
)

type ReadFrame func(r io.Reader) (*unpi.Frame, error)
type WriteFrame func(w io.Writer, frame *unpi.Frame) error
