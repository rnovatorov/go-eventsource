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
	pool            *pgxpool.Pool
	aggregatesTable string
	eventsTable     string
}

func New(pool *pgxpool.Pool, opts ...option) *Store {
	cfg := newConfig(opts...)

	return &Store{
		pool:            pool,
		aggregatesTable: append(cfg.schema, "aggregates").Sanitize(),
		eventsTable:     append(cfg.schema, "events").Sanitize(),
	}
}

func (s *Store) MigrateDatabase(ctx context.Context) error {
	return pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `
			CREATE TABLE IF NOT EXISTS `+s.aggregatesTable+` (
				id TEXT NOT NULL,
				version INT NOT NULL CHECK (version >= 0),
				PRIMARY KEY (id)
			)
		`); err != nil {
			return fmt.Errorf("create aggregates table: %w", err)
		}

		if _, err := tx.Exec(ctx, `
			CREATE TABLE IF NOT EXISTS `+s.eventsTable+` (
				id TEXT NOT NULL,
				aggregate_id TEXT NOT NULL REFERENCES `+s.aggregatesTable+` (id),
				aggregate_version INT NOT NULL,
				timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
				metadata JSONB NOT NULL,
				data JSONB NOT NULL,
				PRIMARY KEY (id),
				UNIQUE (aggregate_id, aggregate_version)
			)
		`); err != nil {
			return fmt.Errorf("create events table: %w", err)
		}

		return nil
	})
}

func (s *Store) ListEvents(
	ctx context.Context, aggregateID string,
) (eventsource.Events, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, aggregate_version, timestamp, metadata, data
		FROM `+s.eventsTable+`
		WHERE aggregate_id = $1
	`, aggregateID)
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
			if _, err := tx.Exec(ctx, `
				INSERT INTO `+s.aggregatesTable+` (id, version)
				VALUES ($1, 0)
				ON CONFLICT DO NOTHING
			`, aggregateID); err != nil {
				return fmt.Errorf("insert aggregate: %w", err)
			}
		}

		newVersion := expectedAggregateVersion + len(events)
		if ct, err := tx.Exec(ctx, `
			UPDATE `+s.aggregatesTable+`
			SET version = $3
			WHERE id = $1 AND version = $2
		`, aggregateID, expectedAggregateVersion, newVersion); err != nil {
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

	if _, err := tx.Exec(ctx, `
		INSERT INTO `+s.eventsTable+` (id, aggregate_id, aggregate_version,
			timestamp, metadata, data)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, event.ID, event.AggregateID, event.AggregateVersion,
		event.Timestamp, string(metadataBytes), string(dataBytes)); err != nil {
		return fmt.Errorf("insert: %w", err)
	}

	return nil
}
