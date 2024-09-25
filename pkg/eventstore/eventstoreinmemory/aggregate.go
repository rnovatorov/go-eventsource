package eventstoreinmemory

import (
	"sync"

	"github.com/rnovatorov/go-eventsource/pkg/eventsource"
)

type aggregate struct {
	sync.RWMutex
	version int
	events  eventsource.Events
}
