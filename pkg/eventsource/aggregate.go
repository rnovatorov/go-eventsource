package eventsource

type Aggregate[T any, R aggregateRoot[T]] struct {
	id           string
	version      int
	root         R
	stateChanges StateChanges
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

func (a *Aggregate[T, R]) changeState(cmd Command) error {
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
