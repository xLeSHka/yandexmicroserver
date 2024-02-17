package orchestrator

import (
	"context"
	"errors"

	"github.com/xleshka/distributedcalc/backend/internal/agent"
	"github.com/xleshka/distributedcalc/backend/internal/application/cache"
)

type Repository interface {
	Add(ctx context.Context, expression *Expression, cache *cache.Cache) error
	GetAllExpressions(ctx context.Context, cache *cache.Cache) ([]Expression, error)
	GetExpressionById(ctx context.Context, id string) (Expression, error)
	CheckExists(ctx context.Context, expression string) (Expression, bool)
	SetExpression(ctx context.Context, expression *Expression, cache *cache.Cache) error
	GetAllOperations(ctx context.Context, cache *cache.Cache) ([]Operation, error)
	AddOperation(ctx context.Context, operaion *Operation, cache *cache.Cache) error
	SetExecutionTime(ctx context.Context, operaion *Operation, cache *cache.Cache) error
	GetAllAgents(ctx context.Context, cache *cache.Cache) ([]agent.Agent, error)
	SetAgent(ctx context.Context, id, status string, cache *cache.Cache) error
	DeleteAgent(ctx context.Context, id string, cache *cache.Cache) error
}

var (
	ErrExpressionExists = errors.New("expression exist")
)
