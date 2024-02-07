package postgresql

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xleshka/distributedcalc/backend/internal/config"
	"github.com/xleshka/distributedcalc/backend/internal/orchestrator"
	repeatable "github.com/xleshka/distributedcalc/backend/pkg/utils"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}
type Expression struct {
	orchestrator.Expression
}
type Operation struct {
	orchestrator.Operation
}

func NewClient(ctx context.Context, sc config.StorageConfig, log *slog.Logger) (pool *pgxpool.Pool, err error) {
	const fn = "internal.postgresql.newclient"

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		sc.Username, sc.Password, sc.Host, sc.Port, sc.Database)

	err = repeatable.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		pool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			return err
		}
		return nil
	}, sc.MaxAttempts, 5*time.Second)
	if err != nil {
		log.With(
			slog.String("%w", fn),
			slog.String("%s", err.Error()),
		)

		return nil, err
	}
	return pool, err
}
