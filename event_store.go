package eventsourcing

import (
	"context"
)

type EventStore interface {
	ListEvents(
		ctx context.Context, aggregateID string,
	) ([]*Event, error)
	SaveEvents(
		ctx context.Context, aggregateID string, expectedAggregateVersion int,
		events []*Event,
	) error
}
