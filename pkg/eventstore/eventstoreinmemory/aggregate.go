package eventstoreinmemory

import (
	"sync"

	"github.com/rnovatorov/go-eventsource/pkg/eventstore"
)

type aggregate struct {
	sync.RWMutex
	version int
	events  eventstore.Events
}
