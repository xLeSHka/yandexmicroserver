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

func GetExpressionHandler(ctx context.Context, log *slog.Logger, rep orch.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "bad post expression method type", http.StatusBadRequest)
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

func PostExpressionsHandler(ctx context.Context, log *slog.Logger, rep orch.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			log.Error("bad get expressions method type")
			http.Error(w, "bad get expressions method type", http.StatusBadRequest)
			return
		}
		expressions, err := app.AllExpressions(ctx, log, rep)
		if err != nil {
			log.Error(fmt.Sprintf("failed bd get expressions: %v", err))
			http.Error(w, fmt.Sprintf("failed bd get expressions: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "Application/json")
		if err := json.NewEncoder(w).Encode(expressions); err != nil {
			log.Error(fmt.Sprintf("failed encode expressions: %v", err))
			http.Error(w, fmt.Sprintf("failed encode expressions: %v", err), http.StatusInternalServerError)
		}
	}
}

func PostOperationsHandler(ctx context.Context, log *slog.Logger, rep orch.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			log.Error("bad get operations method type")
			http.Error(w, "ad method get operations", http.StatusBadRequest)
			return
		}
		operations, err := app.AllOperations(ctx, log, rep)
		if err != nil {
			log.Error(fmt.Sprintf("failed bd get operaions: %v", err))
			http.Error(w, fmt.Sprintf("failed bd get operaions: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "Application/json")
		if err := json.NewEncoder(w).Encode(operations); err != nil {
			log.Error(fmt.Sprintf("failed encode operations: %v", err))
			http.Error(w, fmt.Sprintf("failed encode operations: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

func GetOperationHandler(ctx context.Context, log *slog.Logger, rep orch.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Error("bad post expression method type")
			http.Error(w, "bad post expression method type", http.StatusBadRequest)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error("failed read req body", sl.Err(err))
			http.Error(w, "failed read req body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()
		var oper orch.Operation
		err = json.Unmarshal(body, &oper)
		if err != nil {
			log.Error("failed unmarshal req body", sl.Err(err))
			http.Error(w, "failed unmarshal req body", http.StatusInternalServerError)
			return
		}
		fmt.Println(oper)
		err = app.SetOperation(ctx, log, oper, rep)
		if err != nil {
			log.Error("failed set operation to bd", sl.Err(err))
			http.Error(w, "failed set operation to bd", http.StatusInternalServerError)
			return
		}
		log.Info("req body decoded", slog.Any("req", oper))

	}
}

func PostAgentsHandler(ctx context.Context, log *slog.Logger, rep orch.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			log.Error("bad get agents method type")
			http.Error(w, fmt.Sprintf("bad get agents method type"), http.StatusBadRequest)
		}
		agents, err := app.AllAgents(ctx, log, rep)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, fmt.Sprintf("failed get agents: %v", err), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "Application/json")
		if err := json.NewEncoder(w).Encode(agents); err != nil {
			log.Error(err.Error())
			http.Error(w, fmt.Sprintf("failed encode agents to json: %v", err), http.StatusInternalServerError)
		}
	}
}
