package zstack

import (
	. "github.com/shimmeringbee/unpi"
	. "github.com/shimmeringbee/unpi/library"
)

func registerMessages(l *Library) {
	l.Add(AREQ, SYS, SysResetReqID, SysResetReq{})
	l.Add(AREQ, SYS, SysResetIndID, SysResetInd{})

	l.Add(SREQ, SYS, SysOSALNVWriteReqID, SysOSALNVWriteReq{})
	l.Add(SRSP, SYS, SysOSALNVWriteRespID, SysOSALNVWriteResp{})

	l.Add(SREQ, SAPI, SAPIZBStartRequestID, SAPIZBStartRequest{})
	l.Add(SRSP, SAPI, SAPIZBStartResponseID, SAPIZBStartResponse{})

	l.Add(AREQ, SAPI, SAPIZBStartConfirmID, SAPIZBStartConfirm{})
}

type ZStackStatus uint8

type GenericZStackStatus struct {
	Status ZStackStatus
}

const (
	ZSuccess ZStackStatus = 0x00
)
