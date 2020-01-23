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

	t.Run("SysOSALNVWriteReq", func(t *testing.T) {
		identity, found := ml.GetByObject(&SysOSALNVWriteReq{})

		assert.True(t, found)
		assert.Equal(t, SREQ, identity.MessageType)
		assert.Equal(t, SYS, identity.Subsystem)
		assert.Equal(t, uint8(0x09), identity.CommandID)

		ty, found := ml.GetByIdentifier(SREQ, SYS, 0x09)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(SysOSALNVWriteReq{}), ty)
	})

	t.Run("SysOSALNVWriteResp", func(t *testing.T) {
		identity, found := ml.GetByObject(&SysOSALNVWriteResp{})

		assert.True(t, found)
		assert.Equal(t, SRSP, identity.MessageType)
		assert.Equal(t, SYS, identity.Subsystem)
		assert.Equal(t, uint8(0x09), identity.CommandID)

		ty, found := ml.GetByIdentifier(SRSP, SYS, 0x09)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(SysOSALNVWriteResp{}), ty)
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

	t.Run("SAPIZBStartResponse", func(t *testing.T) {
		identity, found := ml.GetByObject(&SAPIZBStartResponse{})

		assert.True(t, found)
		assert.Equal(t, SRSP, identity.MessageType)
		assert.Equal(t, SAPI, identity.Subsystem)
		assert.Equal(t, uint8(0x00), identity.CommandID)

		ty, found := ml.GetByIdentifier(SRSP, SAPI, 0x00)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(SAPIZBStartResponse{}), ty)
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

	t.Run("SAPIZBPermitJoiningResponse", func(t *testing.T) {
		identity, found := ml.GetByObject(&SAPIZBPermitJoiningResponse{})

		assert.True(t, found)
		assert.Equal(t, SRSP, identity.MessageType)
		assert.Equal(t, SAPI, identity.Subsystem)
		assert.Equal(t, uint8(0x08), identity.CommandID)

		ty, found := ml.GetByIdentifier(SRSP, SAPI, 0x08)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(SAPIZBPermitJoiningResponse{}), ty)
	})

	t.Run("SAPIZBGetDeviceInfoReq", func(t *testing.T) {
		identity, found := ml.GetByObject(&SAPIZBGetDeviceInfoReq{})

		assert.True(t, found)
		assert.Equal(t, SREQ, identity.MessageType)
		assert.Equal(t, SAPI, identity.Subsystem)
		assert.Equal(t, uint8(0x06), identity.CommandID)

		ty, found := ml.GetByIdentifier(SREQ, SAPI, 0x06)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(SAPIZBGetDeviceInfoReq{}), ty)
	})

	t.Run("SAPIZBGetDeviceInfoResp", func(t *testing.T) {
		identity, found := ml.GetByObject(&SAPIZBGetDeviceInfoResp{})

		assert.True(t, found)
		assert.Equal(t, SRSP, identity.MessageType)
		assert.Equal(t, SAPI, identity.Subsystem)
		assert.Equal(t, uint8(0x06), identity.CommandID)

		ty, found := ml.GetByIdentifier(SRSP, SAPI, 0x06)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(SAPIZBGetDeviceInfoResp{}), ty)
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

	t.Run("ZdoMGMTLQIReqResp", func(t *testing.T) {
		identity, found := ml.GetByObject(&ZdoMGMTLQIReqResp{})

		assert.True(t, found)
		assert.Equal(t, SRSP, identity.MessageType)
		assert.Equal(t, ZDO, identity.Subsystem)
		assert.Equal(t, uint8(0x31), identity.CommandID)

		ty, found := ml.GetByIdentifier(SRSP, ZDO, 0x31)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(ZdoMGMTLQIReqResp{}), ty)
	})

	t.Run("ZdoMGMTLQIResp", func(t *testing.T) {
		identity, found := ml.GetByObject(&ZdoMGMTLQIResp{})

		assert.True(t, found)
		assert.Equal(t, AREQ, identity.MessageType)
		assert.Equal(t, ZDO, identity.Subsystem)
		assert.Equal(t, uint8(0xb1), identity.CommandID)

		ty, found := ml.GetByIdentifier(AREQ, ZDO, 0xb1)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(ZdoMGMTLQIResp{}), ty)
	})

	t.Run("AFRegisterReq", func(t *testing.T) {
		identity, found := ml.GetByObject(&AFRegisterReq{})

		assert.True(t, found)
		assert.Equal(t, SREQ, identity.MessageType)
		assert.Equal(t, AF, identity.Subsystem)
		assert.Equal(t, uint8(0x00), identity.CommandID)

		ty, found := ml.GetByIdentifier(SREQ, AF, 0x00)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(AFRegisterReq{}), ty)
	})

	t.Run("AFRegisterResp", func(t *testing.T) {
		identity, found := ml.GetByObject(&AFRegisterResp{})

		assert.True(t, found)
		assert.Equal(t, SRSP, identity.MessageType)
		assert.Equal(t, AF, identity.Subsystem)
		assert.Equal(t, uint8(0x00), identity.CommandID)

		ty, found := ml.GetByIdentifier(SRSP, AF, 0x00)

		assert.True(t, found)
		assert.Equal(t, reflect.TypeOf(AFRegisterResp{}), ty)
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
