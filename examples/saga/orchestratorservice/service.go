package orchestratorservice

import (
	"context"
	"fmt"

	"github.com/rnovatorov/go-eventsource/examples/saga/orchestratorservice/internal/model"
	"github.com/rnovatorov/go-eventsource/examples/saga/sagapb"
	"github.com/rnovatorov/go-eventsource/pkg/eventsource"
	"github.com/rnovatorov/go-eventsource/pkg/eventstore"
)

type Service struct {
	sagaRepository *eventsource.AggregateRepository[model.Saga, *model.Saga]
	requestSender  RequestSender
}

type Params struct {
	EventStore    eventstore.Interface
	RequestSender RequestSender
}

func New(p Params) *Service {
	return &Service{
		sagaRepository: eventsource.NewAggregateRepository[model.Saga](p.EventStore),
		requestSender:  p.RequestSender,
	}
}

func (s *Service) BeginSaga(
	ctx context.Context, sagaID string, taskDefinitions []*sagapb.TaskDefinition,
) error {
	_, err := s.sagaRepository.Create(ctx, sagaID, model.BeginSaga{
		TaskDefinitions: taskDefinitions,
	})
	return err
}

func (s *Service) HandleSagaEvent(ctx context.Context, event *eventstore.Event) error {
	data, err := event.Data.UnmarshalNew()
	if err != nil {
		return fmt.Errorf("unmarshal data: %w", err)
	}

	switch d := data.(type) {
	case *sagapb.TaskBegun:
		s.handleTaskBegun(ctx, event, d)
	case *sagapb.CompensationBegun:
		s.handleCompensationBegun(ctx, event, d)
	}

	return nil
}

func (s *Service) handleTaskBegun(
	ctx context.Context, e *eventstore.Event, d *sagapb.TaskBegun,
) error {
	saga, err := s.sagaRepository.Get(ctx, e.AggregateID)
	if err != nil {
		return fmt.Errorf("get saga: %w", err)
	}

	def := saga.Root().TaskDefinition(d.Id)

	return s.requestSender.SendTaskRequest(ctx, &sagapb.TaskRequest{
		SagaId:        saga.ID(),
		TaskId:        d.Id,
		TaskMethod:    def.Method,
		TaskArguments: saga.Root().TaskArguments(d.Id),
	})
}

func (s *Service) handleCompensationBegun(
	ctx context.Context, e *eventstore.Event, d *sagapb.CompensationBegun,
) error {
	saga, err := s.sagaRepository.Get(ctx, e.AggregateID)
	if err != nil {
		return fmt.Errorf("get saga: %w", err)
	}

	def := saga.Root().TaskDefinition(d.Id)

	return s.requestSender.SendCompensationRequest(ctx, &sagapb.CompensationRequest{
		SagaId:                    saga.ID(),
		TaskId:                    d.Id,
		TaskArguments:             saga.Root().TaskArguments(d.Id),
		TaskResult:                saga.Root().TaskResult(d.Id),
		TaskCompensationMethod:    def.CompensationMethod,
		TaskCompensationArguments: saga.Root().TaskCompensationArguments(d.Id),
	})
}

func (s *Service) HandleTaskResponse(
	ctx context.Context, resp *sagapb.TaskResponse,
) error {
	if resp.TaskAborted {
		_, err := s.sagaRepository.Update(ctx, resp.SagaId, model.AbortTask{
			ID:     resp.TaskId,
			Reason: resp.TaskAbortReason,
		})
		return err
	}
	_, err := s.sagaRepository.Update(ctx, resp.SagaId, model.EndTask{
		ID:     resp.TaskId,
		Result: resp.TaskResult,
	})
	return err
}

func (s *Service) HandleCompensationResponse(
	ctx context.Context, resp *sagapb.CompensationResponse,
) error {
	_, err := s.sagaRepository.Update(ctx, resp.SagaId, model.EndCompensation{
		ID: resp.TaskId,
	})
	return err
}
