package sagaservice

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/structpb"

	"github.com/rnovatorov/go-eventsource/pkg/eventsource"
	"github.com/rnovatorov/go-eventsource/pkg/saga/sagamodel"
)

type Service struct {
	sagaRepository *eventsource.AggregateRepository[sagamodel.Saga, *sagamodel.Saga]
	taskExecutor   TaskExecutor
}

type Params struct {
	EventStore   eventsource.EventStore
	TaskExecutor TaskExecutor
}

func New(p Params) *Service {
	return &Service{
		sagaRepository: eventsource.NewAggregateRepository[sagamodel.Saga](p.EventStore),
		taskExecutor:   p.TaskExecutor,
	}
}

func (s *Service) BeginSaga(
	ctx context.Context, id string, taskDefinitions []*sagamodel.TaskDefinition,
) error {
	_, err := s.sagaRepository.Create(ctx, id, sagamodel.BeginSaga{
		TaskDefinitions: taskDefinitions,
	})
	return err
}

func (s *Service) EndTask(
	ctx context.Context, sagaID string, taskID string, result *structpb.Struct,
) error {
	_, err := s.sagaRepository.Update(ctx, sagaID, sagamodel.EndTask{
		ID:     taskID,
		Result: result,
	})
	return err
}

func (s *Service) AbortTask(
	ctx context.Context, sagaID string, taskID string, reason *structpb.Struct,
) error {
	_, err := s.sagaRepository.Update(ctx, sagaID, sagamodel.AbortTask{
		ID:     taskID,
		Reason: reason,
	})
	return err
}

func (s *Service) EndCompensation(
	ctx context.Context, sagaID string, taskID string,
) error {
	_, err := s.sagaRepository.Update(ctx, sagaID, sagamodel.EndCompensation{
		ID: taskID,
	})
	return err
}

func (s *Service) HandleEvent(ctx context.Context, event *eventsource.Event) error {
	data, err := event.Data.UnmarshalNew()
	if err != nil {
		return fmt.Errorf("unmarshal data: %w", err)
	}

	switch d := data.(type) {
	case *sagamodel.TaskBegun:
		s.handleTaskBegun(ctx, event, d)
	case *sagamodel.CompensationBegun:
		s.handleCompensationBegun(ctx, event, d)
	}

	return nil
}

func (s *Service) handleTaskBegun(
	ctx context.Context, e *eventsource.Event, d *sagamodel.TaskBegun,
) error {
	sagaCtx, err := s.newSagaContext(ctx, e.AggregateID)
	if err != nil {
		return fmt.Errorf("new saga context: %w", err)
	}

	return s.taskExecutor.BeginTask(sagaCtx, d.Id)
}

func (s *Service) handleCompensationBegun(
	ctx context.Context, e *eventsource.Event, d *sagamodel.CompensationBegun,
) error {
	sagaCtx, err := s.newSagaContext(ctx, e.AggregateID)
	if err != nil {
		return fmt.Errorf("new saga context: %w", err)
	}

	return s.taskExecutor.BeginCompensation(sagaCtx, d.Id)
}

func (s *Service) newSagaContext(ctx context.Context, id string) (SagaContext, error) {
	saga, err := s.sagaRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &sagaContext{
		Context: ctx,
		saga:    saga,
	}, nil
}
