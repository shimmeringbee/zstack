package zstack

import (
	"context"
	"errors"
	"fmt"
	"github.com/shimmeringbee/bytecodec"
	"github.com/shimmeringbee/zigbee"
	"reflect"
)

var NVRAMUnsuccessful = errors.New("nvram operation unsuccessful")
var NVRAMUnrecognised = errors.New("nvram structure unrecognised")

func (z *ZStack) writeNVRAM(ctx context.Context, v interface{}) error {
	configId, found := nvramStructToID[reflect.TypeOf(v)]

	if !found {
		return NVRAMUnrecognised
	}

	configValue, err := bytecodec.Marshal(v)
	if err != nil {
		return err
	}

	writeRequest := SysOSALNVWrite{
		NVItemID: configId,
		Offset:   0,
		Value:    configValue,
	}

	writeResponse := SysOSALNVWriteReply{}

	if err := z.requestResponder.RequestResponse(ctx, writeRequest, &writeResponse); err != nil {
		return err
	}

	if writeResponse.Status != ZSuccess {
		return fmt.Errorf("%w: status = %v", NVRAMUnsuccessful, writeResponse.Status)
	}

	return nil
}

func (z *ZStack) readNVRAM(ctx context.Context, v interface{}) error {
	vType := reflect.TypeOf(v)

	if vType.Kind() == reflect.Ptr {
		vType = vType.Elem()
	}

	configId, found := nvramStructToID[vType]

	if !found {
		return NVRAMUnrecognised
	}

	readRequest := SysOSALNVRead{
		NVItemID: configId,
		Offset:   0,
	}

	readResponse := SysOSALNVReadReply{}

	if err := z.requestResponder.RequestResponse(ctx, readRequest, &readResponse); err != nil {
		return err
	}

	if readResponse.Status != ZSuccess {
		return fmt.Errorf("%w: status = %v", NVRAMUnsuccessful, readResponse.Status)
	}

	return bytecodec.Unmarshal(readResponse.Value, v)
}

type SysOSALNVWrite struct {
	NVItemID uint16
	Offset   uint8
	Value    []byte `bcsliceprefix:"8"`
}

const SysOSALNVWriteID uint8 = 0x09

type SysOSALNVWriteReply GenericZStackStatus

const SysOSALNVWriteReplyID uint8 = 0x09

type SysOSALNVRead struct {
	NVItemID uint16
	Offset   uint8
}

const SysOSALNVReadID uint8 = 0x08

type SysOSALNVReadReply struct {
	Status ZStackStatus
	Value  []byte `bcsliceprefix:"8"`
}

const SysOSALNVReadReplyID uint8 = 0x08

var nvramStructToID = map[reflect.Type]uint16{
	reflect.TypeOf(ZCDNVStartUpOption{}):    ZCDNVStartUpOptionID,
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

const ZCDNVStartUpOptionID uint16 = 0x0003

type ZCDNVStartUpOption struct {
	StartOption uint8
}

const ZCDNVLogicalTypeID uint16 = 0x0087

type ZCDNVLogicalType struct {
	LogicalType zigbee.LogicalType
}

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
