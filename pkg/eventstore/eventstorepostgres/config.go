package eventstorepostgres

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/rnovatorov/go-eventsource/pkg/eventsource"
)

type config struct {
	schema            string
	projectionUpdater ProjectionUpdater
}

func newConfig(opts ...option) config {
	cfg := config{
		schema: "eventsource",
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

type option func(*config)

type ProjectionUpdater func(context.Context, pgx.Tx, *eventsource.Event) error

func WithProjectionUpdater(u ProjectionUpdater) option {
	return func(cfg *config) {
		cfg.projectionUpdater = u
	}
}
