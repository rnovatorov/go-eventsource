package eventsource

import "errors"

var (
	ErrConcurrentUpdate       = errors.New("concurrent update")
	ErrAggregateAlreadyExists = errors.New("aggregate already exists")
	ErrAggregateDoesNotExist  = errors.New("aggregate does not exist")
	ErrEmptyID                = errors.New("empty ID")
	ErrUnknownCommand         = errors.New("unknown command")
)
