package eventsource

import (
	"context"
	"fmt"

	"github.com/rnovatorov/go-eventsource/pkg/eventstore"
)

type Aggregate[T any, R aggregateRoot[T]] struct {
	id           string
	version      int
	root         R
	stateChanges StateChanges
	causationIDs map[string]struct{}
}

func NewAggregate[T any, R aggregateRoot[T]](id string) *Aggregate[T, R] {
	agg, err := RehydrateAggregate[T, R](id, nil)
	if err != nil {
		panic(err)
	}

	return agg
}

func RehydrateAggregate[T any, R aggregateRoot[T]](
	id string, events eventstore.Events,
) (*Aggregate[T, R], error) {
	var root R = new(T)
	var version int
	causationIDs := make(map[string]struct{}, len(events))

	for _, event := range events {
		stateChange, err := event.Data.UnmarshalNew()
		if err != nil {
			return nil, fmt.Errorf("unmarshal state change: %w", err)
		}

		root.ApplyStateChange(stateChange)
		version = event.AggregateVersion

		if cid := event.Metadata.CausationID(); cid != "" {
			causationIDs[cid] = struct{}{}
		}
	}

	return &Aggregate[T, R]{
		id:           id,
		version:      version,
		root:         root,
		stateChanges: nil,
		causationIDs: causationIDs,
	}, nil
}

func (a *Aggregate[T, R]) ID() string {
	return a.id
}

func (a *Aggregate[T, R]) Version() int {
	return a.version
}

func (a *Aggregate[T, R]) Root() R {
	return a.root
}

func (a *Aggregate[T, R]) ProcessCommand(ctx context.Context, cmd Command) error {
	metadata := eventstore.MetadataFromContext(ctx)

	if _, ok := a.causationIDs[metadata.CausationID()]; ok {
		return ErrCommandAlreadyProcessed
	}

	stateChanges, err := a.root.ProcessCommand(cmd)
	if err != nil {
		return fmt.Errorf("%T: %w", cmd, err)
	}

	for _, stateChange := range stateChanges {
		a.root.ApplyStateChange(stateChange)
		a.stateChanges = append(a.stateChanges, stateChange)
		a.version++
	}

	if cid := metadata.CausationID(); cid != "" {
		a.causationIDs[cid] = struct{}{}
	}

	return nil
}
