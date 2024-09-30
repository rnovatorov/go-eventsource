package eventsource

import (
	"time"

	"google.golang.org/protobuf/types/known/anypb"
)

type Event struct {
	ID               string
	AggregateID      string
	AggregateVersion int
	Timestamp        time.Time
	Metadata         Metadata
	Data             *anypb.Any
}

type Events []*Event
