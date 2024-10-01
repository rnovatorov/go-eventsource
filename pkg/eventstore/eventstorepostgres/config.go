package eventstorepostgres

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/rnovatorov/go-eventsource/pkg/eventsource"
)

type config struct {
	schema        string
	saveEventHook SaveEventHook
}

func newConfig(opts ...option) config {
	cfg := config{
		schema: "eventsource",
		saveEventHook: func(context.Context, pgx.Tx, *eventsource.Event) error {
			// noop
			return nil
		},
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

type option func(*config)

type SaveEventHook func(context.Context, pgx.Tx, *eventsource.Event) error

func WithSaveEventHook(hook SaveEventHook) option {
	return func(cfg *config) {
		cfg.saveEventHook = hook
	}
}
