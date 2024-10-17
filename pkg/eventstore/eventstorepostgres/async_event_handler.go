package eventstorepostgres

import (
	"context"

	"github.com/rnovatorov/go-eventsource/pkg/eventsource"
)

type AsyncEventHandler interface {
	SubscriptionID() string
	HandleEvent(context.Context, *eventsource.Event) error
}
