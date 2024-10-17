package sagacoordinator

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/anypb"

	"github.com/rnovatorov/go-eventsource/examples/saga/sagamodel"
	"github.com/rnovatorov/go-eventsource/examples/saga/sagapb"
	"github.com/rnovatorov/go-eventsource/pkg/eventsource"
)

type Coordinator struct {
	sagaRepository *eventsource.AggregateRepository[sagamodel.Saga, *sagamodel.Saga]
	requestSender  RequestSender
}

type Params struct {
	EventStore    eventsource.EventStore
	RequestSender RequestSender
}

func New(p Params) *Coordinator {
	return &Coordinator{
		sagaRepository: eventsource.NewAggregateRepository[sagamodel.Saga](p.EventStore),
		requestSender:  p.RequestSender,
	}
}

func (c *Coordinator) BeginSaga(
	ctx context.Context, sagaID string, workflowID string,
	taskDefinitions []*sagapb.TaskDefinition,
) error {
	_, err := c.sagaRepository.Create(ctx, sagaID, sagamodel.BeginSaga{
		TaskDefinitions: taskDefinitions,
	})
	return err
}

func (c *Coordinator) HandleSagaEvent(ctx context.Context, event *eventsource.Event) error {
	data, err := event.Data.UnmarshalNew()
	if err != nil {
		return fmt.Errorf("unmarshal data: %w", err)
	}

	switch d := data.(type) {
	case *sagapb.TaskBegun:
		c.handleTaskBegun(ctx, event, d)
	case *sagapb.CompensationBegun:
		c.handleCompensationBegun(ctx, event, d)
	}

	return nil
}

func (c *Coordinator) handleTaskBegun(
	ctx context.Context, e *eventsource.Event, d *sagapb.TaskBegun,
) error {
	saga, err := c.sagaRepository.Get(ctx, e.AggregateID)
	if err != nil {
		return fmt.Errorf("get saga: %w", err)
	}

	def := saga.Root().TaskDefinition(d.Id)

	req, err := anypb.New(&sagapb.PerformTaskRequest{
		SagaId:        saga.ID(),
		TaskId:        d.Id,
		TaskMethod:    def.Method,
		TaskArguments: saga.Root().TaskArguments(d.Id),
	})
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}

	return c.requestSender.SendRequest(ctx, req)
}

func (c *Coordinator) handleCompensationBegun(
	ctx context.Context, e *eventsource.Event, d *sagapb.CompensationBegun,
) error {
	saga, err := c.sagaRepository.Get(ctx, e.AggregateID)
	if err != nil {
		return fmt.Errorf("get saga: %w", err)
	}

	def := saga.Root().TaskDefinition(d.Id)

	req, err := anypb.New(&sagapb.PerformCompensationRequest{
		SagaId:                    saga.ID(),
		TaskId:                    d.Id,
		TaskArguments:             saga.Root().TaskArguments(d.Id),
		TaskResult:                saga.Root().TaskResult(d.Id),
		TaskCompensationMethod:    def.CompensationMethod,
		TaskCompensationArguments: saga.Root().TaskCompensationArguments(d.Id),
	})
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}

	return c.requestSender.SendRequest(ctx, req)
}

func (c *Coordinator) HandlePerformTaskResponse(
	ctx context.Context, resp *sagapb.PerformTaskResponse,
) error {
	if resp.TaskAborted {
		_, err := c.sagaRepository.Update(ctx, resp.SagaId, sagamodel.AbortTask{
			ID:     resp.TaskId,
			Reason: resp.TaskAbortReason,
		})
		return err
	}
	_, err := c.sagaRepository.Update(ctx, resp.SagaId, sagamodel.EndTask{
		ID:     resp.TaskId,
		Result: resp.TaskResult,
	})
	return err
}

func (c *Coordinator) HandlePerformCompensationResponse(
	ctx context.Context, resp *sagapb.PerformCompensationResponse,
) error {
	_, err := c.sagaRepository.Update(ctx, resp.SagaId, sagamodel.EndCompensation{
		ID: resp.TaskId,
	})
	return err
}
