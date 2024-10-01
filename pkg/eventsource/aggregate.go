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
	if id, ok := metadata[CausationID].(string); ok {
		if _, ok := a.causationIDs[id]; ok {
			return ErrDuplicateCommand
		}
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

	return nil
}
