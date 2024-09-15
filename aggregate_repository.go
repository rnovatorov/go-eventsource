package eventsourcing

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/anypb"
)

func NewAggregateRepository[T any, R aggregateRoot[T]](
	eventStore EventStore,
) *AggregateRepository[T, R] {
	return &AggregateRepository[T, R]{
		eventStore: eventStore,
	}
}

type AggregateRepository[T any, R aggregateRoot[T]] struct {
	eventStore EventStore
}

func (r *AggregateRepository[T, R]) Get(
	ctx context.Context, id string,
) (*Aggregate[T, R], error) {
	agg, err := r.load(ctx, id)
	if err != nil {
		return nil, err
	}

	if agg.Version() == 0 {
		return nil, ErrAggregateDoesNotExist
	}

	return agg, nil
}

func (r *AggregateRepository[T, R]) Create(
	ctx context.Context, id string, cmd any,
) (*Aggregate[T, R], error) {
	agg, err := r.load(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}

	if agg.Version() != 0 {
		return nil, ErrAggregateAlreadyExists
	}

	if err := agg.ChangeState(cmd); err != nil {
		return nil, fmt.Errorf("change state: %w", err)
	}

	if err := r.Save(ctx, agg); err != nil {
		if errors.Is(err, ErrConcurrentUpdate) {
			return nil, ErrAggregateAlreadyExists
		}
		return nil, fmt.Errorf("save: %w", err)
	}

	return agg, nil
}

func (r *AggregateRepository[T, R]) GetOrCreate(
	ctx context.Context, id string, cmd any,
) (*Aggregate[T, R], error) {
	agg, err := r.load(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}

	if agg.Version() > 0 {
		return agg, nil
	}

	if err := agg.ChangeState(cmd); err != nil {
		return nil, fmt.Errorf("change state: %w", err)
	}

	if err := r.Save(ctx, agg); err != nil {
		if errors.Is(err, ErrConcurrentUpdate) {
			agg, err = r.load(ctx, id)
			if err != nil {
				return nil, fmt.Errorf("load: %w", err)
			}
			return agg, nil
		}
		return nil, fmt.Errorf("save: %w", err)
	}

	return agg, nil
}

func (r *AggregateRepository[T, R]) Update(
	ctx context.Context, id string, cmd any,
) (*Aggregate[T, R], error) {
	agg, err := r.load(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}

	if agg.Version() == 0 {
		return nil, ErrAggregateDoesNotExist
	}

	if err := agg.ChangeState(cmd); err != nil {
		return nil, fmt.Errorf("change state: %w", err)
	}

	if err := r.Save(ctx, agg); err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}

	return agg, nil
}

func (r *AggregateRepository[T, R]) load(
	ctx context.Context, id string,
) (*Aggregate[T, R], error) {
	var root R = new(T)
	var version int

	events, err := r.eventStore.ListEvents(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("list events: %w", err)
	}

	for _, event := range events {
		stateChange, err := event.Data.UnmarshalNew()
		if err != nil {
			return nil, fmt.Errorf("unmarshal state change: %w", err)
		}
		root.ApplyStateChange(stateChange)
		version = event.AggregateVersion
	}

	return &Aggregate[T, R]{
		id:           id,
		version:      version,
		root:         root,
		stateChanges: nil,
	}, nil
}

func (r *AggregateRepository[T, R]) Save(
	ctx context.Context, agg *Aggregate[T, R],
) error {
	if len(agg.StateChanges()) == 0 {
		return nil
	}

	originalVersion := agg.Version() - len(agg.StateChanges())
	events := make([]*Event, 0, len(agg.StateChanges()))

	for i, stateChange := range agg.StateChanges() {
		id, err := uuid.NewRandom()
		if err != nil {
			return fmt.Errorf("generate event ID: %w", err)
		}
		data, err := anypb.New(stateChange)
		if err != nil {
			return fmt.Errorf("marshal state change: %w", err)
		}
		events = append(events, &Event{
			ID:               id.String(),
			AggregateVersion: originalVersion + i + 1,
			Timestamp:        time.Now(),
			Metadata:         MetadataFromContext(ctx),
			Data:             data,
		})
	}

	if err := r.eventStore.SaveEvents(
		ctx, agg.ID(), agg.Version(), events,
	); err != nil {
		return fmt.Errorf("save events: %w", err)
	}

	agg.stateChanges = nil

	return nil
}
