package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	app "github.com/xleshka/distributedcalc/backend/internal/application/app"
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

func AddExpressionHandler(ctx context.Context, log *slog.Logger, rep orch.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "invalid method", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil && err != io.EOF {
			log.Error("failed decode request body", sl.Err(err))
			http.Error(w, fmt.Sprintf("body read error: %v", err), http.StatusBadRequest)
			return
		}
		expr := string(body)
		if !app.ValidExpression(expr) {
			log.Error("invalid exprassion")
			http.Error(w, fmt.Sprintf("invalid exprassion: %v", err), http.StatusBadRequest)
			return
		}

		expression := orch.Expression{}
		expression.Expression = expr
		t := time.Now()
		expression.CreatedTime = t
		expression.Status = "wait"

		expression, err = app.AddExpression(ctx, expression, log, rep)
		if err != nil {
			log.Error("failed add expression to bd")
			http.Error(w, fmt.Sprintf("failed add expression to bd: %v", err), http.StatusInternalServerError)
			return
		}
		log.Info("req body decoded", slog.Any("req", expression))
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(expression); err != nil {
			log.Error("json encode error")
			http.Error(w, fmt.Sprintf("json encode error: %v", err), http.StatusInternalServerError)
			return
		}

	}
}

func PostExpression(ctx context.Context, log *slog.Logger, rep orch.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Error("failed post expression")
			http.Error(w, "failed post expression", http.StatusBadRequest)
			return
		}
		expressions, err := rep.GetAllExpressions(ctx)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, fmt.Sprintf("failed get expressions: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "Application/json")
		if err := json.NewEncoder(w).Encode(expressions); err != nil {
			log.Error(err.Error())
			http.Error(w, fmt.Sprintf("failed encode expressions: %v", err), http.StatusInternalServerError)
		}
	}
}
