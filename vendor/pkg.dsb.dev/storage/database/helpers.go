// Package database contains utilities for interacting with databases.
package database

import (
	"context"
	"database/sql"
	"time"

	"pkg.dsb.dev/multierror"

	"github.com/jackc/pgtype"
)

// FromTextArray converts an instance of pgtype.TextArray to a string array.
func FromTextArray(ta pgtype.TextArray) []string {
	out := make([]string, len(ta.Elements))
	for i, elem := range ta.Elements {
		out[i] = elem.String
	}
	return out
}

// ToTextArray converts a given string slice to a pgtype.TextArray instance.
// Will panic if Set returns an error.
func ToTextArray(arr []string) pgtype.TextArray {
	tags := pgtype.TextArray{}
	if err := tags.Set(arr); err != nil {
		panic(err)
	}
	return tags
}

// ToInterval converts a given time.Duration to a pgtype.Interval instance.
// Will panic if Set returns an error.
func ToInterval(dur time.Duration) pgtype.Interval {
	interval := pgtype.Interval{}
	if err := interval.Set(dur); err != nil {
		panic(err)
	}
	return interval
}

// FromInterval converts a given pgtype.Interval instance into a time.Duration
// instance. It assumes the interval is set using microseconds, as that is how
// `ToInterval` will create them.
func FromInterval(interval pgtype.Interval) time.Duration {
	const conversion = 1000
	return time.Duration(interval.Microseconds * conversion)
}

// WithinTransaction invokes 'cb' within an SQL transaction. If the callback returns an error, the transaction
// is rolled back.
func WithinTransaction(ctx context.Context, db *sql.DB, cb func(ctx context.Context, tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := cb(ctx, tx); err != nil {
		return multierror.Append(err, tx.Rollback())
	}

	return tx.Commit()
}
