package application

import (
	"context"
	"fmt"
	"time"

	"github.com/rnovatorov/go-eventsource/examples/accounting/model"
	"github.com/rnovatorov/go-eventsource/pkg/eventsource"
)

type App struct {
	bookRepository *eventsource.AggregateRepository[model.Book, *model.Book]
}

type Params struct {
	EventStore eventsource.EventStore
}

func New(p Params) *App {
	return &App{
		bookRepository: eventsource.NewAggregateRepository[model.Book](p.EventStore),
	}
}

func (a *App) CreateBook(
	ctx context.Context, bookID string, bookDescription string,
) (string, error) {
	book, err := a.bookRepository.Create(ctx, bookID,
		model.BookCreate{
			Description: bookDescription,
		},
	)
	if err != nil {
		return "", err
	}

	return book.ID(), nil
}

func (a *App) CloseBook(
	ctx context.Context, bookID string,
) error {
	_, err := a.bookRepository.Update(ctx, bookID, model.BookClose{})
	return err
}

func (a *App) AddBookAccount(
	ctx context.Context, bookID string, accountName string,
	accountType model.AccountType,
) error {
	_, err := a.bookRepository.Update(ctx, bookID,
		model.BookAccountAdd{
			Name: accountName,
			Type: accountType,
		},
	)
	return err
}

func (a *App) GetBookAccountBalance(
	ctx context.Context, bookID string, accountName string,
) (uint64, error) {
	book, err := a.bookRepository.Get(ctx, bookID)
	if err != nil {
		return 0, fmt.Errorf("get book: %w", err)
	}

	account, err := book.Root().AccountByName(accountName)
	if err != nil {
		return 0, fmt.Errorf("get account by name: %w", err)
	}

	return account.Balance(), nil
}

func (a *App) EnterBookTransaction(
	ctx context.Context, bookID string, timestamp time.Time,
	accountDebited string, accountCredited string, amount uint64,
) error {
	_, err := a.bookRepository.Update(ctx, bookID,
		model.BookTransactionEnter{
			Timestamp:       timestamp,
			AccountDebited:  accountDebited,
			AccountCredited: accountCredited,
			Amount:          amount,
		},
	)
	return err
}
