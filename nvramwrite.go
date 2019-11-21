package zstack

import (
	"context"
	"errors"
	"fmt"
	"github.com/shimmeringbee/bytecodec"
	"github.com/shimmeringbee/zigbee"
	"reflect"
)

var NVRAMWriteUnsuccessful = errors.New("nvram write unsuccessful")
var NVRAMUnrecognised = errors.New("nvram write structure unrecognised")

func (z *ZStack) writeNVRAM(ctx context.Context, v interface{}) error {
	configId, found := nvMapping[reflect.TypeOf(v)]

	if !found {
		return NVRAMUnrecognised
	}

	configValue, err := bytecodec.Marshall(v)
	if err != nil {
		return err
	}

	writeRequest := SysOSALNVWriteReq{
		NVItemID: configId,
		Offset:   0,
		Value:    configValue,
	}

	writeResponse := SysOSALNVWriteResp{}

	if err := z.RequestResponder.MessageRequestResponse(ctx, writeRequest, &writeResponse); err != nil {
		return err
	}

	if writeResponse.Status != ZSuccess {
		return fmt.Errorf("%w: status = %v", NVRAMWriteUnsuccessful, writeResponse.Status)
	}

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

var nvMapping = map[reflect.Type]uint16{
	reflect.TypeOf(NCDNVStartUpOption{}):    NCDNVStartUpOptionID,
	reflect.TypeOf(ZCDNVLogicalType{}):      ZCDNVLogicalTypeID,
	reflect.TypeOf(ZCDNVSecurityMode{}):     ZCDNVSecurityModeID,
	reflect.TypeOf(ZCDNVPreCfgKeysEnable{}): ZCDNVPreCfgKeysEnableID,
	reflect.TypeOf(ZCDNVPreCfgKey{}):        ZCDNVPreCfgKeyID,
	reflect.TypeOf(ZCDNVZDODirectCB{}):      ZCDNVZDODirectCBID,
	reflect.TypeOf(ZCDNVChanList{}):         ZCDNVChanListID,
	reflect.TypeOf(ZCDNVPANID{}):            ZCDNVPANIDID,
	reflect.TypeOf(ZCDNVExtPANID{}):         ZCDNVExtPANIDID,
	reflect.TypeOf(ZCDNVUseDefaultTCLK{}):   ZCDNVUseDefaultTCLKID,
	reflect.TypeOf(ZCDNVTCLKTableStart{}):   ZCDNVTCLKTableStartID,
}

const NCDNVStartUpOptionID uint16 = 0x0003

type NCDNVStartUpOption struct {
	StartOption uint8
}

const ZCDNVLogicalTypeID uint16 = 0x0087

type ZCDNVLogicalType struct {
	LogicalType LogicalType
}

type LogicalType uint8

const (
	Coordinator LogicalType = 0x00
	Router      LogicalType = 0x01
	EndDevice   LogicalType = 0x02
)

const ZCDNVSecurityModeID uint16 = 0x0064

type ZCDNVSecurityMode struct {
	Enabled uint8
}

const ZCDNVPreCfgKeysEnableID uint16 = 0x0063

type ZCDNVPreCfgKeysEnable struct {
	Enabled uint8
}

const ZCDNVPreCfgKeyID uint16 = 0x0062

type ZCDNVPreCfgKey struct {
	NetworkKey zigbee.NetworkKey
}

const ZCDNVZDODirectCBID uint16 = 0x008f

type ZCDNVZDODirectCB struct {
	Enabled uint8
}

const ZCDNVChanListID uint16 = 0x0084

type ZCDNVChanList struct {
	Channels [4]byte
}

const ZCDNVPANIDID uint16 = 0x0083

type ZCDNVPANID struct {
	PANID zigbee.PANID
}

const ZCDNVExtPANIDID uint16 = 0x002d

type ZCDNVExtPANID struct {
	ExtendedPANID zigbee.ExtendedPANID
}

const ZCDNVUseDefaultTCLKID uint16 = 0x006d

type ZCDNVUseDefaultTCLK struct {
	Enabled uint8
}

const ZCDNVTCLKTableStartID uint16 = 0x0101

type ZCDNVTCLKTableStart struct {
	Address        zigbee.IEEEAddress
	NetworkKey     zigbee.NetworkKey
	TXFrameCounter uint32
	RXFrameCounter uint32
}
