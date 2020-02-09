package zstack // import "github.com/shimmeringbee/zstack"

import (
	"context"
	"github.com/shimmeringbee/unpi/broker"
	"github.com/shimmeringbee/unpi/library"
	"github.com/shimmeringbee/zigbee"
	"io"
	"time"
)

type RequestResponder interface {
	RequestResponse(ctx context.Context, req interface{}, resp interface{}) error
}

type Awaiter interface {
	Await(ctx context.Context, resp interface{}) error
}

type Subscriber interface {
	Subscribe(message interface{}, callback func(v interface{})) (error, func())
}

type ZStack struct {
	requestResponder RequestResponder
	awaiter          Awaiter
	subscriber       Subscriber

	NetworkProperties NetworkProperties

	events chan interface{}

	networkManagerStop     chan bool
	networkManagerIncoming chan interface{}

	messageReceiverStop func()

	devices          map[zigbee.IEEEAddress]*Device
	devicesByNetAddr map[zigbee.NetworkAddress]zigbee.IEEEAddress
}

type JoinState uint8

const (
	Off           JoinState = 0x00
	OnCoordinator JoinState = 0x01
	OnAllRouters  JoinState = 0x02
)

type NetworkProperties struct {
	NetworkAddress zigbee.NetworkAddress
	IEEEAddress    zigbee.IEEEAddress
	PANID          zigbee.PANID
	ExtendedPANID  zigbee.ExtendedPANID
	NetworkKey     zigbee.NetworkKey
	Channel        uint8
	JoinState      JoinState
}

const DefaultZStackTimeout = 5 * time.Second
const DefaultZStackRetries = 3
const DefaultInflightEvents = 50

func New(uart io.ReadWriter) *ZStack {
	ml := library.NewLibrary()
	registerMessages(ml)

	znp := broker.NewBroker(uart, uart, ml)

	return &ZStack{
		requestResponder:       znp,
		awaiter:                znp,
		subscriber:             znp,
		events:                 make(chan interface{}, DefaultInflightEvents),
		networkManagerStop:     make(chan bool, 1),
		networkManagerIncoming: make(chan interface{}, DefaultInflightEvents),
		devices:                map[zigbee.IEEEAddress]*Device{},
		devicesByNetAddr:       map[zigbee.NetworkAddress]zigbee.IEEEAddress{},
	}
}

func (z *ZStack) Stop() {
	z.stopNetworkManager()
	z.stopMessageReceiver()
}
