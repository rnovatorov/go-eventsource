package eventstorepostgres

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/rnovatorov/go-eventsource/pkg/eventstore"
)

type SaveEventHook = func(context.Context, pgx.Tx, *eventstore.Event) error
