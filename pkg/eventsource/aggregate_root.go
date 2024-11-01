package eventsource

type aggregateRoot[T any] interface {
	*T
	ProcessCommand(Command) (StateChanges, error)
	ApplyStateChange(StateChange)
}
