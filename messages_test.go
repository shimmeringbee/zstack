package zstack

import (
	. "github.com/shimmeringbee/unpi"
	. "github.com/shimmeringbee/unpi/library"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func Test_registerMessages(t *testing.T) {
	ml := NewLibrary()
	registerMessages(ml)

	t.Run("SysResetReq", func(t *testing.T) {
		identity, found := ml.GetByObject(&SysResetReq{})

		assert.True(t, found)
		assert.Equal(t, AREQ, identity.MessageType)
		assert.Equal(t, SYS, identity.Subsystem)
		assert.Equal(t, uint8(0x00), identity.CommandID)

		ty, found := ml.GetByIdentifier(AREQ, SYS, 0x00)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(SysResetReq{}), ty)
	})

	t.Run("SysResetInd", func(t *testing.T) {
		identity, found := ml.GetByObject(&SysResetInd{})

		assert.True(t, found)
		assert.Equal(t, AREQ, identity.MessageType)
		assert.Equal(t, SYS, identity.Subsystem)
		assert.Equal(t, uint8(0x80), identity.CommandID)

		ty, found := ml.GetByIdentifier(AREQ, SYS, 0x80)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(SysResetInd{}), ty)
	})

	t.Run("SysOSALNVWrite", func(t *testing.T) {
		identity, found := ml.GetByObject(&SysOSALNVWrite{})

		assert.True(t, found)
		assert.Equal(t, SREQ, identity.MessageType)
		assert.Equal(t, SYS, identity.Subsystem)
		assert.Equal(t, uint8(0x09), identity.CommandID)

		ty, found := ml.GetByIdentifier(SREQ, SYS, 0x09)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(SysOSALNVWrite{}), ty)
	})

	t.Run("SysOSALNVWriteReply", func(t *testing.T) {
		identity, found := ml.GetByObject(&SysOSALNVWriteReply{})

		assert.True(t, found)
		assert.Equal(t, SRSP, identity.MessageType)
		assert.Equal(t, SYS, identity.Subsystem)
		assert.Equal(t, uint8(0x09), identity.CommandID)

		ty, found := ml.GetByIdentifier(SRSP, SYS, 0x09)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(SysOSALNVWriteReply{}), ty)
	})

	t.Run("SAPIZBStartRequest", func(t *testing.T) {
		identity, found := ml.GetByObject(&SAPIZBStartRequest{})

		assert.True(t, found)
		assert.Equal(t, SREQ, identity.MessageType)
		assert.Equal(t, SAPI, identity.Subsystem)
		assert.Equal(t, uint8(0x00), identity.CommandID)

		ty, found := ml.GetByIdentifier(SREQ, SAPI, 0x00)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(SAPIZBStartRequest{}), ty)
	})

	t.Run("SAPIZBStartRequestReply", func(t *testing.T) {
		identity, found := ml.GetByObject(&SAPIZBStartRequestReply{})

		assert.True(t, found)
		assert.Equal(t, SRSP, identity.MessageType)
		assert.Equal(t, SAPI, identity.Subsystem)
		assert.Equal(t, uint8(0x00), identity.CommandID)

		ty, found := ml.GetByIdentifier(SRSP, SAPI, 0x00)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(SAPIZBStartRequestReply{}), ty)
	})

	t.Run("SAPIZBPermitJoiningRequest", func(t *testing.T) {
		identity, found := ml.GetByObject(&SAPIZBPermitJoiningRequest{})

		assert.True(t, found)
		assert.Equal(t, SREQ, identity.MessageType)
		assert.Equal(t, SAPI, identity.Subsystem)
		assert.Equal(t, uint8(0x08), identity.CommandID)

		ty, found := ml.GetByIdentifier(SREQ, SAPI, 0x08)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(SAPIZBPermitJoiningRequest{}), ty)
	})

	t.Run("SAPIZBPermitJoiningRequestReply", func(t *testing.T) {
		identity, found := ml.GetByObject(&SAPIZBPermitJoiningRequestReply{})

		assert.True(t, found)
		assert.Equal(t, SRSP, identity.MessageType)
		assert.Equal(t, SAPI, identity.Subsystem)
		assert.Equal(t, uint8(0x08), identity.CommandID)

		ty, found := ml.GetByIdentifier(SRSP, SAPI, 0x08)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(SAPIZBPermitJoiningRequestReply{}), ty)
	})

	t.Run("SAPIZBGetDeviceInfo", func(t *testing.T) {
		identity, found := ml.GetByObject(&SAPIZBGetDeviceInfo{})

		assert.True(t, found)
		assert.Equal(t, SREQ, identity.MessageType)
		assert.Equal(t, SAPI, identity.Subsystem)
		assert.Equal(t, uint8(0x06), identity.CommandID)

		ty, found := ml.GetByIdentifier(SREQ, SAPI, 0x06)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(SAPIZBGetDeviceInfo{}), ty)
	})

	t.Run("SAPIZBGetDeviceInfoReply", func(t *testing.T) {
		identity, found := ml.GetByObject(&SAPIZBGetDeviceInfoReply{})

		assert.True(t, found)
		assert.Equal(t, SRSP, identity.MessageType)
		assert.Equal(t, SAPI, identity.Subsystem)
		assert.Equal(t, uint8(0x06), identity.CommandID)

		ty, found := ml.GetByIdentifier(SRSP, SAPI, 0x06)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(SAPIZBGetDeviceInfoReply{}), ty)
	})

	t.Run("ZDOStateChangeInd", func(t *testing.T) {
		identity, found := ml.GetByObject(&ZDOStateChangeInd{})

		assert.True(t, found)
		assert.Equal(t, AREQ, identity.MessageType)
		assert.Equal(t, ZDO, identity.Subsystem)
		assert.Equal(t, uint8(0xc0), identity.CommandID)

		ty, found := ml.GetByIdentifier(AREQ, ZDO, 0xc0)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(ZDOStateChangeInd{}), ty)
	})

	t.Run("ZdoEndDeviceAnnceInd", func(t *testing.T) {
		identity, found := ml.GetByObject(&ZdoEndDeviceAnnceInd{})

		assert.True(t, found)
		assert.Equal(t, AREQ, identity.MessageType)
		assert.Equal(t, ZDO, identity.Subsystem)
		assert.Equal(t, uint8(0xc1), identity.CommandID)

		ty, found := ml.GetByIdentifier(AREQ, ZDO, 0xc1)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(ZdoEndDeviceAnnceInd{}), ty)
	})

	t.Run("ZdoLeaveInd", func(t *testing.T) {
		identity, found := ml.GetByObject(&ZdoLeaveInd{})

		assert.True(t, found)
		assert.Equal(t, AREQ, identity.MessageType)
		assert.Equal(t, ZDO, identity.Subsystem)
		assert.Equal(t, uint8(0xc9), identity.CommandID)

		ty, found := ml.GetByIdentifier(AREQ, ZDO, 0xc9)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(ZdoLeaveInd{}), ty)
	})

	t.Run("ZdoTcDevInd", func(t *testing.T) {
		identity, found := ml.GetByObject(&ZdoTcDevInd{})

		assert.True(t, found)
		assert.Equal(t, AREQ, identity.MessageType)
		assert.Equal(t, ZDO, identity.Subsystem)
		assert.Equal(t, uint8(0xca), identity.CommandID)

		ty, found := ml.GetByIdentifier(AREQ, ZDO, 0xca)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(ZdoTcDevInd{}), ty)
	})

	t.Run("ZdoMGMTLQIReq", func(t *testing.T) {
		identity, found := ml.GetByObject(&ZdoMGMTLQIReq{})

		assert.True(t, found)
		assert.Equal(t, SREQ, identity.MessageType)
		assert.Equal(t, ZDO, identity.Subsystem)
		assert.Equal(t, uint8(0x31), identity.CommandID)

		ty, found := ml.GetByIdentifier(SREQ, ZDO, 0x31)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(ZdoMGMTLQIReq{}), ty)
	})

	t.Run("ZdoMGMTLQIReqReply", func(t *testing.T) {
		identity, found := ml.GetByObject(&ZdoMGMTLQIReqReply{})

		assert.True(t, found)
		assert.Equal(t, SRSP, identity.MessageType)
		assert.Equal(t, ZDO, identity.Subsystem)
		assert.Equal(t, uint8(0x31), identity.CommandID)

		ty, found := ml.GetByIdentifier(SRSP, ZDO, 0x31)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(ZdoMGMTLQIReqReply{}), ty)
	})

	t.Run("ZdoMGMTLQIRsp", func(t *testing.T) {
		identity, found := ml.GetByObject(&ZdoMGMTLQIRsp{})

		assert.True(t, found)
		assert.Equal(t, AREQ, identity.MessageType)
		assert.Equal(t, ZDO, identity.Subsystem)
		assert.Equal(t, uint8(0xb1), identity.CommandID)

		ty, found := ml.GetByIdentifier(AREQ, ZDO, 0xb1)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(ZdoMGMTLQIRsp{}), ty)
	})

	t.Run("AFRegister", func(t *testing.T) {
		identity, found := ml.GetByObject(&AFRegister{})

		assert.True(t, found)
		assert.Equal(t, SREQ, identity.MessageType)
		assert.Equal(t, AF, identity.Subsystem)
		assert.Equal(t, uint8(0x00), identity.CommandID)

		ty, found := ml.GetByIdentifier(SREQ, AF, 0x00)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(AFRegister{}), ty)
	})

	t.Run("AFRegisterReply", func(t *testing.T) {
		identity, found := ml.GetByObject(&AFRegisterReply{})

		assert.True(t, found)
		assert.Equal(t, SRSP, identity.MessageType)
		assert.Equal(t, AF, identity.Subsystem)
		assert.Equal(t, uint8(0x00), identity.CommandID)

		ty, found := ml.GetByIdentifier(SRSP, AF, 0x00)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(AFRegisterReply{}), ty)
	})

	t.Run("ZdoActiveEpReq", func(t *testing.T) {
		identity, found := ml.GetByObject(&ZdoActiveEpReq{})

		assert.True(t, found)
		assert.Equal(t, SREQ, identity.MessageType)
		assert.Equal(t, ZDO, identity.Subsystem)
		assert.Equal(t, uint8(0x05), identity.CommandID)

		ty, found := ml.GetByIdentifier(SREQ, ZDO, 0x05)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(ZdoActiveEpReq{}), ty)
	})

	t.Run("ZdoActiveEpReqReply", func(t *testing.T) {
		identity, found := ml.GetByObject(&ZdoActiveEpReqReply{})

		assert.True(t, found)
		assert.Equal(t, SRSP, identity.MessageType)
		assert.Equal(t, ZDO, identity.Subsystem)
		assert.Equal(t, uint8(0x05), identity.CommandID)

		ty, found := ml.GetByIdentifier(SRSP, ZDO, 0x05)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(ZdoActiveEpReqReply{}), ty)
	})

	t.Run("ZdoActiveEpRsp", func(t *testing.T) {
		identity, found := ml.GetByObject(&ZdoActiveEpRsp{})

		assert.True(t, found)
		assert.Equal(t, AREQ, identity.MessageType)
		assert.Equal(t, ZDO, identity.Subsystem)
		assert.Equal(t, uint8(0x85), identity.CommandID)

		ty, found := ml.GetByIdentifier(AREQ, ZDO, 0x85)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(ZdoActiveEpRsp{}), ty)
	})
}
