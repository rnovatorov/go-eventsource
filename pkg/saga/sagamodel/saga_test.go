package sagamodel_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/rnovatorov/go-eventsource/pkg/eventsource"
	"github.com/rnovatorov/go-eventsource/pkg/saga/sagamodel"
)

type SagaSuite struct {
	suite.Suite
	saga *sagamodel.SagaAggregate
}

func (s *SagaSuite) SetupTest() {
	s.saga = eventsource.NewAggregate[sagamodel.Saga](uuid.NewString())
}

func (s *SagaSuite) TestSequence() {
	s.shouldProcessCommand(sagamodel.BeginSaga{
		TaskDefinitions: sagamodel.TaskDefinitions{
			s.task("a"),
			s.task("b", "a"),
			s.task("c", "b"),
		},
	})
	s.Require().True(s.saga.Root().InProgress())
	s.Require().True(s.saga.Root().TaskInProgress("a"))
	s.Require().False(s.saga.Root().TaskInProgress("b"))
	s.Require().False(s.saga.Root().TaskInProgress("c"))

	s.shouldProcessCommand(sagamodel.EndTask{ID: "a", Result: &structpb.Struct{}})
	s.Require().True(s.saga.Root().InProgress())
	s.Require().True(s.saga.Root().TaskEnded("a"))
	s.Require().True(s.saga.Root().TaskInProgress("b"))
	s.Require().False(s.saga.Root().TaskInProgress("c"))

	s.shouldProcessCommand(sagamodel.EndTask{ID: "b", Result: &structpb.Struct{}})
	s.Require().True(s.saga.Root().InProgress())
	s.Require().True(s.saga.Root().TaskEnded("a"))
	s.Require().True(s.saga.Root().TaskEnded("b"))
	s.Require().True(s.saga.Root().TaskInProgress("c"))

	s.shouldProcessCommand(sagamodel.EndTask{ID: "c", Result: &structpb.Struct{}})
	s.Require().True(s.saga.Root().Ended())
	s.Require().True(s.saga.Root().TaskEnded("a"))
	s.Require().True(s.saga.Root().TaskEnded("b"))
	s.Require().True(s.saga.Root().TaskEnded("c"))
}

func (s *SagaSuite) TestDiamond() {
	s.shouldProcessCommand(sagamodel.BeginSaga{
		TaskDefinitions: sagamodel.TaskDefinitions{
			s.task("d", "b", "c"),
			s.task("b", "a"),
			s.task("c", "a"),
			s.task("a"),
		},
	})
	s.Require().True(s.saga.Root().InProgress())
	s.Require().True(s.saga.Root().TaskInProgress("a"))
	s.Require().False(s.saga.Root().TaskInProgress("b"))
	s.Require().False(s.saga.Root().TaskInProgress("c"))
	s.Require().False(s.saga.Root().TaskInProgress("d"))

	s.shouldProcessCommand(sagamodel.EndTask{ID: "a", Result: &structpb.Struct{}})
	s.Require().True(s.saga.Root().InProgress())
	s.Require().True(s.saga.Root().TaskEnded("a"))
	s.Require().True(s.saga.Root().TaskInProgress("b"))
	s.Require().True(s.saga.Root().TaskInProgress("c"))
	s.Require().False(s.saga.Root().TaskInProgress("d"))

	s.shouldProcessCommand(sagamodel.EndTask{ID: "b", Result: &structpb.Struct{}})
	s.Require().True(s.saga.Root().InProgress())
	s.Require().True(s.saga.Root().TaskEnded("a"))
	s.Require().True(s.saga.Root().TaskEnded("b"))
	s.Require().True(s.saga.Root().TaskInProgress("c"))
	s.Require().False(s.saga.Root().TaskInProgress("d"))

	s.shouldProcessCommand(sagamodel.EndTask{ID: "c", Result: &structpb.Struct{}})
	s.Require().True(s.saga.Root().InProgress())
	s.Require().True(s.saga.Root().TaskEnded("a"))
	s.Require().True(s.saga.Root().TaskEnded("b"))
	s.Require().True(s.saga.Root().TaskEnded("c"))
	s.Require().True(s.saga.Root().TaskInProgress("d"))

	s.shouldProcessCommand(sagamodel.EndTask{ID: "d", Result: &structpb.Struct{}})
	s.Require().True(s.saga.Root().Ended())
	s.Require().True(s.saga.Root().TaskEnded("a"))
	s.Require().True(s.saga.Root().TaskEnded("b"))
	s.Require().True(s.saga.Root().TaskEnded("c"))
	s.Require().True(s.saga.Root().TaskEnded("d"))
}

func (s *SagaSuite) TestTwoTasksWithSameID() {
	err := s.saga.ProcessCommand(context.Background(), sagamodel.BeginSaga{
		TaskDefinitions: sagamodel.TaskDefinitions{
			s.task("a"),
			s.task("b"),
			s.task("c"),
			s.task("a"),
		},
	})
	s.Require().ErrorIs(err, sagamodel.ErrVertexAlreadyAdded)
}

func (s *SagaSuite) TestDependencyNotFound() {
	err := s.saga.ProcessCommand(context.Background(), sagamodel.BeginSaga{
		TaskDefinitions: sagamodel.TaskDefinitions{
			s.task("a"),
			s.task("b", "a"),
			s.task("c", "d"),
		},
	})
	s.Require().ErrorIs(err, sagamodel.ErrVertexNotFound)
}

func (s *SagaSuite) TestShortCycle() {
	err := s.saga.ProcessCommand(context.Background(), sagamodel.BeginSaga{
		TaskDefinitions: sagamodel.TaskDefinitions{
			s.task("a", "c"),
			s.task("b", "a"),
			s.task("c", "b"),
		},
	})
	s.Require().ErrorIs(err, sagamodel.ErrDependencyGraphCyclic)
}

func (s *SagaSuite) TestLongCycle() {
	err := s.saga.ProcessCommand(context.Background(), sagamodel.BeginSaga{
		TaskDefinitions: sagamodel.TaskDefinitions{
			s.task("a", "b"),
			s.task("b", "c", "d"),
			s.task("c", "e"),
			s.task("d", "e"),
			s.task("e", "b"),
		},
	})
	s.Require().ErrorIs(err, sagamodel.ErrDependencyGraphCyclic)
}

func (s *SagaSuite) TestDisconnectedSequences() {
	s.shouldProcessCommand(sagamodel.BeginSaga{
		TaskDefinitions: sagamodel.TaskDefinitions{
			s.task("a"),
			s.task("b", "a"),
			s.task("x"),
			s.task("y", "x"),
		},
	})
	s.Require().True(s.saga.Root().InProgress())
	s.Require().True(s.saga.Root().TaskInProgress("a"))
	s.Require().True(s.saga.Root().TaskInProgress("x"))
	s.Require().False(s.saga.Root().TaskInProgress("b"))
	s.Require().False(s.saga.Root().TaskInProgress("y"))

	s.shouldProcessCommand(sagamodel.EndTask{ID: "a", Result: &structpb.Struct{}})
	s.Require().True(s.saga.Root().InProgress())
	s.Require().True(s.saga.Root().TaskEnded("a"))
	s.Require().True(s.saga.Root().TaskInProgress("x"))
	s.Require().True(s.saga.Root().TaskInProgress("b"))
	s.Require().False(s.saga.Root().TaskInProgress("y"))

	s.shouldProcessCommand(sagamodel.EndTask{ID: "x", Result: &structpb.Struct{}})
	s.Require().True(s.saga.Root().InProgress())
	s.Require().True(s.saga.Root().TaskEnded("a"))
	s.Require().True(s.saga.Root().TaskEnded("x"))
	s.Require().True(s.saga.Root().TaskInProgress("b"))
	s.Require().True(s.saga.Root().TaskInProgress("y"))

	s.shouldProcessCommand(sagamodel.EndTask{ID: "y", Result: &structpb.Struct{}})
	s.Require().True(s.saga.Root().InProgress())
	s.Require().True(s.saga.Root().TaskEnded("a"))
	s.Require().True(s.saga.Root().TaskEnded("x"))
	s.Require().True(s.saga.Root().TaskInProgress("b"))
	s.Require().True(s.saga.Root().TaskEnded("y"))

	s.shouldProcessCommand(sagamodel.EndTask{ID: "b", Result: &structpb.Struct{}})
	s.Require().True(s.saga.Root().Ended())
	s.Require().True(s.saga.Root().TaskEnded("a"))
	s.Require().True(s.saga.Root().TaskEnded("x"))
	s.Require().True(s.saga.Root().TaskEnded("b"))
	s.Require().True(s.saga.Root().TaskEnded("y"))
}

func (s *SagaSuite) shouldProcessCommand(cmd eventsource.Command) {
	err := s.saga.ProcessCommand(context.Background(), cmd)
	s.Require().NoError(err)
}

func (s *SagaSuite) task(id string, deps ...string) *sagamodel.TaskDefinition {
	return &sagamodel.TaskDefinition{
		Id:           id,
		Arguments:    nil,
		Dependencies: deps,
	}
}

func TestSaga(t *testing.T) {
	suite.Run(t, new(SagaSuite))
}
