package eventsourcing

type Aggregate[T any, R aggregateRoot[T]] struct {
	id           string
	version      int
	root         R
	stateChanges []StateChange
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

func (a *Aggregate[T, R]) StateChanges() []StateChange {
	return a.stateChanges
}

func (a *Aggregate[T, R]) ChangeState(cmd any) error {
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
