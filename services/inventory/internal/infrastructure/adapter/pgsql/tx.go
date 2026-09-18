package pgsql

import (
	"context"
	"errors"
	mathrand "math/rand/v2"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

const maxOptimisticRetries = 5

const (
	optimisticRetryBaseDelay = 5 * time.Millisecond
	optimisticRetryMaxDelay  = 80 * time.Millisecond
)

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

func optimisticRetryDelay(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}
	delay := optimisticRetryBaseDelay
	for i := 0; i < attempt && delay < optimisticRetryMaxDelay; i++ {
		delay *= 2
	}
	if delay > optimisticRetryMaxDelay {
		delay = optimisticRetryMaxDelay
	}
	half := delay / 2
	return half + time.Duration(mathrand.Int64N(int64(half)+1))
}

func waitBeforeOptimisticRetry(ctx context.Context, attempt int) error {
	timer := time.NewTimer(optimisticRetryDelay(attempt))
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
