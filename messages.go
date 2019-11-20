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

type GenericZStackStatus struct {
	Status uint8
}
