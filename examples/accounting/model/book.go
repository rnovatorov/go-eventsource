package model

import (
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/rnovatorov/go-eventsource/examples/accounting/accountingpb"
	"github.com/rnovatorov/go-eventsource/pkg/eventsource"
)

type Book struct {
	created     bool
	closed      bool
	description string
	accounts    map[string]*Account
}

func (b *Book) Closed() bool {
	return b.closed
}

func (b *Book) Description() string {
	return b.description
}

func (b *Book) AccountByName(name string) (*Account, error) {
	if account, ok := b.accounts[name]; ok {
		return account, nil
	}

	return nil, ErrAccountNotFound
}

func (b *Book) ProcessCommand(
	command eventsource.Command,
) (eventsource.StateChanges, error) {
	switch cmd := command.(type) {
	case BookCreate:
		return b.processCreate(cmd)
	case BookClose:
		return b.processClose(cmd)
	case BookAccountAdd:
		return b.processAccountAdd(cmd)
	case BookTransactionEnter:
		return b.processTransactionEnter(cmd)
	default:
		return nil, fmt.Errorf("%w: %T", eventsource.ErrUnknownCommand, cmd)
	}
}

func (b *Book) processCreate(cmd BookCreate) (eventsource.StateChanges, error) {
	if b.created {
		return nil, ErrBookAlreadyCreated
	}

	if b.closed {
		return nil, ErrBookClosed
	}

	return eventsource.StateChanges{
		&accountingpb.Book_Created{
			Description: cmd.Description,
		},
	}, nil
}

func (b *Book) processClose(BookClose) (eventsource.StateChanges, error) {
	if b.closed {
		return nil, ErrBookClosed
	}

	return eventsource.StateChanges{
		&accountingpb.Book_Closed{},
	}, nil
}

func (b *Book) processAccountAdd(cmd BookAccountAdd) (eventsource.StateChanges, error) {
	if b.closed {
		return nil, ErrBookClosed
	}

	if _, ok := b.accounts[cmd.Name]; ok {
		return nil, ErrAccountNameConflict
	}

	if cmd.Name == "" {
		return nil, ErrAccountNameEmpty
	}

	if cmd.Type == AccountTypeUnknown {
		return nil, ErrAccountTypeUnknown
	}

	return eventsource.StateChanges{
		&accountingpb.Book_AccountAdded{
			Name: cmd.Name,
			Type: cmd.Type,
		},
	}, nil
}

func (b *Book) processTransactionEnter(
	cmd BookTransactionEnter,
) (eventsource.StateChanges, error) {
	if b.closed {
		return nil, ErrBookClosed
	}

	accountDebited, ok := b.accounts[cmd.AccountDebited]
	if !ok {
		return nil, ErrAccountDebitedNotFound
	}

	accountCredited, ok := b.accounts[cmd.AccountCredited]
	if !ok {
		return nil, ErrAccountCreditedNotFound
	}

	if _, err := accountDebited.canDebit(cmd.Amount); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrAccountDebitDeclined, err)
	}

	if _, err := accountCredited.canCredit(cmd.Amount); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrAccountCreditDeclined, err)
	}

	return eventsource.StateChanges{
		&accountingpb.Book_TransactionEntered{
			Timestamp:       timestamppb.New(cmd.Timestamp),
			AccountDebited:  cmd.AccountDebited,
			AccountCredited: cmd.AccountCredited,
			Amount:          cmd.Amount,
		},
	}, nil
}

func (b *Book) ApplyStateChange(stateChange eventsource.StateChange) {
	switch sc := stateChange.(type) {
	case *accountingpb.Book_Created:
		b.applyCreated(sc)
	case *accountingpb.Book_Closed:
		b.applyClosed(sc)
	case *accountingpb.Book_AccountAdded:
		b.applyAccountAdded(sc)
	case *accountingpb.Book_TransactionEntered:
		b.applyTransactionEntered(sc)
	default:
		panic(fmt.Sprintf("unexpected state change: %T", sc))
	}
}

func (b *Book) applyCreated(sc *accountingpb.Book_Created) {
	b.created = true
	b.description = sc.Description
	b.accounts = make(map[string]*Account)
}

func (b *Book) applyClosed(*accountingpb.Book_Closed) {
	b.closed = true
}

func (b *Book) applyAccountAdded(sc *accountingpb.Book_AccountAdded) {
	b.accounts[sc.Name] = &Account{
		name:    sc.Name,
		type_:   sc.Type,
		balance: 0,
	}
}

func (b *Book) applyTransactionEntered(sc *accountingpb.Book_TransactionEntered) {
	b.accounts[sc.AccountDebited].mustDebit(sc.Amount)
	b.accounts[sc.AccountCredited].mustCredit(sc.Amount)
}
