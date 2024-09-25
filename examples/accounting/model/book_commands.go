package model

import "time"

type BookCreate struct {
	Description string
}

type BookClose struct{}

type BookAccountAdd struct {
	Name string
	Type AccountType
}

type BookTransactionEnter struct {
	Timestamp       time.Time
	AccountDebited  string
	AccountCredited string
	Amount          uint64
}
