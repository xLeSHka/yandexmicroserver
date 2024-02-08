package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	orch "github.com/xleshka/distributedcalc/backend/internal/orchestrator"
	"github.com/xleshka/distributedcalc/backend/pkg/postgresql"
)

type repository struct {
	client postgresql.Client
	logger *slog.Logger
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", " "), "\n", " ")
}

func (r *repository) Add(ctx context.Context, expression *orch.Expression) error {
	q := `
		INSERT INTO public.expressions 
			(expression,expression_status,created_at) 
		VALUES ($1,$2,$3) 
		RETURNING id
	`

	r.logger.Info(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	if err := r.client.QueryRow(ctx, q, expression.Expression,
		expression.Status, expression.CreatedTime).Scan(&expression.Id); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf("sql error:%s, Detail: %s Where: %s, Code: %s, SQLState: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			r.logger.Error(newErr.Error())
			return newErr
		}
		r.logger.Error(err.Error())
		return err
	}
	return nil
}

func (r *repository) GetAllExpressions(ctx context.Context) ([]orch.Expression, error) {
	q := `
	SELECT id,expression,expression_status,created_at
		FROM public.expressions
	`

	r.logger.Info(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	expressions := make([]orch.Expression, 0)

	for rows.Next() {
		var expr orch.Expression

		err := rows.Scan(&expr.Id, &expr.Expression, &expr.Status,
			&expr.CreatedTime)
		if err != nil {
			r.logger.Error(err.Error())
			return nil, err
		}
		expressions = append(expressions, expr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return expressions, nil
}

func (r *repository) GetExpressionById(ctx context.Context, id string) (orch.Expression, error) {
	q := `
	SELECT id,expression,expression_status,created_at
		FROM public.expressions WHERE id = $1
	`

	r.logger.Info(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var expr orch.Expression

	err := r.client.QueryRow(ctx, q, id).Scan(&expr.Id, &expr.Expression, &expr.Status,
		&expr.CreatedTime)
	if err != nil {
		r.logger.Error(err.Error())
		return orch.Expression{}, err
	}
	return expr, nil
}

func (r *repository) CheckExists(ctx context.Context, expression string) (orch.Expression, bool) {
	q := ` 
		SELECT id,expression,expression_status,created_at
			FROM public.expressions WHERE expression = $1
	`
	r.logger.Info(fmt.Sprintf("SQL Suery: %s", formatQuery(q)))

	var expr orch.Expression

	err := r.client.QueryRow(ctx, q, expression).Scan(&expr.Id, &expr.Expression, &expr.Status,
		&expr.CreatedTime)

	if err == sql.ErrNoRows {
		r.logger.Info("false suc")
		return orch.Expression{}, false
	} else if err != nil {
		r.logger.Info("unsuccesful")
		r.logger.Error(err.Error())
		return orch.Expression{}, false
	}
	r.logger.Info("true suc")
	return expr, true

}

func (r *repository) SetExpression(ctx context.Context, expression orch.Expression) error {
	q := `
		INSERT INTO public.exprassions (expression,expression_status,created_at) 
		VALUES ($1,$2,$3) WHERE id = $4
	`

	r.logger.Info(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	r.client.QueryRow(ctx, q, expression.Expression,
		expression.Status, expression.CreatedTime, expression.Id)

	return nil
}

func (r *repository) GetAllOperations(ctx context.Context) ([]orch.Operation, error) {
	q := `
		SELECT operation, execution_time_by_milliseconds FROM public.operations
	`
	r.logger.Info(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		r.logger.Error(err.Error())
		return nil, err
	}

	operations := make([]orch.Operation, 0)

	for rows.Next() {
		var operation orch.Operation
		err := rows.Scan(&operation.Operation, &operation.ExecutionTimeByMilliseconds)
		if err != nil {
			r.logger.Error(err.Error())
			return nil, err
		}
		operations = append(operations, operation)
	}
	if err := rows.Err(); err != nil {
		r.logger.Error(err.Error())
		return nil, err
	}
	return operations, nil
}

func (r *repository) SetExecutionTime(ctx context.Context, operation string, timeInSeconds int) {
	q := `
		INSERT INTO 
		public.operations (execution_time_by_milliseconds) 
		VALUES ($1) WHERE operation = $2
	`

	r.logger.Info(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	r.client.QueryRow(ctx, q, timeInSeconds, operation)

}

func NewRepository(client postgresql.Client, log *slog.Logger) orch.Repository {
	return &repository{
		client: client,
		logger: log,
	}
}
