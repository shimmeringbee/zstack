package zstack

import (
	. "github.com/shimmeringbee/unpi"
	. "github.com/shimmeringbee/unpi/library"
)

func registerMessages(library *Library) {
	library.Add(AREQ, SYS, SysResetReqID, SysResetReq{})
	library.Add(AREQ, SYS, SysResetIndID, SysResetInd{})

	library.Add(SREQ, SYS, SysOSALNVWriteReqID, SysOSALNVWriteReq{})
	library.Add(SRSP, SYS, SysOSALNVWriteRespID, SysOSALNVWriteResp{})
}

type ZStackStatus uint8

type GenericZStackStatus struct {
	Status ZStackStatus
}

const (
	ZSuccess ZStackStatus = 0x00
)
