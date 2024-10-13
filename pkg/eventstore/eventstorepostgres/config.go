package eventstorepostgres

import (
	"context"
	"io"
	"log/slog"
)

type config struct {
	context            context.Context
	logger             *slog.Logger
	syncEventHandler   SyncEventHandler
	asyncEventHandlers []AsyncEventHandler
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

func WithSyncEventHandler(handler SyncEventHandler) option {
	return func(cfg *config) {
		cfg.syncEventHandler = handler
	}
}

func WithAsyncEventHandlers(handlers ...AsyncEventHandler) option {
	seen := make(map[string]bool)
	for _, h := range handlers {
		if id := h.SubscriptionID(); !seen[id] {
			seen[id] = true
			continue
		}
		panic("duplicate subscription ID")
	}
	return func(cfg *config) {
		cfg.asyncEventHandlers = handlers
	}
}
