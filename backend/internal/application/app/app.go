package application

import (
	"context"
	"log/slog"
	"unicode"

	orch "github.com/xleshka/distributedcalc/backend/internal/orchestrator"
	expr "github.com/xleshka/distributedcalc/backend/internal/orchestrator/expression"
)

func ValidExpression(expression string) bool {
	for _, r := range expression {
		if !unicode.Is(unicode.Digit, r) {
			/*проверяем наличие любых символов кроме ( ) + - / *  */
			if !(r >= 0x08 && r <= 0x0B || r == 0x0D || r == 0x0F) {
				return false
			}
		}
	}
	return true
}
func AddExpression(ctx context.Context, express orch.Expression, log *slog.Logger) {
	expression, exist := ctx.Value("cache").data[]
	if exist {
		/*post expression*/
	}
	id, err := expr.Add(ctx, express)
	if err != nil {
		log.With(
			slog.String("failed add expression to BD"),
		)
		return 
	}

}
