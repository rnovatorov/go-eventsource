package sagamodel

//go:generate protoc -I ../../../proto --go_out=. --go_opt=Msaga_events.proto=../sagamodel saga_events.proto

import (
	"fmt"
	"sort"

	"google.golang.org/protobuf/types/known/structpb"

	"github.com/rnovatorov/go-eventsource/pkg/eventsource"
)

type SagaAggregate = eventsource.Aggregate[Saga, *Saga]

type Saga struct {
	tasks   map[string]*task
	begun   bool
	aborted bool
	ended   bool
}

func (s *Saga) Copy() *Saga {
	tasks := make(map[string]*task, len(s.tasks))
	for id, t := range s.tasks {
		tasks[id] = t.copy()
	}
	return &Saga{
		tasks:   s.tasks,
		begun:   s.begun,
		aborted: s.aborted,
		ended:   s.ended,
	}
}

func (s *Saga) Begun() bool {
	return s.begun
}

func (s *Saga) Ended() bool {
	return s.ended
}

func (s *Saga) Aborted() bool {
	return s.aborted
}

func (s *Saga) InProgress() bool {
	return s.begun && !s.ended && !s.aborted
}

func (s *Saga) CompensationInProgress() bool {
	for _, t := range s.tasks {
		if t.compensationInProgress() {
			return true
		}
	}
	return false
}

func (s *Saga) TaskEnded(id string) bool {
	t, ok := s.tasks[id]
	return ok && t.ended
}

func (s *Saga) TaskInProgress(id string) bool {
	t, ok := s.tasks[id]
	return ok && t.inProgress()
}

func (s *Saga) TaskArguments(id string) *structpb.Struct {
	t, ok := s.tasks[id]
	if !ok {
		return nil
	}
	return t.arguments
}

func (s *Saga) TaskResult(id string) *structpb.Struct {
	t, ok := s.tasks[id]
	if !ok {
		return nil
	}
	return t.result
}

func (s *Saga) CompensationBegun(id string) bool {
	t, ok := s.tasks[id]
	return ok && t.compensationBegun
}

func (s *Saga) CompensationEnded(id string) bool {
	t, ok := s.tasks[id]
	return ok && t.compensationEnded
}

func (s *Saga) ProcessCommand(
	command eventsource.Command,
) (eventsource.StateChanges, error) {
	switch cmd := command.(type) {
	case BeginSaga:
		return s.processBegin(cmd)
	case EndTask:
		return s.processEndTask(cmd)
	case AbortTask:
		return s.processAbortTask(cmd)
	case EndCompensation:
		return s.processEndCompensation(cmd)
	default:
		return nil, fmt.Errorf("%w: %T", eventsource.ErrCommandUnknown, cmd)
	}
}

func (s *Saga) processBegin(cmd BeginSaga) (eventsource.StateChanges, error) {
	if s.begun {
		return nil, ErrSagaAlreadyBegun
	}

	if len(cmd.TaskDefinitions) == 0 {
		return nil, ErrTaskDefinitionsMissing
	}

	graph, err := buildDependencyGraph(cmd.TaskDefinitions)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDependencyGraphBuildingFailed, err)
	}
	if graph.hasCycles() {
		return nil, ErrDependencyGraphCyclic
	}

	stateChanges := eventsource.StateChanges{&SagaBegun{
		TaskDefinitions: cmd.TaskDefinitions,
	}}

	eventsource.Given(s, stateChanges, func() {
		for _, id := range s.beginnableTasks() {
			stateChanges = append(stateChanges, &TaskBegun{
				Id: id,
			})
		}
	})

	return stateChanges, nil
}

func (s *Saga) processEndTask(cmd EndTask) (eventsource.StateChanges, error) {
	if s.ended {
		return nil, ErrSagaEnded
	}

	t, ok := s.tasks[cmd.ID]
	if !ok {
		return nil, ErrTaskNotDefined
	}
	if !t.begun {
		return nil, ErrTaskNotBegun
	}
	if t.ended {
		return nil, ErrTaskAlreadyEnded
	}
	if t.aborted {
		return nil, ErrTaskAborted
	}

	if cmd.Result == nil {
		return nil, ErrTaskResultMissing
	}

	stateChanges := eventsource.StateChanges{&TaskEnded{
		Id:     cmd.ID,
		Result: cmd.Result,
	}}

	eventsource.Given(s, stateChanges, func() {
		if s.aborted {
			for _, id := range s.beginnableCompensations() {
				stateChanges = append(stateChanges, &CompensationBegun{
					Id: id,
				})
			}
			return
		}
		taskBegun := false
		for _, id := range s.beginnableTasks() {
			stateChanges = append(stateChanges, &TaskBegun{
				Id: id,
			})
			taskBegun = true
		}
		if !taskBegun && !s.hasTasksInProgress() {
			stateChanges = append(stateChanges, &SagaEnded{})
		}
	})

	return stateChanges, nil
}

func (s *Saga) processAbortTask(cmd AbortTask) (eventsource.StateChanges, error) {
	if s.ended {
		return nil, ErrSagaEnded
	}

	t, ok := s.tasks[cmd.ID]
	if !ok {
		return nil, ErrTaskNotDefined
	}
	if !t.begun {
		return nil, ErrTaskNotBegun
	}
	if t.ended {
		return nil, ErrTaskEnded
	}
	if t.aborted {
		return nil, ErrTaskAlreadyAborted
	}

	if cmd.Reason == nil {
		return nil, ErrTaskAbortReasonMissing
	}

	stateChanges := eventsource.StateChanges{&TaskAborted{
		Id:     cmd.ID,
		Reason: cmd.Reason,
	}}

	eventsource.Given(s, stateChanges, func() {
		compensationBegun := false
		for _, id := range s.beginnableCompensations() {
			stateChanges = append(stateChanges, &CompensationBegun{
				Id: id,
			})
			compensationBegun = true
		}
		if !compensationBegun && !s.hasTasksInProgress() && !s.hasCompensationsInProgress() {
			stateChanges = append(stateChanges, &SagaEnded{})
		}
	})

	return stateChanges, nil
}

func (s *Saga) processEndCompensation(
	cmd EndCompensation,
) (eventsource.StateChanges, error) {
	if s.ended {
		return nil, ErrSagaEnded
	}

	t, ok := s.tasks[cmd.ID]
	if !ok {
		return nil, ErrTaskNotDefined
	}
	if !t.compensationBegun {
		return nil, ErrCompensationNotBegun
	}
	if t.compensationEnded {
		return nil, ErrCompensationAlreadyEnded
	}

	stateChanges := eventsource.StateChanges{&CompensationEnded{
		Id: cmd.ID,
	}}

	eventsource.Given(s, stateChanges, func() {
		compensationBegun := false
		for _, id := range s.beginnableCompensations() {
			stateChanges = append(stateChanges, &CompensationBegun{
				Id: id,
			})
			compensationBegun = true
		}
		if !compensationBegun && !s.hasTasksInProgress() && !s.hasCompensationsInProgress() {
			stateChanges = append(stateChanges, &SagaEnded{})
		}
	})

	return stateChanges, nil
}

func (s *Saga) hasTasksInProgress() bool {
	for _, t := range s.tasks {
		if t.inProgress() {
			return true
		}
	}
	return false
}

func (s *Saga) hasCompensationsInProgress() bool {
	for _, t := range s.tasks {
		if t.compensationInProgress() {
			return true
		}
	}
	return false
}

func (s *Saga) beginnableTasks() (ids []string) {
	if !s.begun {
		return nil
	}
	if s.aborted {
		return nil
	}
	for id, t := range s.tasks {
		if t.begun {
			continue
		}
		depsSatisfied := true
		for _, dep := range t.dependencies {
			if !s.tasks[dep].ended {
				depsSatisfied = false
				break
			}
		}
		if depsSatisfied {
			ids = append(ids, id)
		}
	}
	sort.Strings(ids)
	return ids
}

func (s *Saga) beginnableCompensations() (ids []string) {
	if !s.aborted {
		return nil
	}
	for id, t := range s.tasks {
		if !t.begun {
			continue
		}
		if t.aborted {
			continue
		}
		if !t.ended {
			continue
		}
		if t.compensationBegun {
			continue
		}
		depsSatisfied := true
		for _, dep := range t.dependencies {
			if !s.tasks[dep].begun {
				continue
			}
			if s.tasks[dep].aborted {
				continue
			}
			if s.tasks[dep].compensationEnded {
				continue
			}
			depsSatisfied = false
			break
		}
		if depsSatisfied {
			ids = append(ids, id)
		}
	}
	sort.Strings(ids)
	return ids
}

func (s *Saga) ApplyStateChange(stateChange eventsource.StateChange) {
	switch sc := stateChange.(type) {
	case *SagaBegun:
		s.applySagaBegun(sc)
	case *SagaEnded:
		s.applySagaEnded(sc)
	case *TaskBegun:
		s.applyTaskBegun(sc)
	case *TaskEnded:
		s.applyTaskEnded(sc)
	case *TaskAborted:
		s.applyTaskAborted(sc)
	case *CompensationBegun:
		s.applyCompensationBegun(sc)
	case *CompensationEnded:
		s.applyCompensationEnded(sc)
	default:
		panic(fmt.Sprintf("unexpected state change: %T", sc))
	}
}

func (s *Saga) applySagaBegun(sc *SagaBegun) {
	s.begun = true
	s.tasks = make(map[string]*task, len(sc.TaskDefinitions))
	for _, def := range sc.TaskDefinitions {
		s.tasks[def.Id] = &task{
			arguments:    def.Arguments,
			dependencies: def.Dependencies,
		}
	}
}

func (s *Saga) applySagaEnded(*SagaEnded) {
	s.ended = true
}

func (s *Saga) applyTaskBegun(sc *TaskBegun) {
	s.tasks[sc.Id].begun = true
}

func (s *Saga) applyTaskEnded(sc *TaskEnded) {
	s.tasks[sc.Id].result = sc.Result
	s.tasks[sc.Id].ended = true
}

func (s *Saga) applyTaskAborted(sc *TaskAborted) {
	s.tasks[sc.Id].aborted = true
	s.tasks[sc.Id].abortReason = sc.Reason
	if !s.aborted {
		s.aborted = true
	}
}

func (s *Saga) applyCompensationBegun(sc *CompensationBegun) {
	s.tasks[sc.Id].compensationBegun = true
}

func (s *Saga) applyCompensationEnded(sc *CompensationEnded) {
	s.tasks[sc.Id].compensationEnded = true
}
