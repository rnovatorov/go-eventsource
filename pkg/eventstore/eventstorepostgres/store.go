package eventstorepostgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rnovatorov/go-routine"
	"github.com/rnovatorov/pgxlisten"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/rnovatorov/go-eventsource/pkg/eventstore"
)

var _ eventstore.Interface = (*Store)(nil)

type Store struct {
	routines                   *routine.Group
	pool                       *pgxpool.Pool
	config                     config
	eventsSequencedFanout      *pgxlisten.Fanout
	eventsSequencedFanoutReady chan struct{}
}

func Start(pool *pgxpool.Pool, opts ...option) *Store {
	cfg := newConfig(opts...)

	s := &Store{
		routines:                   routine.NewGroup(cfg.context),
		pool:                       pool,
		config:                     cfg,
		eventsSequencedFanoutReady: make(chan struct{}),
	}
	s.routines.Go(s.run)

	return s
}

func (s *Store) Stop() {
	s.routines.Stop()
}

func (s *Store) run(ctx context.Context) error {
	listener := pgxlisten.StartListener(s.pool)
	defer listener.Stop()

	eventsSequenced := listener.Listen("events.sequenced")
	defer eventsSequenced.Unlisten()

	s.eventsSequencedFanout = pgxlisten.StartFanout(eventsSequenced)
	defer s.eventsSequencedFanout.Stop()

	close(s.eventsSequencedFanoutReady)

	s.sequenceEventsLoop(ctx, listener)
	return nil
}

func (s *Store) Subscribe(
	ctx context.Context, subscriptionID string, handler eventstore.EventHandler,
) error {
	if _, err := s.pool.Exec(ctx, createSubscriptionQuery, pgx.NamedArgs{
		"subscription_id": subscriptionID,
	}); err != nil {
		return fmt.Errorf("create subscription: %w", err)
	}

	s.routines.Go(func(ctx context.Context) error {
		s.handleSubscriptionEventsLoop(ctx, subscriptionID, handler)
		return nil
	})

	return nil
}

func (s *Store) sequenceEventsLoop(ctx context.Context, listener *pgxlisten.Listener) {
	eventsInserted := listener.Listen("events.inserted")
	defer eventsInserted.Unlisten()

	// FIXME: Hard-code.
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		case <-eventsInserted.Notifications():
		}
		if err := s.sequenceEvents(ctx); err != nil {
			s.config.logger.ErrorContext(ctx,
				"failed to sequence events",
				slog.String("error", err.Error()))
		}
	}
}

func (s *Store) sequenceEvents(ctx context.Context) error {
	return pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		if _, err := tx.Exec(
			ctx, acquireEventsAdvisoryLockQuery,
		); err != nil {
			return fmt.Errorf("acquire events advisory lock: %w", err)
		}

		if ct, err := tx.Exec(ctx, sequenceEventsQuery); err != nil {
			return err
		} else if ct.RowsAffected() > 0 {
			if _, err := tx.Exec(
				ctx, notifyEventsSequencedQuery,
			); err != nil {
				return fmt.Errorf("notify events sequenced: %w", err)
			}
		}

		return nil
	})
}

func (s *Store) handleSubscriptionEventsLoop(
	ctx context.Context, subscriptionID string, handler eventstore.EventHandler,
) {
	select {
	case <-ctx.Done():
		return
	case <-s.eventsSequencedFanoutReady:
	}

	eventsSequenced := s.eventsSequencedFanout.Listen()
	defer eventsSequenced.Unlisten()

	// FIXME: Hard-code.
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		if err := s.handleSubscriptionEvents(
			ctx, subscriptionID, handler,
		); err != nil {
			s.config.logger.ErrorContext(ctx,
				"failed to handle subscription events",
				slog.String("error", err.Error()),
				slog.String("subscription_id", subscriptionID))
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		case <-eventsSequenced.Notifications():
		}
	}
}

func (s *Store) handleSubscriptionEvents(
	ctx context.Context, subscriptionID string, handler eventstore.EventHandler,
) error {
	for {
		if err := s.handleSubscriptionEvent(
			ctx, subscriptionID, handler,
		); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil
			}
			return err
		}
	}
}

func (s *Store) handleSubscriptionEvent(
	ctx context.Context, subscriptionID string, handler eventstore.EventHandler,
) error {
	return pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		rows, _ := tx.Query(ctx, selectNextSubscriptionEventQuery,
			pgx.NamedArgs{
				"subscription_id": subscriptionID,
			})
		event, err := pgx.CollectExactlyOneRow(rows, s.collectEvent)
		if err != nil {
			return fmt.Errorf("select next subscription event: %w", err)
		}

		if err := handler(ctx, event); err != nil {
			return fmt.Errorf("event handler: %w", err)
		}

		if _, err := tx.Exec(ctx, advanceSubscriptionPositionQuery,
			pgx.NamedArgs{
				"subscription_id": subscriptionID,
			},
		); err != nil {
			return fmt.Errorf("advance subscription position: %w", err)
		}

		return nil
	})
}

func (s *Store) ListEvents(
	ctx context.Context, aggregateID string,
) (eventstore.Events, error) {
	rows, _ := s.pool.Query(ctx, listEventsQuery, pgx.NamedArgs{
		"aggregate_id": aggregateID,
	})

	return pgx.CollectRows(rows, s.collectEvent)
}

func (s *Store) collectEvent(row pgx.CollectableRow) (*eventstore.Event, error) {
	var id string
	var aggregateID string
	var aggregateVersion int
	var timestamp time.Time
	var metadataBytes []byte
	var dataBytes []byte

	if err := row.Scan(
		&id, &aggregateID, &aggregateVersion, &timestamp, &metadataBytes,
		&dataBytes,
	); err != nil {
		return nil, fmt.Errorf("scan row: %w", err)
	}

	var metadata eventstore.Metadata
	if err := json.Unmarshal(metadataBytes, &metadata); err != nil {
		return nil, fmt.Errorf("unmarshal metadata: %w", err)
	}

	var data anypb.Any
	if err := protojson.Unmarshal(dataBytes, &data); err != nil {
		return nil, fmt.Errorf("unmarshal data: %w", err)
	}

	return &eventstore.Event{
		ID:               id,
		AggregateID:      aggregateID,
		AggregateVersion: aggregateVersion,
		Timestamp:        timestamp,
		Metadata:         metadata,
		Data:             &data,
	}, nil
}

func (s *Store) SaveEvents(
	ctx context.Context, aggregateID string, expectedAggregateVersion int,
	events eventstore.Events,
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
			return eventstore.ErrConcurrentUpdate
		}

		for i, event := range events {
			if err := s.saveEvent(ctx, tx, event); err != nil {
				return fmt.Errorf("%d: %w", i, err)
			}
		}

		if _, err := tx.Exec(ctx, notifyEventsInsertedQuery); err != nil {
			return fmt.Errorf("notify events inserted: %w", err)
		}

		return nil
	})
}

func (s *Store) saveEvent(
	ctx context.Context, tx pgx.Tx, event *eventstore.Event,
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

	if hook := s.config.saveEventHook; hook != nil {
		if err := hook(ctx, tx, event); err != nil {
			return fmt.Errorf("save event hook: %w", err)
		}
	}

	return nil
}
