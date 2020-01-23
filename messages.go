package zstack

import (
	"errors"
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

	l.Add(SREQ, ZDO, ZdoMGMTLQIReqID, ZdoMGMTLQIReq{})
	l.Add(SRSP, ZDO, ZdoMGMTLQIReqRespID, ZdoMGMTLQIReqResp{})
	l.Add(AREQ, ZDO, ZdoMGMTLQIRespID, ZdoMGMTLQIResp{})

	l.Add(SREQ, AF, AFRegisterReqID, AFRegisterReq{})
	l.Add(SRSP, AF, AFRegisterRespID, AFRegisterResp{})
}

type ZStackStatus uint8

type GenericZStackStatus struct {
	Status ZStackStatus
}

var ErrorZFailure = errors.New("ZStack has returned a failure")

const (
	ZSuccess ZStackStatus = 0x00
	ZFailure ZStackStatus = 0x01
)
