package model

type Account struct {
	name    string
	type_   AccountType
	balance uint64
}

func (a *Account) Name() string {
	return a.name
}

func (a *Account) Type() AccountType {
	return a.type_
}

func (a *Account) Balance() uint64 {
	return a.balance
}

func (a *Account) mustDebit(amount uint64) {
	newBalance, err := a.canDebit(amount)
	if err != nil {
		panic(err)
	}
	a.balance = newBalance
}

func (a *Account) mustCredit(amount uint64) {
	newBalance, err := a.canCredit(amount)
	if err != nil {
		panic(err)
	}
	a.balance = newBalance
}

func (a *Account) canDebit(amount uint64) (newBalance uint64, err error) {
	switch a.type_ {
	case AccountType_CAPITAL:
		return a.canDecreaseBalance(amount)
	case AccountType_ASSET:
		return a.canIncreaseBalance(amount)
	case AccountType_LIABILITY:
		return a.canDecreaseBalance(amount)
	case AccountType_INCOME:
		return a.canDecreaseBalance(amount)
	case AccountType_EXPENSE:
		return a.canIncreaseBalance(amount)
	case AccountType_UNKNOWN:
		fallthrough
	default:
		return 0, ErrAccountTypeUnknown
	}
}

func (a *Account) canCredit(amount uint64) (newBalance uint64, err error) {
	switch a.type_ {
	case AccountType_CAPITAL:
		return a.canIncreaseBalance(amount)
	case AccountType_ASSET:
		return a.canDecreaseBalance(amount)
	case AccountType_LIABILITY:
		return a.canIncreaseBalance(amount)
	case AccountType_INCOME:
		return a.canIncreaseBalance(amount)
	case AccountType_EXPENSE:
		return a.canDecreaseBalance(amount)
	case AccountType_UNKNOWN:
		fallthrough
	default:
		return 0, ErrAccountTypeUnknown
	}
}

func (a *Account) canIncreaseBalance(amount uint64) (newBalance uint64, err error) {
	return a.balance + amount, nil
}

func (a *Account) canDecreaseBalance(amount uint64) (newBalance uint64, err error) {
	if amount > a.balance {
		return 0, ErrAccountOverdrawn
	}

	return a.balance - amount, nil
}
