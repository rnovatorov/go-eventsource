package eventstorepostgres

import (
	"context"
	"io"
	"log/slog"
)

type config struct {
	context       context.Context
	logger        *slog.Logger
	saveEventHook SaveEventHook
}

func newConfig(opts ...option) config {
	cfg := config{
		context: context.Background(),
		logger:  slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

type option func(*config)

func WithContext(ctx context.Context) option {
	return func(cfg *config) {
		cfg.context = ctx
	}
}

func WithLogger(logger *slog.Logger) option {
	return func(cfg *config) {
		cfg.logger = logger
	}
}

func WithSaveEventHook(hook SaveEventHook) option {
	return func(cfg *config) {
		cfg.saveEventHook = hook
	}
}
