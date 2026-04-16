package database

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

// IsNullErr is error when scaning is there are no rows
func IsNullErr(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
