package application

import (
	"context"
	"log/slog"
	"unicode"

	orch "github.com/xleshka/distributedcalc/backend/internal/orchestrator"
)

func ValidExpression(expression string) bool {
	for _, r := range expression {
		/*проверяем число ли символ, если нет, то проверяем это один из валидных символов или нет*/
		if !unicode.Is(unicode.Digit, r) {
			if !(r >= 0x08 && r <= 0x0B || r == 0x0D || r == 0x0F) {
				return false
			}
		}
	}
	return true
}
func AddExpression(ctx context.Context, express orch.Expression, log *slog.Logger, rep orch.Repository) (orch.Expression, error) {

	expression, exist := rep.CheckExists(ctx, express.Expression)
	if exist {
		return expression, nil
	}
	err := rep.Add(ctx, &express)
	if err != nil {
		log.Error("failed add expression to bd")
		return orch.Expression{}, err
	}
	return express, nil
}
