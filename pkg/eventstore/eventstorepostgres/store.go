package eventstorepostgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/rnovatorov/go-eventsource/pkg/eventsource"
)

var _ eventsource.EventStore = (*Store)(nil)

type Store struct {
	pool          *pgxpool.Pool
	saveEventHook SaveEventHook
}

func New(pool *pgxpool.Pool, opts ...option) *Store {
	cfg := newConfig(opts...)

	return &Store{
		pool:          pool,
		saveEventHook: cfg.saveEventHook,
	}
}

func (s *Store) ListEvents(
	ctx context.Context, aggregateID string,
) (eventsource.Events, error) {
	rows, err := s.pool.Query(ctx, listEventsQuery, pgx.NamedArgs{
		"aggregate_id": aggregateID,
	})
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return pgx.CollectRows(rows, func(
		row pgx.CollectableRow,
	) (*eventsource.Event, error) {
		var id string
		var aggregateVersion int
		var timestamp time.Time
		var metadataBytes []byte
		var dataBytes []byte

		if err := row.Scan(
			&id, &aggregateVersion, &timestamp, &metadataBytes, &dataBytes,
		); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}

		var metadata eventsource.Metadata
		if err := json.Unmarshal(metadataBytes, &metadata); err != nil {
			return nil, fmt.Errorf("unmarshal metadata: %w", err)
		}

		var data anypb.Any
		if err := protojson.Unmarshal(dataBytes, &data); err != nil {
			return nil, fmt.Errorf("unmarshal data: %w", err)
		}

		return &eventsource.Event{
			ID:               id,
			AggregateID:      aggregateID,
			AggregateVersion: aggregateVersion,
			Timestamp:        timestamp,
			Metadata:         metadata,
			Data:             &data,
		}, nil
	})
}

func (s *Store) SaveEvents(
	ctx context.Context, aggregateID string, expectedAggregateVersion int,
	events eventsource.Events,
) error {
	return pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		if expectedAggregateVersion == 0 {
			if _, err := tx.Exec(ctx, createAggregateQuery, pgx.NamedArgs{
				"aggregate_id": aggregateID,
			}); err != nil {
				return fmt.Errorf("create aggregate: %w", err)
			}
		}

		newVersion := expectedAggregateVersion + len(events)

		if ct, err := tx.Exec(ctx, updateAggregateVersionQuery, pgx.NamedArgs{
			"aggregate_id":               aggregateID,
			"expected_aggregate_version": expectedAggregateVersion,
			"new_aggregate_version":      newVersion,
		}); err != nil {
			return fmt.Errorf("update aggregate version: %w", err)
		} else if ct.RowsAffected() == 0 {
			return eventsource.ErrConcurrentUpdate
		}

		for i, event := range events {
			if err := s.saveEvent(ctx, tx, event); err != nil {
				return fmt.Errorf("%d: %w", i, err)
			}
		}

		return nil
	})
}

func (s *Store) saveEvent(
	ctx context.Context, tx pgx.Tx, event *eventsource.Event,
) error {
	dataBytes, err := protojson.Marshal(event.Data)
	if err != nil {
		return fmt.Errorf("marshal data: %w", err)
	}

	metadataBytes, err := json.Marshal(event.Metadata)
	if err != nil {
		return fmt.Errorf("marshal metadata: %w", err)
	}

	if _, err := tx.Exec(ctx, saveEventQuery, pgx.NamedArgs{
		"id":                event.ID,
		"aggregate_id":      event.AggregateID,
		"aggregate_version": event.AggregateVersion,
		"timestamp":         event.Timestamp,
		"metadata":          string(metadataBytes),
		"data":              string(dataBytes),
	}); err != nil {
		return err
	}

	if err := s.saveEventHook(ctx, tx, event); err != nil {
		return fmt.Errorf("save event hook: %w", err)
	}

	return nil
}
