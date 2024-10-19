package application

import (
	"context"
	"time"

	"github.com/rnovatorov/go-eventsource/examples/accounting/accountingpb"
	"github.com/rnovatorov/go-eventsource/examples/accounting/model"
	"github.com/rnovatorov/go-eventsource/pkg/eventsource"
	"github.com/rnovatorov/go-eventsource/pkg/eventstore"
)

type App struct {
	bookRepository    *eventsource.AggregateRepository[model.Book, *model.Book]
	projectionQueries ProjectionQueries
}

type Params struct {
	EventStore        eventstore.Interface
	ProjectionQueries ProjectionQueries
}

func New(p Params) *App {
	return &App{
		bookRepository:    eventsource.NewAggregateRepository[model.Book](p.EventStore),
		projectionQueries: p.ProjectionQueries,
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
	accountType accountingpb.AccountType,
) error {
	_, err := a.bookRepository.Update(ctx, bookID,
		model.BookAccountAdd{
			AccountName: accountName,
			AccountType: accountType,
		},
	)
	return err
}

func (a *App) GetBookAccountBalance(
	ctx context.Context, bookID string, accountName string,
) (uint64, error) {
	return a.projectionQueries.GetAccountBalance(ctx, bookID, accountName)
}

func (a *App) EnterBookTransaction(
	ctx context.Context, bookID string, timestamp time.Time,
	accountDebited string, accountCredited string, amount uint64,
) error {
	_, err := a.bookRepository.Update(ctx, bookID,
		model.BookTransactionEnter{Transaction: model.Transaction{
			Timestamp:       timestamp,
			AccountDebited:  accountDebited,
			AccountCredited: accountCredited,
			Amount:          amount,
		}},
	)
	return err
}
