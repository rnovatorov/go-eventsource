package eventsource

import "context"

type Aggregate[T any, R aggregateRoot[T]] struct {
	id           string
	version      int
	root         R
	stateChanges StateChanges
	causationIDs map[string]struct{}
}

func (a *Aggregate[T, R]) ID() string {
	return a.id
}

func (a *Aggregate[T, R]) Version() int {
	return a.version
}

func (a *Aggregate[T, R]) Root() R {
	return a.root
}

func (a *Aggregate[T, R]) ChangeState(ctx context.Context, cmd Command) error {
	metadata := MetadataFromContext(ctx)

	causationID, _ := metadata[CausationID].(string)
	if _, ok := a.causationIDs[causationID]; ok {
		return ErrDuplicateCommand
	}

	stateChanges, err := a.root.ProcessCommand(cmd)
	if err != nil {
		return err
	}

	for _, stateChange := range stateChanges {
		a.root.ApplyStateChange(stateChange)
		a.stateChanges = append(a.stateChanges, stateChange)
		a.version++
	}

	a.causationIDs[causationID] = struct{}{}

	return nil
}
