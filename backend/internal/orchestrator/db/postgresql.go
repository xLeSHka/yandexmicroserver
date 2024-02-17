package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/xleshka/distributedcalc/backend/internal/agent"
	"github.com/xleshka/distributedcalc/backend/internal/application/cache"
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

func (r *repository) Add(ctx context.Context, expression *orch.Expression, cache *cache.Cache) error {
	q := `
		INSERT INTO public.expressions 
			(expression,expression_status,created_at) 
		VALUES ($1,$2,$3) 
		RETURNING id
	`

	r.logger.Info(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	if err := r.client.QueryRow(ctx, q, expression.Expression,
		expression.Status, expression.CreatedTime).Scan(&expression.Id); err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf("sql error:%s, Detail: %s Where: %s, Code: %s, SQLState: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			r.logger.Error(newErr.Error())
			return newErr
		}
		r.logger.Error(err.Error())
		return err
	}
	cache.Set(expression.Id, *expression)
	return nil
}

func (r *repository) GetAllExpressions(ctx context.Context, cache *cache.Cache) ([]orch.Expression, error) {
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
		cache.Set(expr.Id, expr)
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

func (r *repository) SetExpression(ctx context.Context, expression *orch.Expression, cache *cache.Cache) error {
	q := `
	UPDATE public.exprassions SET (expression,expression_status,created_at) 
	= ($1,$2,$3) WHERE id = $4;
		
	`

	r.logger.Info(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	r.client.Exec(ctx, q, expression.Expression,
		expression.Status, expression.CreatedTime, expression.Id)
	cache.Set(expression.Id, expression)
	return nil
}

func (r *repository) GetAllOperations(ctx context.Context, cache *cache.Cache) ([]orch.Operation, error) {
	q := `
		SELECT operation, execution_time_by_milliseconds 
		FROM public.operations
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
		cache.Set(operation.Operation, operation)
	}
	if err := rows.Err(); err != nil {
		r.logger.Error(err.Error())
		return nil, err
	}
	return operations, nil
}

func (r *repository) AddOperation(ctx context.Context, operation *orch.Operation, cache *cache.Cache) error {
	q := `
	INSERT INTO public.operations (operation,execution_time_by_milliseconds) 
		VALUES ($1,$2)
	`

	r.logger.Info("SQL Query: %s", formatQuery(q))

	r.client.QueryRow(ctx, q, operation.Operation, operation.ExecutionTimeByMilliseconds)
	cache.Set(operation.Operation, *operation)
	return nil
}

func (r *repository) SetExecutionTime(ctx context.Context, operaion *orch.Operation, cache *cache.Cache) error {
	q := `
	UPDATE
	public.operations SET execution_time_by_milliseconds
	= $1 WHERE operation = $2
	`

	r.logger.Info(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	r.client.Exec(ctx, q, operaion.ExecutionTimeByMilliseconds, operaion.Operation)
	cache.Set(operaion.Operation, *operaion)
	return nil
}

func (r *repository) GetAllAgents(ctx context.Context, cache *cache.Cache) ([]agent.Agent, error) {
	q := `SELECT id,status_code FROM public.agents`

	r.logger.Info(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		r.logger.Error(err.Error())
		return nil, err
	}
	agents := make([]agent.Agent, 0)
	for rows.Next() {
		var agent agent.Agent
		err := rows.Scan(&agent.ID, &agent.Status)
		if err != nil {
			r.logger.Error(err.Error())
			return nil, err
		}
		cache.Set(agent.ID, agent)
		agents = append(agents, agent)
	}
	if err := rows.Err(); err != nil {
		r.logger.Info(err.Error())
		return nil, err
	}
	return agents, nil
}
func (r *repository) AddAgent(ctx context.Context, agent *agent.Agent, cache *cache.Cache) error {
	q := `
	INSERT INTO public.agents (status_code) 
	VALUES ($1) RETURNING id
	`
	r.logger.Info("SQL Query: %s", formatQuery(q))
	if err := r.client.QueryRow(ctx, q, agent.Status).Scan(&agent.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf("sql error:%s, Detail: %s Where: %s, Code: %s, SQLState: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			r.logger.Error(newErr.Error())
			return newErr
		}
		r.logger.Error(err.Error())
		return err

	}
	cache.Set(agent.ID, *agent)
	return nil
}
func (r *repository) SetAgent(ctx context.Context, id, status string, cache *cache.Cache) error {
	q := `UPDATE
	public.agents SET status
	= $1 WHERE id = $2`

	r.logger.Info(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	r.client.Exec(ctx, q, id)
	cache.Set(id, agent.Agent{ID: id, Status: status})
	return nil
}

func (r *repository) DeleteAgent(ctx context.Context, id string, cache *cache.Cache) error {
	q := `DELETE FROM public.agents 
		WHERE id = $1 
	`
	r.logger.Info(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	row := r.client.QueryRow(ctx, q, id)
	if row != nil {
		r.logger.Info(fmt.Sprintf("delete %s", row))

	}
	delete(cache.Data, id)
	return nil
}
func NewRepository(client postgresql.Client, log *slog.Logger) orch.Repository {
	return &repository{
		client: client,
		logger: log,
	}
}
