package eventstore

import (
	"context"
)

type EventHandler func(context.Context, *Event) error
