package postgresadapter

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/riverqueue/river"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/rnovatorov/go-eventsource/examples/saga/sagapb"
)

type SagaWorkerArgs struct {
	RequestData  json.RawMessage `json:"request_data"`
	Compensation bool            `json:"compensation"`
}

func (SagaWorkerArgs) Kind() string { return "saga_task" }

type SagaWorker struct {
	river.WorkerDefaults[SagaWorkerArgs]
	requestHandler  RequestHandler
	responseHandler ResponseHandler
}

func (w *SagaWorker) Work(
	ctx context.Context, job *river.Job[SagaWorkerArgs],
) error {
	if job.Args.Compensation {
		return w.executeCompensation(ctx, job.Args.RequestData)
	} else {
		return w.executeTask(ctx, job.Args.RequestData)
	}
}

func (w *SagaWorker) executeTask(ctx context.Context, data json.RawMessage) error {
	var req sagapb.TaskRequest
	if err := protojson.Unmarshal(data, &req); err != nil {
		return fmt.Errorf("unmarshal request: %w", err)
	}

	resp, err := w.requestHandler.HandleTaskRequest(ctx, &req)
	if err != nil {
		return fmt.Errorf("handle request: %w", err)
	}

	if err := w.responseHandler.HandleTaskResponse(ctx, resp); err != nil {
		return fmt.Errorf("handle response: %w", err)
	}

	return nil
}

func (w *SagaWorker) executeCompensation(ctx context.Context, data json.RawMessage) error {
	var req sagapb.CompensationRequest
	if err := protojson.Unmarshal(data, &req); err != nil {
		return fmt.Errorf("unmarshal request: %w", err)
	}

	resp, err := w.requestHandler.HandleCompensationRequest(ctx, &req)
	if err != nil {
		return fmt.Errorf("handle request: %w", err)
	}

	if err := w.responseHandler.HandleCompensationResponse(ctx, resp); err != nil {
		return fmt.Errorf("handle response: %w", err)
	}

	return nil
}
