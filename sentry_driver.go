package entgosentry

import (
	"context"
	"database/sql"
	"fmt"

	"entgo.io/ent/dialect"
	"github.com/getsentry/sentry-go"
)

const (
	op string = "db"
)

type SentryDriver struct {
	dialect.Driver
}

// Trace возвращает Driver для трассировки запросов к базе данныхgo
func Trace(d dialect.Driver) dialect.Driver {
	drv := &SentryDriver{d}
	return drv
}

func setSpan(ctx context.Context, query string) *sentry.Span {
	span := sentry.StartSpan(ctx, op)
	span.Description = query
	return span
}

func (d *SentryDriver) Exec(ctx context.Context, query string, args, v any) error {
	span := setSpan(ctx, query)
	defer span.Finish()
	return d.Driver.Exec(ctx, query, args, v)
}

func (d *SentryDriver) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	drv, ok := d.Driver.(interface {
		ExecContext(context.Context, string, ...any) (sql.Result, error)
	})
	if !ok {
		fmt.Println("Warning: Driver.ExecContext is not supported")
	}
	span := setSpan(ctx, query)
	defer span.Finish()
	return drv.ExecContext(ctx, query, args...)
}

func (d *SentryDriver) Query(ctx context.Context, query string, args, v any) error {
	span := setSpan(ctx, query)
	defer span.Finish()
	return d.Driver.Query(ctx, query, args, v)
}

func (d *SentryDriver) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	drv, ok := d.Driver.(interface {
		QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	})
	if !ok {
		fmt.Println("Warning: Driver.QueryContext is not supported")
	}
	span := setSpan(ctx, query)
	defer span.Finish()
	return drv.QueryContext(ctx, query, args...)
}
