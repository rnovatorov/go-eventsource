package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rnovatorov/go-eventsource/examples/accounting/application"
	"github.com/rnovatorov/go-eventsource/examples/accounting/httpadapter"
	"github.com/rnovatorov/go-eventsource/examples/accounting/postgresadapter"
	"github.com/rnovatorov/go-eventsource/pkg/eventstore/eventstorepostgres"
)

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("new database pool: %w", err)
	}
	defer pool.Close()

	eventStore := eventstorepostgres.New(pool,
		eventstorepostgres.WithSaveEventHook(postgresadapter.UpdateProjections))

	app := application.New(application.Params{
		EventStore:        eventStore,
		ProjectionQueries: postgresadapter.NewProjectionQueries(pool),
	})

	server := &http.Server{
		Addr:        os.Getenv("HTTP_SERVER_LISTEN_ADDRESS"),
		Handler:     httpadapter.NewHandler(app),
		BaseContext: func(net.Listener) context.Context { return ctx },
	}
	go func() {
		<-ctx.Done()
		server.Shutdown(ctx)
	}()
	return server.ListenAndServe()
}
