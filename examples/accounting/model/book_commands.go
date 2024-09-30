package model

type BookCreate struct {
	Description string
}

type BookClose struct{}

type BookAccountAdd struct {
	AccountName string
	AccountType AccountType
}

type BookTransactionEnter struct {
	Transaction Transaction
}
