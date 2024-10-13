package eventstorepostgres

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/rnovatorov/go-eventsource/pkg/eventsource"
)

type SyncEventHandler interface {
	HandleEvent(context.Context, pgx.Tx, *eventsource.Event) error
}

type SyncEventHandlerFunc func(context.Context, pgx.Tx, *eventsource.Event) error

func (f SyncEventHandlerFunc) HandleEvent(
	ctx context.Context, tx pgx.Tx, e *eventsource.Event,
) error {
	return f(ctx, tx, e)
}
