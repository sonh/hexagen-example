package database

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// queryer is an interface for Query
type queryer interface {
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
}

// execer is an interface for Exec
type execer interface {
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

type QueryExecer interface {
	queryer
	execer
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
}

func GetFieldNames(e interface{}) []string {

	return nil
}

func GetScanFields(e interface{}, reqlist []string) []interface{} {
	return nil
}
