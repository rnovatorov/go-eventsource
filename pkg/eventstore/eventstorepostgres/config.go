package eventstorepostgres

import "github.com/jackc/pgx/v5"

type config struct {
	schema pgx.Identifier
}

func newConfig(opts ...option) config {
	c := config{
		schema: pgx.Identifier{},
	}
	for _, opt := range opts {
		opt(&c)
	}
	return c
}

type option func(*config)

func WithSchema(schema string) option {
	return func(p *config) {
		p.schema = pgx.Identifier{schema}
	}
}
