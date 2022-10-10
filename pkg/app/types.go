package app

import (
	"github.com/BrobridgeOrg/gravity-exporter-jetstream/pkg/eventbus"
)

type App interface {
	GetEventBus() eventbus.EventBus
}
