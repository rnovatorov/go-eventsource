package eventstoreinmemory

import (
	"context"
	"sync"

	"github.com/rnovatorov/go-eventsource/pkg/eventstore"
)

var _ eventstore.Interface = (*Store)(nil)

type Store struct {
	mu         sync.RWMutex
	aggregates map[string]*aggregate
}

func New() *Store {
	return &Store{
		aggregates: make(map[string]*aggregate),
	}
}

func (s *Store) ListEvents(
	ctx context.Context, aggregateID string,
) (eventstore.Events, error) {
	agg := s.getAggregate(aggregateID)
	if agg == nil {
		return nil, nil
	}

	agg.RLock()
	defer agg.RUnlock()

	return agg.events, nil
}

func (s *Store) SaveEvents(
	ctx context.Context, aggregateID string, expectedAggregateVersion int,
	events eventstore.Events,
) error {
	agg := s.getOrCreateAggregate(aggregateID)

	agg.Lock()
	defer agg.Unlock()

	if agg.version != expectedAggregateVersion {
		return eventstore.ErrConcurrentUpdate
	}

	for _, event := range events {
		agg.events = append(agg.events, event)
		agg.version++
	}

	return nil
}

func (s *Store) getAggregate(aggregateID string) *aggregate {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.aggregates[aggregateID]
}

func (s *Store) getOrCreateAggregate(aggregateID string) *aggregate {
	if agg := s.getAggregate(aggregateID); agg != nil {
		return agg
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if agg := s.aggregates[aggregateID]; agg != nil {
		return agg
	}

	agg := new(aggregate)
	s.aggregates[aggregateID] = agg
	return agg
}
