package eventsourcing

type aggregateRoot[T any] interface {
	*T
	ProcessCommand(any) ([]StateChange, error)
	ApplyStateChange(StateChange)
}
