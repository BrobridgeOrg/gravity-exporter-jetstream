package instance

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	eventbus "github.com/BrobridgeOrg/gravity-exporter-jetstream/pkg/eventbus/service"
	subscriber "github.com/BrobridgeOrg/gravity-exporter-jetstream/pkg/subscriber/service"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type AppInstance struct {
	done       chan os.Signal
	eventBus   *eventbus.EventBus
	subscriber *subscriber.Subscriber
}

func NewAppInstance() *AppInstance {

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill, syscall.SIGTERM)
	a := &AppInstance{
		done: sig,
	}

	return a
}

func (a *AppInstance) Init() error {

	log.WithFields(log.Fields{
		"max_procs": runtime.GOMAXPROCS(0),
	}).Info("Starting application")

	// Initializing modules
	a.eventBus = eventbus.NewEventBus(
		a,
		viper.GetString("jetstream.host"),
		eventbus.EventBusHandler{
			Reconnect: func(natsConn *nats.Conn) {
				log.Warn("re-connected to event server")
			},
			Disconnect: func(natsConn *nats.Conn) {
				log.Error("event server was disconnected")
			},
		},
		eventbus.Options{
			PingInterval:        time.Duration(viper.GetInt64("jetstream.pingInterval")),
			MaxPingsOutstanding: viper.GetInt("jetstream.maxPingsOutstanding"),
			MaxReconnects:       viper.GetInt("jetstream.maxReconnects"),
		},
	)

	a.subscriber = subscriber.NewSubscriber(a)

	// Initializing EventBus
	err := a.initEventBus()
	if err != nil {
		return err
	}

	err = a.subscriber.Init()
	if err != nil {
		return err
	}

	return nil
}

func (a *AppInstance) Uninit() {
}

func (a *AppInstance) Run() error {

	err := a.subscriber.Run()
	if err != nil {
		return err
	}

	<-a.done
	a.subscriber.Stop()
	time.Sleep(5 * time.Second)
	log.Error("Bye!")

	return nil
}
