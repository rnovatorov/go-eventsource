package eventsource

type aggregateRoot[T any] interface {
	*T
	ProcessCommand(Command) (StateChanges, error)
	ApplyStateChange(StateChange)
}

type copiableAggregateRoot[T any] interface {
	aggregateRoot[T]
	Copy() *T
}

func Given[T any, R copiableAggregateRoot[T]](root R, changes StateChanges, f func()) {
	bak := root.Copy()
	for _, change := range changes {
		root.ApplyStateChange(change)
	}
	f()
	*root = *bak
}
