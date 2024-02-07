package expression

import (
	"context"
	"errors"

	orch "github.com/xleshka/distributedcalc/backend/internal/orchestrator"
)

type Repository interface {
	Add(ctx context.Context, expression orch.Expression) (id int, err error)
	GetAllExpressions(ctx context.Context) ([]orch.Expression, error)
	GetExpressionById(ctx context.Context, id string) (orch.Expression, error)
	GetOneExpression(ctx context.Context, expression string) (orch.Expression, bool)
	SetExpression(ctx context.Context, exprassion orch.Expression) error
	GetAllOperations(ctx context.Context) ([]orch.Operation, error)
	SetExecutionTime(ctx context.Context, name rune, timeInSeconds int)
	Delete(ctx context.Context, id string) error
	DeleteDaemon(ctx context.Context, id string) error
}

var (
	ErrExpressionExists = errors.New("expression exist")
)
