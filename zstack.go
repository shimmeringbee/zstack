package zstack // import "github.com/shimmeringbee/zstack"

import (
	"github.com/shimmeringbee/unpi"
	"github.com/shimmeringbee/unpi/library"
)

func PopulateMessageLibrary() *library.Library {
	l := library.NewLibrary()

	l.Add(unpi.AREQ, unpi.SYS, SysResetRequestID, SysResetReq{})
	l.Add(unpi.AREQ, unpi.SYS, SysResetIndicationCommandID, SysResetInd{})

	l.Add(unpi.SREQ, unpi.SYS, SysOSALNVWriteRequestID, SysOSALNVWrite{})
	l.Add(unpi.SRSP, unpi.SYS, SysOSALNVWriteResponseID, SysOSALNVWriteResponse{})

	return l
}

type ResetType uint8

const (
	Hard ResetType = 0
	Soft ResetType = 1
)

type SysResetReq struct {
	ResetType ResetType
}

const SysResetRequestID uint8 = 0x00

type ResetReason uint8

const (
	PowerUp  ResetReason = 0
	External ResetReason = 1
	Watchdog ResetReason = 2
)

type SysResetInd struct {
	Reason            ResetReason
	TransportRevision uint8
	ProductID         uint8
	MajorRelease      uint8
	MinorRelease      uint8
	HardwareRevision  uint8
}

const SysResetIndicationCommandID uint8 = 0x80

type SysOSALNVWrite struct {
	NVItemID uint16
	Offset   uint8
	Value    []byte `bclength:"8"`
}

const SysOSALNVWriteRequestID uint8 = 0x09

type SysOSALNVWriteResponse struct {
	Status uint8
}

const SysOSALNVWriteResponseID uint8 = 0x09
