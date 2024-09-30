package model

import "time"

type Transaction struct {
	Timestamp       time.Time
	AccountDebited  string
	AccountCredited string
	Amount          uint64
}
