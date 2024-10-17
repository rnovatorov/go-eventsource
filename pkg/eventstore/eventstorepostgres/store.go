package eventstorepostgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/rnovatorov/go-eventsource/pkg/eventsource"
	"github.com/rnovatorov/pgxlisten"
)

var _ eventsource.EventStore = (*Store)(nil)

type Store struct {
	pool    *pgxpool.Pool
	config  config
	cancel  context.CancelFunc
	stopped chan struct{}
}

func Start(pool *pgxpool.Pool, opts ...option) (*Store, error) {
	cfg := newConfig(opts...)

	for _, handler := range cfg.asyncEventHandlers {
		if _, err := pool.Exec(cfg.context, createSubscriptionQuery,
			pgx.NamedArgs{
				"subscription_id": handler.SubscriptionID(),
			},
		); err != nil {
			return nil, fmt.Errorf("create subscription: %s: %w",
				handler.SubscriptionID(), err)
		}
	}

	ctx, cancel := context.WithCancel(cfg.context)

	s := &Store{
		pool:    pool,
		config:  cfg,
		cancel:  cancel,
		stopped: make(chan struct{}),
	}
	go s.run(ctx)

	return s, nil
}

func (s *Store) Stop() {
	s.cancel()
	<-s.stopped
}

func (s *Store) run(ctx context.Context) {
	defer close(s.stopped)

	listener := pgxlisten.StartListener(s.pool)
	defer listener.Stop()

	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.sequenceEventsLoop(ctx, listener)
	}()

	eventsSequenced := listener.Listen("events.sequenced")
	defer eventsSequenced.Unlisten()

	eventsSequencedFanout := pgxlisten.StartFanout(eventsSequenced)
	defer eventsSequencedFanout.Stop()

	for _, eventHandler := range s.config.asyncEventHandlers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			eventHandler := eventHandler
			s.handleSubscriptionEventsLoop(ctx, eventsSequencedFanout,
				eventHandler)
		}()
	}

	<-ctx.Done()
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
			s.config.logger.ErrorContext(ctx, "failed to sequence events",
				slog.String("error", err.Error()))
		}
	}
}

func (s *Store) sequenceEvents(ctx context.Context) error {
	return pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, acquireEventsAdvisoryLockQuery); err != nil {
			return fmt.Errorf("acquire events advisory lock: %w", err)
		}

		if ct, err := tx.Exec(ctx, sequenceEventsQuery); err != nil {
			return err
		} else if ct.RowsAffected() > 0 {
			if _, err := tx.Exec(ctx, notifyEventsSequencedQuery); err != nil {
				return fmt.Errorf("notify events sequenced: %w", err)
			}
		}

		return nil
	})
}

func (s *Store) handleSubscriptionEventsLoop(
	ctx context.Context, eventsSequencedFanout *pgxlisten.Fanout,
	handler AsyncEventHandler,
) {
	eventsSequenced := eventsSequencedFanout.Listen()
	defer eventsSequenced.Unlisten()

	// FIXME: Hard-code.
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		if err := s.handleSubscriptionEvents(ctx, handler); err != nil {
			s.config.logger.ErrorContext(ctx,
				"failed to handle subscription events",
				slog.String("error", err.Error()),
				slog.String("subscription_id", handler.SubscriptionID()))
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
	ctx context.Context, handler AsyncEventHandler,
) error {
	for {
		if err := s.handleSubscriptionEvent(ctx, handler); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil
			}
			return err
		}
	}
}

func (s *Store) handleSubscriptionEvent(
	ctx context.Context, handler AsyncEventHandler,
) error {
	return pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		rows, _ := tx.Query(ctx, selectNextSubscriptionEventQuery, pgx.NamedArgs{
			"subscription_id": handler.SubscriptionID(),
		})
		event, err := pgx.CollectExactlyOneRow(rows, s.collectEvent)
		if err != nil {
			return fmt.Errorf("select next subscription event: %w", err)
		}

		if err := handler.HandleEvent(ctx, event); err != nil {
			return fmt.Errorf("handle event: %w", err)
		}

		if _, err := tx.Exec(ctx, advanceSubscriptionPositionQuery, pgx.NamedArgs{
			"subscription_id": handler.SubscriptionID(),
		}); err != nil {
			return fmt.Errorf("advance subscription position: %w", err)
		}

		return nil
	})
}

func (s *Store) ListEvents(
	ctx context.Context, aggregateID string,
) (eventsource.Events, error) {
	rows, _ := s.pool.Query(ctx, listEventsQuery, pgx.NamedArgs{
		"aggregate_id": aggregateID,
	})

	return pgx.CollectRows(rows, s.collectEvent)
}

func (s *Store) collectEvent(row pgx.CollectableRow) (*eventsource.Event, error) {
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

		if _, err := tx.Exec(ctx, notifyEventsInsertedQuery); err != nil {
			return fmt.Errorf("notify events inserted: %w", err)
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

	if handler := s.config.syncEventHandler; handler != nil {
		if err := handler.HandleEvent(ctx, tx, event); err != nil {
			return fmt.Errorf("sync event handler: %w", err)
		}
	}

	return nil
}
