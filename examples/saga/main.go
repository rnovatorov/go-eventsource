package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/rnovatorov/go-eventsource/examples/accounting/application"
	"github.com/rnovatorov/go-eventsource/examples/accounting/httpadapter"
	"github.com/rnovatorov/go-eventsource/examples/accounting/postgresadapter"
	"github.com/rnovatorov/go-eventsource/examples/saga/sagacoordinator"
	"github.com/rnovatorov/go-eventsource/examples/saga/sagapb"
	"github.com/rnovatorov/go-eventsource/examples/saga/watermilladapter"
	"github.com/rnovatorov/go-eventsource/pkg/eventsource"
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

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("new database pool: %w", err)
	}
	defer pool.Close()

	requestsChannel := gochannel.NewGoChannel(gochannel.Config{
		OutputChannelBuffer: 10,
		Persistent:          true,
	}, nil)
	defer requestsChannel.Close()

	responsesChannel := gochannel.NewGoChannel(gochannel.Config{
		OutputChannelBuffer: 10,
		Persistent:          true,
	}, nil)
	defer responsesChannel.Close()

	router, err := message.NewRouter(message.RouterConfig{},
		watermill.NewSlogLogger(logger))
	if err != nil {
		return fmt.Errorf("new router: %w", err)
	}
	router.AddHandler("task_processor",
		"test", requestsChannel,
		"test", responsesChannel,
		func(msg *message.Message) ([]*message.Message, error) {
			var reqAny anypb.Any
			if err := protojson.Unmarshal(msg.Payload, &reqAny); err != nil {
				return nil, fmt.Errorf("unmarshal payload: %w", err)
			}

			req, err := reqAny.UnmarshalNew()
			if err != nil {
				return nil, fmt.Errorf("unmarshal request: %w", err)
			}

			switch r := req.(type) {
			case *sagapb.PerformTaskRequest:
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-time.After(time.Second):
				}
				resultAny, err := anypb.New(wrapperspb.String("result"))
				if err != nil {
					return nil, fmt.Errorf("marshal result: %w", err)
				}
				respAny, err := anypb.New(&sagapb.PerformTaskResponse{
					SagaId:          r.SagaId,
					TaskId:          r.TaskDefinition.Id,
					TaskAborted:     false,
					TaskResult:      resultAny,
					TaskAbortReason: nil,
				})
				if err != nil {
					return nil, fmt.Errorf("marshal response: %w", err)
				}
				payload, err := protojson.Marshal(respAny)
				if err != nil {
					return nil, fmt.Errorf("marshal payload: %w", err)
				}
				return []*message.Message{{
					UUID:    uuid.NewString(),
					Payload: payload,
				}}, nil
			case *sagapb.PerformCompensationRequest:
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-time.After(time.Second):
				}
				respAny, err := anypb.New(&sagapb.PerformCompensationResponse{
					SagaId: r.SagaId,
					TaskId: r.TaskDefinition.Id,
				})
				if err != nil {
					return nil, fmt.Errorf("marshal response: %w", err)
				}
				payload, err := protojson.Marshal(respAny)
				if err != nil {
					return nil, fmt.Errorf("marshal payload: %w", err)
				}
				return []*message.Message{{
					UUID:    uuid.NewString(),
					Payload: payload,
				}}, nil
			default:
				return nil, errors.New("not implemented")
			}
		},
	)
	go router.Run(ctx)
	defer router.Close()

	h := new(H)

	eventStore, err := eventstorepostgres.Start(pool,
		eventstorepostgres.WithLogger(logger),
		eventstorepostgres.WithAsyncEventHandlers(h),
		eventstorepostgres.WithSyncEventHandler(
			postgresadapter.ProjectionUpdater{}))
	if err != nil {
		return fmt.Errorf("start postgres event store: %w", err)
	}
	defer eventStore.Stop()

	app := application.New(application.Params{
		EventStore:        eventStore,
		ProjectionQueries: postgresadapter.NewProjectionQueries(pool),
	})

	requestSender := watermilladapter.NewRequestSender(
		watermilladapter.RequestSenderParams{
			Publisher: requestsChannel,
		})

	coordinator := sagacoordinator.New(sagacoordinator.Params{
		EventStore:    eventStore,
		RequestSender: requestSender,
	})
	h.coord = coordinator

	responseReceiver, err := watermilladapter.StartResponseReceiver(ctx,
		watermilladapter.ResponseReceiverParams{
			Subscriber:     responsesChannel,
			SubscribeTopic: "test",
			Handler:        coordinator,
			Logger:         logger,
		})
	if err != nil {
		return fmt.Errorf("start response receiver: %w", err)
	}
	defer responseReceiver.Stop()

	if err := coordinator.BeginSaga(ctx, "mysagaid14", "wfid",
		[]*sagapb.TaskDefinition{
			{Id: "task1"},
			{Id: "task2", Dependencies: []string{"task1"}},
			{Id: "task3", Dependencies: []string{"task1"}},
			{Id: "task4", Dependencies: []string{"task2", "task3"}},
		},
	); err != nil {
		return fmt.Errorf("begin saga: %w", err)
	}

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

type H struct {
	coord *sagacoordinator.Coordinator
}

func (*H) SubscriptionID() string {
	return "h_sub_id"
}

func (h *H) HandleEvent(ctx context.Context, event *eventsource.Event) error {
	return h.coord.HandleSagaEvent(ctx, event)
}
