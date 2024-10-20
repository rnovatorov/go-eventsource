package model_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/rnovatorov/go-eventsource/examples/saga/orchestratorservice/internal/model"
	"github.com/rnovatorov/go-eventsource/examples/saga/sagapb"
	"github.com/rnovatorov/go-eventsource/pkg/eventsource"
)

type SagaSuite struct {
	suite.Suite
	saga *model.SagaAggregate
}

func (s *SagaSuite) SetupTest() {
	s.saga = eventsource.NewAggregate[model.Saga](uuid.NewString())
}

func (s *SagaSuite) TestSequence() {
	s.shouldProcessCommand(model.BeginSaga{
		TaskDefinitions: []*sagapb.TaskDefinition{
			s.task("a"),
			s.task("b", "a"),
			s.task("c", "b"),
		},
	})
	s.Require().Equal([]string(nil), s.saga.Root().TaskTransitiveDependencies("a"))
	s.Require().Equal([]string{"a"}, s.saga.Root().TaskTransitiveDependencies("b"))
	s.Require().Equal([]string{"a", "b"}, s.saga.Root().TaskTransitiveDependencies("c"))
	s.Require().True(s.saga.Root().InProgress())
	s.Require().True(s.saga.Root().TaskInProgress("a"))
	s.Require().False(s.saga.Root().TaskInProgress("b"))
	s.Require().False(s.saga.Root().TaskInProgress("c"))

	s.shouldProcessCommand(model.EndTask{ID: "a", Result: &structpb.Value{}})
	s.Require().True(s.saga.Root().InProgress())
	s.Require().True(s.saga.Root().TaskEnded("a"))
	s.Require().True(s.saga.Root().TaskInProgress("b"))
	s.Require().False(s.saga.Root().TaskInProgress("c"))

	s.shouldProcessCommand(model.EndTask{ID: "b", Result: &structpb.Value{}})
	s.Require().True(s.saga.Root().InProgress())
	s.Require().True(s.saga.Root().TaskEnded("a"))
	s.Require().True(s.saga.Root().TaskEnded("b"))
	s.Require().True(s.saga.Root().TaskInProgress("c"))

	s.shouldProcessCommand(model.EndTask{ID: "c", Result: &structpb.Value{}})
	s.Require().True(s.saga.Root().Ended())
	s.Require().True(s.saga.Root().TaskEnded("a"))
	s.Require().True(s.saga.Root().TaskEnded("b"))
	s.Require().True(s.saga.Root().TaskEnded("c"))
}

func (s *SagaSuite) TestDiamond() {
	s.shouldProcessCommand(model.BeginSaga{
		TaskDefinitions: []*sagapb.TaskDefinition{
			s.task("d", "b", "c"),
			s.task("b", "a"),
			s.task("c", "a"),
			s.task("a"),
		},
	})
	s.Require().Equal([]string(nil), s.saga.Root().TaskTransitiveDependencies("a"))
	s.Require().Equal([]string{"a"}, s.saga.Root().TaskTransitiveDependencies("b"))
	s.Require().Equal([]string{"a"}, s.saga.Root().TaskTransitiveDependencies("c"))
	s.Require().Equal([]string{"a", "b", "c"}, s.saga.Root().TaskTransitiveDependencies("d"))
	s.Require().True(s.saga.Root().InProgress())
	s.Require().True(s.saga.Root().TaskInProgress("a"))
	s.Require().False(s.saga.Root().TaskInProgress("b"))
	s.Require().False(s.saga.Root().TaskInProgress("c"))
	s.Require().False(s.saga.Root().TaskInProgress("d"))

	s.shouldProcessCommand(model.EndTask{ID: "a", Result: &structpb.Value{}})
	s.Require().True(s.saga.Root().InProgress())
	s.Require().True(s.saga.Root().TaskEnded("a"))
	s.Require().True(s.saga.Root().TaskInProgress("b"))
	s.Require().True(s.saga.Root().TaskInProgress("c"))
	s.Require().False(s.saga.Root().TaskInProgress("d"))

	s.shouldProcessCommand(model.EndTask{ID: "b", Result: &structpb.Value{}})
	s.Require().True(s.saga.Root().InProgress())
	s.Require().True(s.saga.Root().TaskEnded("a"))
	s.Require().True(s.saga.Root().TaskEnded("b"))
	s.Require().True(s.saga.Root().TaskInProgress("c"))
	s.Require().False(s.saga.Root().TaskInProgress("d"))

	s.shouldProcessCommand(model.EndTask{ID: "c", Result: &structpb.Value{}})
	s.Require().True(s.saga.Root().InProgress())
	s.Require().True(s.saga.Root().TaskEnded("a"))
	s.Require().True(s.saga.Root().TaskEnded("b"))
	s.Require().True(s.saga.Root().TaskEnded("c"))
	s.Require().True(s.saga.Root().TaskInProgress("d"))

	s.shouldProcessCommand(model.EndTask{ID: "d", Result: &structpb.Value{}})
	s.Require().True(s.saga.Root().Ended())
	s.Require().True(s.saga.Root().TaskEnded("a"))
	s.Require().True(s.saga.Root().TaskEnded("b"))
	s.Require().True(s.saga.Root().TaskEnded("c"))
	s.Require().True(s.saga.Root().TaskEnded("d"))
}

func (s *SagaSuite) TestTwoTasksWithSameID() {
	err := s.saga.ProcessCommand(context.Background(), model.BeginSaga{
		TaskDefinitions: []*sagapb.TaskDefinition{
			s.task("a"),
			s.task("b"),
			s.task("c"),
			s.task("a"),
		},
	})
	s.Require().ErrorIs(err, model.ErrVertexAlreadyAdded)
}

func (s *SagaSuite) TestDependencyNotFound() {
	err := s.saga.ProcessCommand(context.Background(), model.BeginSaga{
		TaskDefinitions: []*sagapb.TaskDefinition{
			s.task("a"),
			s.task("b", "a"),
			s.task("c", "d"),
		},
	})
	s.Require().ErrorIs(err, model.ErrVertexNotFound)
}

func (s *SagaSuite) TestShortCycle() {
	err := s.saga.ProcessCommand(context.Background(), model.BeginSaga{
		TaskDefinitions: []*sagapb.TaskDefinition{
			s.task("a", "c"),
			s.task("b", "a"),
			s.task("c", "b"),
		},
	})
	s.Require().ErrorIs(err, model.ErrDependencyGraphCyclic)
}

func (s *SagaSuite) TestLongCycle() {
	err := s.saga.ProcessCommand(context.Background(), model.BeginSaga{
		TaskDefinitions: []*sagapb.TaskDefinition{
			s.task("a", "b"),
			s.task("b", "c", "d"),
			s.task("c", "e"),
			s.task("d", "e"),
			s.task("e", "b"),
		},
	})
	s.Require().ErrorIs(err, model.ErrDependencyGraphCyclic)
}

func (s *SagaSuite) TestDisconnectedSequences() {
	s.shouldProcessCommand(model.BeginSaga{
		TaskDefinitions: []*sagapb.TaskDefinition{
			s.task("a"),
			s.task("b", "a"),
			s.task("x"),
			s.task("y", "x"),
		},
	})
	s.Require().Equal([]string(nil), s.saga.Root().TaskTransitiveDependencies("a"))
	s.Require().Equal([]string{"a"}, s.saga.Root().TaskTransitiveDependencies("b"))
	s.Require().Equal([]string(nil), s.saga.Root().TaskTransitiveDependencies("x"))
	s.Require().Equal([]string{"x"}, s.saga.Root().TaskTransitiveDependencies("y"))
	s.Require().True(s.saga.Root().InProgress())
	s.Require().True(s.saga.Root().TaskInProgress("a"))
	s.Require().True(s.saga.Root().TaskInProgress("x"))
	s.Require().False(s.saga.Root().TaskInProgress("b"))
	s.Require().False(s.saga.Root().TaskInProgress("y"))

	s.shouldProcessCommand(model.EndTask{ID: "a", Result: &structpb.Value{}})
	s.Require().True(s.saga.Root().InProgress())
	s.Require().True(s.saga.Root().TaskEnded("a"))
	s.Require().True(s.saga.Root().TaskInProgress("x"))
	s.Require().True(s.saga.Root().TaskInProgress("b"))
	s.Require().False(s.saga.Root().TaskInProgress("y"))

	s.shouldProcessCommand(model.EndTask{ID: "x", Result: &structpb.Value{}})
	s.Require().True(s.saga.Root().InProgress())
	s.Require().True(s.saga.Root().TaskEnded("a"))
	s.Require().True(s.saga.Root().TaskEnded("x"))
	s.Require().True(s.saga.Root().TaskInProgress("b"))
	s.Require().True(s.saga.Root().TaskInProgress("y"))

	s.shouldProcessCommand(model.EndTask{ID: "y", Result: &structpb.Value{}})
	s.Require().True(s.saga.Root().InProgress())
	s.Require().True(s.saga.Root().TaskEnded("a"))
	s.Require().True(s.saga.Root().TaskEnded("x"))
	s.Require().True(s.saga.Root().TaskInProgress("b"))
	s.Require().True(s.saga.Root().TaskEnded("y"))

	s.shouldProcessCommand(model.EndTask{ID: "b", Result: &structpb.Value{}})
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

func (s *SagaSuite) task(id string, deps ...string) *sagapb.TaskDefinition {
	return &sagapb.TaskDefinition{
		Id:           id,
		Dependencies: deps,
	}
}

func TestSaga(t *testing.T) {
	suite.Run(t, new(SagaSuite))
}
