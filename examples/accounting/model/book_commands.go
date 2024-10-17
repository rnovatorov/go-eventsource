package model

import "github.com/rnovatorov/go-eventsource/examples/accounting/accountingpb"

type BookCreate struct {
	Description string
}

type BookClose struct{}

type BookAccountAdd struct {
	AccountName string
	AccountType accountingpb.AccountType
}

type BookTransactionEnter struct {
	Transaction Transaction
}
