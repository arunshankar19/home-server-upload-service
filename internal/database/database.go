package database

import "context"

type DB interface {
	Ping(ctx context.Context) error

	Exec(ctx context.Context, query string, args ...any) (CommandTag, error)
	Query(ctx context.Context, query string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) Row
}
