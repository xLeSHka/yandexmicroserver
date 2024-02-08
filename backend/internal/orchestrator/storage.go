package orchestrator

import (
	"context"
	"errors"
)

type Repository interface {
	Add(ctx context.Context, expression *Expression) error
	GetAllExpressions(ctx context.Context) ([]Expression, error)
	GetExpressionById(ctx context.Context, id string) (Expression, error)
	CheckExists(ctx context.Context, expression string) (Expression, bool)
	SetExpression(ctx context.Context, expression Expression) error
	GetAllOperations(ctx context.Context) ([]Operation, error)
	SetExecutionTime(ctx context.Context, operation string, timeInSeconds int)
	// SetDaemon(ctx context.Context) error
	// DeleteDaemon(ctx context.Context, id string) error
}

var (
	ErrExpressionExists = errors.New("expression exist")
)
