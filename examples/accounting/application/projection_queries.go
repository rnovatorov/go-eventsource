package application

import "context"

type ProjectionQueries interface {
	GetAccountBalance(
		ctx context.Context, bookID string, accountName string,
	) (uint64, error)
}
