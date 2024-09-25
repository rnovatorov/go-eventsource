package model

import "errors"

var (
	ErrAccountNotFound         = errors.New("account not found")
	ErrAccountNameConflict     = errors.New("account name conflict")
	ErrAccountNameEmpty        = errors.New("account name empty")
	ErrAccountTypeUnknown      = errors.New("account type unknown")
	ErrAccountOverdrawn        = errors.New("account overdrawn")
	ErrAccountDebitedNotFound  = errors.New("account debited not found")
	ErrAccountCreditedNotFound = errors.New("account credited not found")
	ErrAccountDebitDeclined    = errors.New("account debit declined")
	ErrAccountCreditDeclined   = errors.New("account credit declined")
	ErrBookClosed              = errors.New("book closed")
	ErrBookAlreadyCreated      = errors.New("book already created")
)
