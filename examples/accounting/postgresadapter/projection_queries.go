package postgresadapter

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectionQueries struct {
	pool *pgxpool.Pool
}

func NewProjectionQueries(pool *pgxpool.Pool) *ProjectionQueries {
	return &ProjectionQueries{
		pool: pool,
	}
}

func (q *ProjectionQueries) GetAccountBalance(
	ctx context.Context, bookID string, accountName string,
) (uint64, error) {
	var balance uint64

	if err := q.pool.QueryRow(ctx, `
		SELECT balance
		FROM accounts
		WHERE book_id = $1 AND name = $2
	`, bookID, accountName).Scan(&balance); err != nil {
		return 0, err
	}

	return balance, nil
}
