package model

import "github.com/rnovatorov/go-eventsource/examples/accounting/accountingpb"

type AccountType = accountingpb.Book_AccountType

const (
	AccountTypeUnknown   = accountingpb.Book_ACCOUNT_TYPE_UNKNOWN
	AccountTypeCapital   = accountingpb.Book_ACCOUNT_TYPE_CAPITAL
	AccountTypeAsset     = accountingpb.Book_ACCOUNT_TYPE_ASSET
	AccountTypeLiability = accountingpb.Book_ACCOUNT_TYPE_LIABILITY
	AccountTypeIncome    = accountingpb.Book_ACCOUNT_TYPE_INCOME
	AccountTypeExpence   = accountingpb.Book_ACCOUNT_TYPE_EXPENSE
)

func NewAccountType(s string) AccountType {
	return AccountType(accountingpb.Book_AccountType_value[s])
}
