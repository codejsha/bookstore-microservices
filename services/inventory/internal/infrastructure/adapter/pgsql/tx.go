package pgsql

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const maxOptimisticRetries = 5

var errVersionConflict = errors.New("pgsql: optimistic version conflict")
var errReservationConflict = errors.New("pgsql: concurrent reservation conflict")
var errStockCreateConflict = errors.New("pgsql: concurrent stock row creation")

func isRetryableConflict(err error) bool {
	return errors.Is(err, errVersionConflict) || errors.Is(err, errStockCreateConflict)
}

const pgUniqueViolationCode = "23505"

func isUniqueViolation(err error, constraint string) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}
	return pgErr.Code == pgUniqueViolationCode && pgErr.ConstraintName == constraint
}
