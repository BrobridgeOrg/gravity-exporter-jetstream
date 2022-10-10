package eventbus

import (
	"time"

	"github.com/BrobridgeOrg/gravity-exporter-jetstream/pkg/app"
	nats "github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

type Options struct {
	PingInterval        time.Duration
	MaxPingsOutstanding int
	MaxReconnects       int
}

type EventBusHandler struct {
	Reconnect  func(natsConn *nats.Conn)
	Disconnect func(natsConn *nats.Conn)
}

type EventBus struct {
	app           app.App
	connection    *nats.Conn
	jetStreamConn nats.JetStreamContext
	host          string
	handler       *EventBusHandler
	options       *Options
}

func NewEventBus(a app.App, host string, handler EventBusHandler, options Options) *EventBus {
	return &EventBus{
		app:        a,
		connection: nil,
		host:       host,
		handler:    &handler,
		options:    &options,
	}
}

func (eb *EventBus) Connect() error {

	log.WithFields(log.Fields{
		"host":                eb.host,
		"PingInterval":        eb.options.PingInterval * time.Second,
		"MaxPingsOutnatsding": eb.options.MaxPingsOutstanding,
		"MaxReconnects":       eb.options.MaxReconnects,
	}).Info("Connecting to NATS server")

	nc, err := nats.Connect(eb.host,
		nats.PingInterval(eb.options.PingInterval*time.Second),
		nats.MaxPingsOutstanding(eb.options.MaxPingsOutstanding),
		nats.MaxReconnects(eb.options.MaxReconnects),
		nats.ReconnectHandler(eb.ReconnectHandler),
		nats.DisconnectHandler(eb.handler.Disconnect),
	)
	if err != nil {
		return err
	}

	eb.connection = nc

	// Connect to NATS JetStream
	err = eb.ConnectToJetStream()
	if err != nil {
		return err
	}

	return nil
}

func (eb *EventBus) ConnectToJetStream() error {

	log.WithFields(log.Fields{}).Info("Connecting to NATS JetStream")

	js, err := eb.connection.JetStream(nats.PublishAsyncMaxPending(102400))
	if err != nil {
		return err
	}

	eb.jetStreamConn = js

	return nil
}

func (eb *EventBus) Close() {
	eb.connection.Close()
}
func (eb *EventBus) ReconnectHandler(natsConn *nats.Conn) {

	// Reconnect to NATS JetStreamg
	err := eb.ConnectToJetStream()
	if err != nil {
		log.Error(err)
		return
	}

	eb.handler.Reconnect(natsConn)
}

func (eb *EventBus) GetConnection() *nats.Conn {
	return eb.connection
}

func (eb *EventBus) GetJetStreamConnection() nats.JetStreamContext {
	return eb.jetStreamConn
}
