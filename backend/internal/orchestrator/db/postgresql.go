package orchestrator

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"

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

func (r *repository) Add(ctx context.Context, expression *orch.Expression, log *slog.Logger) error {
	q := `
		INSERT INTO expressions 
			(expression,expression_status,created_at,completed_at,execution_time_by_milliseconds) 
		VALUES ($1,$2,$3,$4,$5,$6) 
		RETURN id
	`
	executTime := make(map[string]int)
	for _, r := range expression.Expression {
		switch r {
		case '+':
			executTime["+"]++
		case '-':
			executTime["-"]++
		case '/':
			executTime["/"]++
		case '*':
			executTime["*"]++
		}
	}

	res := func(mp map[string]int) int {
		tempRes := 0
		for _, num := range mp {
			tempRes += num
		}
		return tempRes
	}(executTime)

	expression.ExecutionTimeByMilliseconds = res
	expression.CompletedTime = expression.CreatedTime.Add(time.Duration(res) * time.Millisecond)
	ctx = context.WithValue(ctx, "expression", expression)

	r.logger.Info(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	if err := r.client.QueryRow(ctx, q, expression.Expression,
		expression.Status, expression.CreatedTime, expression.CompletedTime,
		expression.ExecutionTimeByMilliseconds).Scan(&expression.Id); err != nil {
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

func (r *repository) GetAllExpressions(ctx context.Context, log *slog.Logger) ([]orch.Expression, error) {
	q := `
	SELECT id,expression,expression_status,created_at,completed_at,execution_time_by_milliseconds 
		FROM expressions
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
			&expr.CreatedTime, &expr.CompletedTime, &expr.ExecutionTimeByMilliseconds)
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

func (r *repository) GetExpressionById(ctx context.Context, id string, log *slog.Logger) (orch.Expression, error) {
	q := `
	SELECT id,expression,expression_status,created_at,completed_at,execution_time_by_milliseconds 
		FROM expressions WHERE id = $1
	`

	r.logger.Info(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var expr orch.Expression

	err := r.client.QueryRow(ctx, q, id).Scan(&expr.Id, &expr.Expression, &expr.Status,
		&expr.CreatedTime, &expr.CompletedTime, &expr.ExecutionTimeByMilliseconds)
	if err != nil {
		r.logger.Error(err.Error())
		return orch.Expression{}, err
	}
	return expr, nil
}

func (r *repository) CheckExists(ctx context.Context, expression string, log *slog.Logger) (orch.Expression, bool) {
	q := `
		SELECT id,expression,expression_status,created_at,completed_at,execution_time_by_milliseconds 
			FROM expressions WHERE expression = $1
	`

	r.logger.Info(fmt.Sprintf("SQL Suery: %s", formatQuery(q)))

	var expr orch.Expression

	err := r.client.QueryRow(ctx, q, expression).Scan(&expr.Id, &expr.Expression, &expr.Status,
		&expr.CreatedTime, &expr.CompletedTime, &expr.ExecutionTimeByMilliseconds)

	if err == sql.ErrNoRows {
		return orch.Expression{}, false
	} else if err != nil {
		r.logger.Error(err.Error())
		return orch.Expression{}, false
	}
	return expr, true

}

func (r *repository) SetExpression(ctx context.Context, expression orch.Expression, log *slog.Logger) error {
	q := `
		INSERT INTO exprassions (expression,expression_status,created_at,completed_at,execution_time_by_milliseconds) 
		VALUES ($1,$2,$3,$4,$5,$6) WHERE id = $7
	`

	r.logger.Info(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	r.client.QueryRow(ctx, q, expression.Expression,
		expression.Status, expression.CreatedTime, expression.CompletedTime,
		expression.ExecutionTimeByMilliseconds, expression.Id)
}

func NewRepository(client postgresql.Client, log *slog.Logger) orch.Repository {
	return &repository{
		client: client,
		logger: log,
	}
}
