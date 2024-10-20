package postgresadapter

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"

	"github.com/rnovatorov/go-eventsource/examples/saga/sagapb"
)

type RequestHandler interface {
	HandleTaskRequest(
		context.Context, *sagapb.TaskRequest,
	) (*sagapb.TaskResponse, error)
	HandleCompensationRequest(
		context.Context, *sagapb.CompensationRequest,
	) (*sagapb.CompensationResponse, error)
}

type ResponseHandler interface {
	HandleTaskResponse(context.Context, *sagapb.TaskResponse) error
	HandleCompensationResponse(context.Context, *sagapb.CompensationResponse) error
}

type ResponseReceiver struct {
	riverClient *river.Client[pgx.Tx]
}

type ResponseReceiverParams struct {
	Pool            *pgxpool.Pool
	RequestHandler  RequestHandler
	ResponseHandler ResponseHandler
}

func StartResponseReceiver(
	ctx context.Context, p ResponseReceiverParams,
) (*RequestSender, error) {
	workers := river.NewWorkers()
	river.AddWorker(workers, &SagaWorker{
		requestHandler:  p.RequestHandler,
		responseHandler: p.ResponseHandler,
	})

	riverClient, err := river.NewClient(riverpgxv5.New(p.Pool), &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 100},
		},
		Workers: workers,
	})
	if err != nil {
		return nil, fmt.Errorf("new river client: %w", err)
	}

	if err := riverClient.Start(ctx); err != nil {
		return nil, fmt.Errorf("start river client: %w", err)
	}

	return &RequestSender{
		riverClient: riverClient,
	}, nil
}

func (r *ResponseReceiver) Stop() {
	const timeout = 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	r.riverClient.Stop(ctx)
}
