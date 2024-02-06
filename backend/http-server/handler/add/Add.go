package server

import (
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	app "github.com/xleshka/distributedcalc/backend/internal/application"
	resp "github.com/xleshka/distributedcalc/backend/internal/lib/api/response"
	"github.com/xleshka/distributedcalc/backend/internal/lib/logger/sl"
	orch "github.com/xleshka/distributedcalc/backend/internal/orchestrator"
)

type Request struct {
	ID     int `json:"id"`
	Status int `json:"status"`
}
type Response struct {
	resp.Response
	Result int `json:"result"`
}

func AddExpressionHandler(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "http.handler.add.addexpression"

		log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		body, err := io.ReadAll(r.Body)
		if err != nil && err != io.EOF {
			log.Error("failed decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("bad req", 400))
			return
		} else if len(body) == 0 {
			log.Error("invalid expression")
			render.JSON(w, r, resp.Error("bad request", 400))
			return
		}
		expr := string(body)
		if !app.ValidExpression(expr) {
			log.Error("invalid exprassion")
			render.JSON(w, r, resp.Error("bad request", 400))
			return
		}
		expression := orch.Expression{}
		expression.Expression = expr
		t := time.Now()
		expression.CreatedTime = t
		expression.Status = orch.Wait
		log.Info("req body decoded", slog.Any("req", expression))
	}
}
