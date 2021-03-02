package zstack // import "github.com/shimmeringbee/zstack"

import (
	"context"
	"github.com/shimmeringbee/logwrap"
	"github.com/shimmeringbee/logwrap/impl/golog"
	"github.com/shimmeringbee/unpi/broker"
	"github.com/shimmeringbee/unpi/library"
	"github.com/shimmeringbee/zigbee"
	"io"
	"log"
	"os"
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

	nodeTable          *NodeTable
	transactionIdStore chan uint8

	logger logwrap.Logger
}

var _ zigbee.Provider = (*ZStack)(nil)

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
const DefaultResolveIEEETimeout = 500 * time.Millisecond
const DefaultZStackRetries = 3
const DefaultInflightEvents = 50
const DefaultInflightTransactions = 20

func New(uart io.ReadWriter, nodeTable *NodeTable) *ZStack {
	ml := library.NewLibrary()
	registerMessages(ml)

	znp := broker.NewBroker(uart, uart, ml)

	transactionIDs := make(chan uint8, DefaultInflightTransactions)

	for i := 0; i < DefaultInflightTransactions; i++ {
		transactionIDs <- uint8(i)
	}

	zstack := &ZStack{
		requestResponder:       znp,
		awaiter:                znp,
		subscriber:             znp,
		events:                 make(chan interface{}, DefaultInflightEvents),
		networkManagerStop:     make(chan bool, 1),
		networkManagerIncoming: make(chan interface{}, DefaultInflightEvents),
		nodeTable:              nodeTable,
		transactionIdStore:     transactionIDs,
	}

	zstack.WithGoLogger(log.New(os.Stderr, "", log.LstdFlags))

	return zstack
}

func (z *ZStack) Stop() {
	z.stopNetworkManager()
	z.stopMessageReceiver()
}

func (z *ZStack) WithGoLogger(parentLogger *log.Logger) {
	z.logger = logwrap.New(golog.Wrap(parentLogger))
}

func (z *ZStack) WithLogWrapLogger(parentLogger logwrap.Logger) {
	z.logger = parentLogger
}
