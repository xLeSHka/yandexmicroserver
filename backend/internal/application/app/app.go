package application

import (
	"context"
	"fmt"
	"log/slog"
	"time"
	"unicode"

	"github.com/xleshka/distributedcalc/backend/internal/agent"
	"github.com/xleshka/distributedcalc/backend/internal/application/cache"
	orch "github.com/xleshka/distributedcalc/backend/internal/orchestrator"
)

var (
	agentCache      *cache.Cache
	operaionsCache  *cache.Cache
	expressionCache *cache.Cache
)

func Initialize() {
	agentCache = cache.NewCache()
	operaionsCache = cache.NewCache()
	expressionCache = cache.NewCache()
}
func ValidExpression(expression string) bool {

	for _, r := range expression {
		/*проверяем число ли символ, если нет, то проверяем это один из валидных символов или нет*/
		switch {
		case unicode.Is(unicode.Digit, r):
			continue
		case r == ' ':
			continue
		case r == '+':
			continue
		case r == '-':
			continue
		case r == '*':
			continue
		case r == '/':
			continue
		default:
			return false
		}
	}
	return true
}
func AllExpressions(ctx context.Context, log *slog.Logger, rep orch.Repository) ([]orch.Expression, error) {
	expressions := make([]orch.Expression, 0)
	if len(expressionCache.Data) == 0 {
		rep.GetAllExpressions(ctx, expressionCache)
	}
	for _, expr := range expressionCache.Data {
		expression, _ := expr.(orch.Expression)
		expressions = append(expressions, expression)
	}
	return expressions, nil
}
func AddExpression(ctx context.Context, express orch.Expression, log *slog.Logger, rep orch.Repository) (orch.Expression, error) {

	expression, exist := rep.CheckExists(ctx, express.Expression)
	if exist {
		return expression, nil
	}
	err := rep.Add(ctx, &express, expressionCache)
	if err != nil {
		log.Error("failed add expression to bd")
		return orch.Expression{}, err
	}
	return express, nil
}

func AllOperations(ctx context.Context, log *slog.Logger, rep orch.Repository) ([]orch.Operation, error) {
	operations := make([]orch.Operation, 0)
	if len(operaionsCache.Data) == 0 {
		rep.GetAllOperations(ctx, operaionsCache)
		if len(operaionsCache.Data) < 4 {
			oprtrs := []string{
				"+", "-", "*", "/",
			}
			for _, o := range oprtrs {
				op := orch.Operation{Operation: o,
					ExecutionTimeByMilliseconds: 2000}
				err := AddOperation(ctx, log, op, rep)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	for _, operation := range operaionsCache.Data {
		op, _ := operation.(orch.Operation)
		operations = append(operations, op)
	}
	log.Info("operations: ", operations)
	return operations, nil
}
func AddOperation(ctx context.Context, log *slog.Logger, operation orch.Operation, rep orch.Repository) error {
	fmt.Println(operation)
	err := rep.AddOperation(ctx, &operation, operaionsCache)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}
func SetOperation(ctx context.Context, log *slog.Logger, operaion orch.Operation, rep orch.Repository) error {

	err := rep.SetExecutionTime(ctx, &operaion, operaionsCache)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func AllAgents(ctx context.Context, log *slog.Logger, rep orch.Repository) ([]agent.Agent, error) {
	agents := make([]agent.Agent, 0)
	if len(agentCache.Data) == 0 {
		rep.GetAllAgents(ctx, agentCache)
	}
	for _, ag := range agentCache.Data {
		a, _ := ag.(agent.Agent)
		if a.Status == "500" {
			go DeleteAgent(ctx, log, a.ID, rep)
			continue
		}
		agents = append(agents, a)
	}
	return agents, nil
}
func SetAgent(ctx context.Context, log *slog.Logger, id, status string, rep orch.Repository) error {
	err := rep.SetAgent(ctx, id, status, agentCache)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func timerStop(v time.Timer, ctx context.Context, id string, rep orch.Repository) error {
	<-v.C
	err := rep.DeleteAgent(ctx, id, agentCache)
	if err != nil {
		return err
	}
	return nil
}
func DeleteAgent(ctx context.Context, log *slog.Logger, id string, rep orch.Repository) {
	timer := time.NewTimer(60 * time.Second)

	err := timerStop(*timer, ctx, id, rep)
	if err != nil {
		log.Error(err.Error())
	}
}
