package eventstore

import (
	"context"
)

type Interface interface {
	ListEvents(
		ctx context.Context, aggregateID string,
	) (Events, error)
	SaveEvents(
		ctx context.Context, aggregateID string, expectedAggregateVersion int,
		events Events,
	) error
}
