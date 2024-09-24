package eventsource

type aggregateRoot[T any] interface {
	*T
	ProcessCommand(Command) ([]StateChange, error)
	ApplyStateChange(StateChange)
}
