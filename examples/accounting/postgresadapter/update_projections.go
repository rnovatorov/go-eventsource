package postgresadapter

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/rnovatorov/go-eventsource/examples/accounting/accountingpb"
	"github.com/rnovatorov/go-eventsource/pkg/eventstore"
)

func UpdateProjections(
	ctx context.Context, tx pgx.Tx, event *eventstore.Event,
) error {
	data, err := event.Data.UnmarshalNew()
	if err != nil {
		return fmt.Errorf("unmarshal data: %w", err)
	}

	switch d := data.(type) {
	case *accountingpb.BookCreated:
		return handleBookCreated(ctx, tx, event, d)
	case *accountingpb.BookClosed:
		return handleBookClosed(ctx, tx, event, d)
	case *accountingpb.BookAccountAdded:
		return handleBookAccountAdded(ctx, tx, event, d)
	case *accountingpb.BookTransactionEntered:
		return handleBookTransactionEntered(ctx, tx, event, d)
	}

	return nil
}

func handleBookCreated(
	ctx context.Context, tx pgx.Tx, e *eventstore.Event,
	d *accountingpb.BookCreated,
) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO books (id, closed, description)
		VALUES ($1, false, $2)
	`, e.AggregateID, d.Description)
	return err
}

func handleBookClosed(
	ctx context.Context, tx pgx.Tx, e *eventstore.Event,
	_ *accountingpb.BookClosed,
) error {
	_, err := tx.Exec(ctx, `
		UPDATE books
		SET closed = true
		WHERE id = $1
	`, e.AggregateID)
	return err
}

func handleBookAccountAdded(
	ctx context.Context, tx pgx.Tx, e *eventstore.Event,
	d *accountingpb.BookAccountAdded,
) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO accounts (book_id, name, type, balance)
		VALUES ($1, $2, $3, 0)
	`, e.AggregateID, d.Name, d.Type.String())
	return err
}

func handleBookTransactionEntered(
	ctx context.Context, tx pgx.Tx, e *eventstore.Event,
	d *accountingpb.BookTransactionEntered,
) error {
	if _, err := tx.Exec(ctx, `
		INSERT INTO transactions (book_id, timestamp,
			account_debited, account_credited, amount)
		VALUES ($1, $2, $3, $4, $5)
	`, e.AggregateID, d.Timestamp.AsTime(),
		d.AccountDebited, d.AccountCredited, d.Amount); err != nil {
		return fmt.Errorf("insert transaction: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		UPDATE accounts
		SET balance = $3
		WHERE book_id = $1 AND name = $2
	`, e.AggregateID, d.AccountDebited, d.AccountDebitedNewBalance); err != nil {
		return fmt.Errorf("update debited account: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		UPDATE accounts
		SET balance = $3
		WHERE book_id = $1 AND name = $2
	`, e.AggregateID, d.AccountCredited, d.AccountCreditedNewBalance); err != nil {
		return fmt.Errorf("update credited account: %w", err)
	}

	return nil
}
