package sagamodel

import (
	"fmt"
	"sort"

	"google.golang.org/protobuf/types/known/structpb"

	"github.com/rnovatorov/go-eventsource/examples/saga/sagapb"
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

	args := new(structpb.Struct)

	for _, arg := range t.definition.Arguments {
		switch v := arg.Value.(type) {
		case *sagapb.ArgumentsDefinition_Static:
			args.Fields[arg.Name] = v.Static
		case *sagapb.ArgumentsDefinition_ResultReference:
			args.Fields[arg.Name] = s.TaskResult(v.ResultReference)
		default:
			panic("unreachable")
		}
	}

	return args
}

func (s *Saga) TaskCompensationArguments(id string) *structpb.Struct {
	t, ok := s.tasks[id]
	if !ok {
		return nil
	}

	args := new(structpb.Struct)

	for _, arg := range t.definition.CompensationArguments {
		switch v := arg.Value.(type) {
		case *sagapb.ArgumentsDefinition_Static:
			args.Fields[arg.Name] = v.Static
		case *sagapb.ArgumentsDefinition_ResultReference:
			args.Fields[arg.Name] = s.TaskResult(v.ResultReference)
		default:
			panic("unreachable")
		}
	}

	return args
}

func (s *Saga) TaskTransitiveDependencies(id string) []string {
	if _, ok := s.tasks[id]; !ok {
		return nil
	}

	var result []string
	seen := make(map[string]bool)

	var dfs func(string)
	dfs = func(id string) {
		for _, dep := range s.tasks[id].definition.Dependencies {
			if !seen[dep] {
				seen[dep] = true
				dfs(dep)
				result = append(result, dep)
			}
		}
	}
	dfs(id)

	return result
}

func (s *Saga) TaskDefinition(id string) *sagapb.TaskDefinition {
	t, ok := s.tasks[id]
	if !ok {
		return nil
	}
	return t.definition
}

func (s *Saga) TaskResult(id string) *structpb.Value {
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

	for _, t := range cmd.TaskDefinitions {
		argDefined := make(map[string]bool)
		for _, arg := range t.Arguments {
			if argDefined[arg.Name] {
				return nil, fmt.Errorf("%w: %s",
					ErrTaskArgumentAlreadyDefined, arg.Name)
			}
			argDefined[arg.Name] = true
			if ref := arg.GetResultReference(); ref != "" {
				if !graph.dependsOn(t.Id, ref) {
					return nil, fmt.Errorf("%w: %s",
						ErrTaskArgumentReferencesIndependentResult, ref)
				}
			}
		}
		compArgDefined := make(map[string]bool)
		for _, arg := range t.CompensationArguments {
			if compArgDefined[arg.Name] {
				return nil, fmt.Errorf("%w: %s",
					ErrTaskCompensationArgumentAlreadyDefined, arg.Name)
			}
			compArgDefined[arg.Name] = true
			if ref := arg.GetResultReference(); ref != "" {
				if ref != t.Id && !graph.dependsOn(t.Id, ref) {
					return nil, fmt.Errorf("%w: %s",
						ErrTaskCompensationArgumentReferencesIndependentResult, ref)
				}
			}
		}
	}

	stateChanges := eventsource.StateChanges{&sagapb.SagaBegun{
		TaskDefinitions: cmd.TaskDefinitions,
	}}

	eventsource.Given(s, stateChanges, func() {
		for _, id := range s.beginnableTasks() {
			stateChanges = append(stateChanges, &sagapb.TaskBegun{
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

	stateChanges := eventsource.StateChanges{&sagapb.TaskEnded{
		Id:     cmd.ID,
		Result: cmd.Result,
	}}

	eventsource.Given(s, stateChanges, func() {
		if s.aborted {
			for _, id := range s.beginnableCompensations() {
				stateChanges = append(stateChanges, &sagapb.CompensationBegun{
					Id: id,
				})
			}
			return
		}
		taskBegun := false
		for _, id := range s.beginnableTasks() {
			stateChanges = append(stateChanges, &sagapb.TaskBegun{
				Id: id,
			})
			taskBegun = true
		}
		if !taskBegun && !s.hasTasksInProgress() {
			stateChanges = append(stateChanges, &sagapb.SagaEnded{})
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

	stateChanges := eventsource.StateChanges{&sagapb.TaskAborted{
		Id:     cmd.ID,
		Reason: cmd.Reason,
	}}

	eventsource.Given(s, stateChanges, func() {
		compensationBegun := false
		for _, id := range s.beginnableCompensations() {
			stateChanges = append(stateChanges, &sagapb.CompensationBegun{
				Id: id,
			})
			compensationBegun = true
		}
		if !compensationBegun && !s.hasTasksInProgress() && !s.hasCompensationsInProgress() {
			stateChanges = append(stateChanges, &sagapb.SagaEnded{})
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

	stateChanges := eventsource.StateChanges{&sagapb.CompensationEnded{
		Id: cmd.ID,
	}}

	eventsource.Given(s, stateChanges, func() {
		compensationBegun := false
		for _, id := range s.beginnableCompensations() {
			stateChanges = append(stateChanges, &sagapb.CompensationBegun{
				Id: id,
			})
			compensationBegun = true
		}
		if !compensationBegun && !s.hasTasksInProgress() && !s.hasCompensationsInProgress() {
			stateChanges = append(stateChanges, &sagapb.SagaEnded{})
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
		for _, dep := range t.definition.Dependencies {
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
		for _, dep := range t.definition.Dependencies {
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
	case *sagapb.SagaBegun:
		s.applySagaBegun(sc)
	case *sagapb.SagaEnded:
		s.applySagaEnded(sc)
	case *sagapb.TaskBegun:
		s.applyTaskBegun(sc)
	case *sagapb.TaskEnded:
		s.applyTaskEnded(sc)
	case *sagapb.TaskAborted:
		s.applyTaskAborted(sc)
	case *sagapb.CompensationBegun:
		s.applyCompensationBegun(sc)
	case *sagapb.CompensationEnded:
		s.applyCompensationEnded(sc)
	default:
		panic(fmt.Sprintf("unexpected state change: %T", sc))
	}
}

func (s *Saga) applySagaBegun(sc *sagapb.SagaBegun) {
	s.begun = true
	s.tasks = make(map[string]*task, len(sc.TaskDefinitions))
	for _, def := range sc.TaskDefinitions {
		s.tasks[def.Id] = &task{
			definition: def,
		}
	}
}

func (s *Saga) applySagaEnded(*sagapb.SagaEnded) {
	s.ended = true
}

func (s *Saga) applyTaskBegun(sc *sagapb.TaskBegun) {
	s.tasks[sc.Id].begun = true
}

func (s *Saga) applyTaskEnded(sc *sagapb.TaskEnded) {
	s.tasks[sc.Id].result = sc.Result
	s.tasks[sc.Id].ended = true
}

func (s *Saga) applyTaskAborted(sc *sagapb.TaskAborted) {
	s.tasks[sc.Id].aborted = true
	s.tasks[sc.Id].abortReason = sc.Reason
	if !s.aborted {
		s.aborted = true
	}
}

func (s *Saga) applyCompensationBegun(sc *sagapb.CompensationBegun) {
	s.tasks[sc.Id].compensationBegun = true
}

func (s *Saga) applyCompensationEnded(sc *sagapb.CompensationEnded) {
	s.tasks[sc.Id].compensationEnded = true
}
