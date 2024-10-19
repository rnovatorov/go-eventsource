package eventsource

import "errors"

var (
	ErrAggregateAlreadyExists  = errors.New("aggregate already exists")
	ErrAggregateDoesNotExist   = errors.New("aggregate does not exist")
	ErrCommandUnknown          = errors.New("command unknown")
	ErrCommandAlreadyProcessed = errors.New("command already processed")
)
