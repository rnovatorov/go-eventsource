package postgresadapter

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/rnovatorov/go-eventsource/examples/saga/sagapb"
)

type RequestSender struct {
	riverClient *river.Client[pgx.Tx]
}

type RequestSenderParams struct {
	Pool *pgxpool.Pool
}

func StartRequestSender(
	ctx context.Context, p RequestSenderParams,
) (*RequestSender, error) {
	riverClient, err := river.NewClient(riverpgxv5.New(p.Pool), &river.Config{})
	if err != nil {
		return nil, err
	}

	return &RequestSender{
		riverClient: riverClient,
	}, nil
}

func (r *RequestSender) Stop() {
	const timeout = 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	r.riverClient.Stop(ctx)
}

func (r *RequestSender) SendTaskRequest(
	ctx context.Context, req *sagapb.TaskRequest,
) error {
	data, err := protojson.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	if _, err := r.riverClient.Insert(ctx, &SagaWorkerArgs{
		RequestData:  data,
		Compensation: false,
	}, nil); err != nil {
		return fmt.Errorf("insert job: %w", err)
	}

	return nil
}

func (r *RequestSender) SendCompensationRequest(
	ctx context.Context, req *sagapb.CompensationRequest,
) error {
	data, err := protojson.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	if _, err := r.riverClient.Insert(ctx, &SagaWorkerArgs{
		RequestData:  data,
		Compensation: true,
	}, nil); err != nil {
		return fmt.Errorf("insert job: %w", err)
	}

	return nil
}
