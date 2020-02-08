package zstack

import (
	"errors"
	. "github.com/shimmeringbee/unpi"
	. "github.com/shimmeringbee/unpi/library"
)

func registerMessages(l *Library) {
	l.Add(AREQ, SYS, SysResetReqID, SysResetReq{})
	l.Add(AREQ, SYS, SysResetIndID, SysResetInd{})

	l.Add(SREQ, SYS, SysOSALNVWriteID, SysOSALNVWrite{})
	l.Add(SRSP, SYS, SysOSALNVWriteReplyID, SysOSALNVWriteReply{})

	l.Add(SREQ, SAPI, SAPIZBStartRequestID, SAPIZBStartRequest{})
	l.Add(SRSP, SAPI, SAPIZBStartRequestReplyID, SAPIZBStartRequestReply{})

	l.Add(SREQ, SAPI, SAPIZBPermitJoiningRequestID, SAPIZBPermitJoiningRequest{})
	l.Add(SRSP, SAPI, SAPIZBPermitJoiningRequestReplyID, SAPIZBPermitJoiningRequestReply{})

	l.Add(SREQ, SAPI, SAPIZBGetDeviceInfoID, SAPIZBGetDeviceInfo{})
	l.Add(SRSP, SAPI, SAPIZBGetDeviceInfoReplyID, SAPIZBGetDeviceInfoReply{})

	l.Add(AREQ, ZDO, ZDOStateChangeIndID, ZDOStateChangeInd{})

	l.Add(AREQ, ZDO, ZdoEndDeviceAnnceIndID, ZdoEndDeviceAnnceInd{})
	l.Add(AREQ, ZDO, ZdoLeaveIndID, ZdoLeaveInd{})
	l.Add(AREQ, ZDO, ZdoTcDevIndID, ZdoTcDevInd{})

	l.Add(SREQ, ZDO, ZdoMGMTLQIReqID, ZdoMGMTLQIReq{})
	l.Add(SRSP, ZDO, ZdoMGMTLQIReqReplyID, ZdoMGMTLQIReqReply{})
	l.Add(AREQ, ZDO, ZdoMGMTLQIRspID, ZdoMGMTLQIRsp{})

	l.Add(SREQ, AF, AFRegisterID, AFRegister{})
	l.Add(SRSP, AF, AFRegisterReplyID, AFRegisterReply{})

	l.Add(SREQ, ZDO, ZdoActiveEpReqID, ZdoActiveEpReq{})
	l.Add(SRSP, ZDO, ZdoActiveEpReqReplyID, ZdoActiveEpReqReply{})
	l.Add(AREQ, ZDO, ZdoActiveEpRspID, ZdoActiveEpRsp{})

	l.Add(SREQ, ZDO, ZdoSimpleDescReqID, ZdoSimpleDescReq{})
	l.Add(SRSP, ZDO, ZdoSimpleDescReqReplyID, ZdoSimpleDescReqReply{})
	l.Add(AREQ, ZDO, ZdoSimpleDescRspID, ZdoSimpleDescRsp{})
}

type ZStackStatus uint8

type HasSuccess interface {
	WasSuccessful() bool
}

type GenericZStackStatus struct {
	Status ZStackStatus
}

func (s GenericZStackStatus) WasSuccessful() bool {
	return s.Status == ZSuccess
}

var ErrorZFailure = errors.New("ZStack has returned a failure")

const (
	ZSuccess ZStackStatus = 0x00
	ZFailure ZStackStatus = 0x01
)
