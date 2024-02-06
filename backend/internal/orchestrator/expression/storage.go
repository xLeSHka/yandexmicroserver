package expression

import (
	"context"

	orch "github.com/xleshka/distributedcalc/backend/internal/orchestrator"
)

type Repository interface {
	Add(ctx context.Context, expression orch.Expression) (string, error)
	GetAllExpressions(ctx context.Context) ([]orch.Expression, error)
	GetOneExpression(ctx context.Context, id string) (orch.Expression, error)
	SetExpression(ctx context.Context, exprassion orch.Expression) error
	GetAllOperations(ctx context.Context) ([]orch.Operation, error)
	SetExecutionTime(ctx context.Context, name rune, timeInSeconds int)
	Delete(ctx context.Context, id string) error
}
