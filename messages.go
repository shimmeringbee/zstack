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

	l.Add(SREQ, SAPI, SAPIZBPermitJoiningRequestID, SAPIZBPermitJoiningRequest{})
	l.Add(SRSP, SAPI, SAPIZBPermitJoiningResponseID, SAPIZBPermitJoiningResponse{})

	l.Add(SREQ, SAPI, SAPIZBGetDeviceInfoReqID, SAPIZBGetDeviceInfoReq{})
	l.Add(SRSP, SAPI, SAPIZBGetDeviceInfoRespID, SAPIZBGetDeviceInfoResp{})

	l.Add(AREQ, ZDO, ZDOStateChangeIndID, ZDOStateChangeInd{})

	l.Add(AREQ, ZDO, ZdoEndDeviceAnnceIndID, ZdoEndDeviceAnnceInd{})
	l.Add(AREQ, ZDO, ZdoLeaveIndID, ZdoLeaveInd{})
	l.Add(AREQ, ZDO, ZdoTcDevIndID, ZdoTcDevInd{})
}

type ZStackStatus uint8

type GenericZStackStatus struct {
	Status ZStackStatus
}

const (
	ZSuccess ZStackStatus = 0x00
)
