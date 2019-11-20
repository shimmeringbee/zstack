package zstack

import "context"

func (z *ZStack) writeNVRAM(ctx context.Context, v interface{}) error {
	return nil
}

type SysOSALNVWriteReq struct {
	NVItemID uint16
	Offset   uint8
	Value    []byte `bclength:"8"`
}

const SysOSALNVWriteReqID uint8 = 0x09

type SysOSALNVWriteResp GenericZStackStatus

const SysOSALNVWriteRespID uint8 = 0x09
