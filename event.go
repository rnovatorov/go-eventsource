package eventsourcing

import (
	"time"

	"google.golang.org/protobuf/types/known/anypb"
)

type Event struct {
	ID               string
	AggregateVersion int
	Timestamp        time.Time
	Metadata         Metadata
	Data             *anypb.Any
}
