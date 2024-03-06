package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/xleshka/distributedcalc/backend/internal/agent"
	app "github.com/xleshka/distributedcalc/backend/internal/application/app"
	orch "github.com/xleshka/distributedcalc/backend/internal/orchestrator"
)

func GetExpressionHandler(ctx context.Context, log *slog.Logger, rep orch.Repository, client *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// if r.Method != http.MethodPost {
		// 	http.Error(w, "bad post expression method type", http.StatusBadRequest)
		// 	return
		// }

		// Устанавливаем заголовки CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")              // Разрешить запросы от всех доменов
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS") // Разрешить GET и OPTIONS запросы
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")  // Разрешить заголовок Content-Type

		// Проверяем метод запроса
		if r.Method == "OPTIONS" {
			// Отправляем пустой ответ для предварительного запроса OPTIONS
			w.WriteHeader(http.StatusOK)
		}

		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil && err != io.EOF {
			log.Error("failed decode request body", err)
			http.Error(w, fmt.Sprintf("body read error: %v", err), http.StatusBadRequest)
			return
		}
		expr := string(body)
		expr = expr[1 : len([]rune(expr))-1]
		log.Info(expr)
		if !app.ValidExpression(expr) {
			log.Error("invalid exprassion")
			http.Error(w, fmt.Sprintf("invalid exprassion: %v", err), http.StatusBadRequest)
			return
		}
		//полный код можно посмотреть на gitHub
		expression := &orch.Expression{}
		expression.Expression = expr
		t := time.Now()
		expression.CreatedTime = t
		expression.CompletedTime = t
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

func PostExpressionsHandler(ctx context.Context, log *slog.Logger, rep orch.Repository, client *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// if r.Method != http.MethodGet {
		// 	log.Error("bad get expressions method type")
		// 	http.Error(w, "bad get expressions method type", http.StatusBadRequest)
		// 	return
		// }
		// Устанавливаем заголовки CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")             // Разрешить запросы от всех доменов
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS") // Разрешить GET и OPTIONS запросы
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // Разрешить заголовок Content-Type

		// Проверяем метод запроса
		if r.Method == "OPTIONS" {
			// Отправляем пустой ответ для предварительного запроса OPTIONS
			w.WriteHeader(http.StatusOK)
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
			return
		}
	}
}

func PostOperationsHandler(ctx context.Context, log *slog.Logger, rep orch.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// if r.Method != http.MethodGet {
		// 	log.Error("bad get operations method type")
		// 	http.Error(w, "ad method get operations", http.StatusBadRequest)
		// 	return
		// }
		// Устанавливаем заголовки CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")             // Разрешить запросы от всех доменов
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS") // Разрешить GET и OPTIONS запросы
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // Разрешить заголовок Content-Type

		// Проверяем метод запроса
		if r.Method == "OPTIONS" {
			// Отправляем пустой ответ для предварительного запроса OPTIONS
			w.WriteHeader(http.StatusOK)
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
		// if r.Method != http.MethodPost {
		// 	log.Error("bad post expression method type")
		// 	http.Error(w, "bad post expression method type", http.StatusBadRequest)
		// 	return
		// }

		// Устанавливаем заголовки CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")              // Разрешить запросы от всех доменов
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS") // Разрешить GET и OPTIONS запросы
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")  // Разрешить заголовок Content-Type

		// Проверяем метод запроса
		if r.Method == "OPTIONS" {
			// Отправляем пустой ответ для предварительного запроса OPTIONS
			w.WriteHeader(http.StatusOK)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error("failed read req body", err)
			http.Error(w, "failed read req body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()
		var oper orch.Operation
		err = json.Unmarshal(body, &oper)
		if err != nil {
			log.Error("failed unmarshal req body", err)
			http.Error(w, "failed unmarshal req body", http.StatusInternalServerError)
			return
		}
		fmt.Println(oper)
		err = app.SetOperation(ctx, log, oper, rep)
		if err != nil {
			log.Error("failed set operation to bd", err)
			http.Error(w, "failed set operation to bd", http.StatusInternalServerError)
			return
		}
		log.Info("req body decoded", slog.Any("req", oper))

	}
}

// func GetSubExprassion(ctx context.Context, log *slog.Logger, ag agent.Agent, heartBeatUrl string, client *http.Client) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method != http.MethodPost {
// 			log.Error("bad post subExprassion method type")
// 			http.Error(w, "bad post subExprassion method type", http.StatusInternalServerError)
// 			return
// 		}
// 		body, err := io.ReadAll(r.Body)
// 		defer r.Body.Close()
// 		if err != nil {
// 			log.Error("failed read req body")
// 			http.Error(w, "failed read req body", http.StatusInternalServerError)
// 			return
// 		}

// 		var expr expparser.SubExpression

//			err = json.Unmarshal(body, &expr)
//			if err != nil {
//				log.Error(err.Error())
//				return
//			}
//			errCh := make(chan struct{})
//			var res string
//			ag.Status = "Busy"
//			app.AgentHeartBeat(ctx, log, ag, heartBeatUrl, client, errCh)
//			res = app.Cacl(expr.Expression, log)
//			ag.Status = "Ok"
//			app.AgentHeartBeat(ctx, log, ag, heartBeatUrl, client, errCh)
//			close(errCh)
//			log.Info(fmt.Sprintf("Calc sub expression %s", res))
//			w.Write([]byte(res))
//		}
//	}
func GetAddAgentHandler(ctx context.Context, log *slog.Logger, rep orch.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Error("bad post add agent method type")
			http.Error(w, "bad post add agent method type", http.StatusInternalServerError)
			return
		}

		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			log.Error("failed read req body: %v", err)
			http.Error(w, "failed read req body", http.StatusInternalServerError)
			return
		}
		var ag agent.Agent
		err = json.Unmarshal(body, &ag)
		if err != nil {
			log.Error("failed unmarshal req body: %v", err)
			http.Error(w, "failed unmarshal req body", http.StatusInternalServerError)
			return
		}
		id, err := app.AddAgent(ctx, log, ag, rep)
		if err != nil {
			log.Error("%v", err)
			http.Error(w, "failed add agent", http.StatusInternalServerError)
			return
		}
		w.Write([]byte(id))
	}
}
func GetAgentStatusHandler(ctx context.Context, log *slog.Logger, rep orch.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Error("bad post agent status method type")
			http.Error(w, "bad post agent status method type", http.StatusInternalServerError)
			return
		}
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			log.Error("failed read req body %v", err)
			http.Error(w, "failed read req body", http.StatusInternalServerError)
			return
		}
		var agent agent.Agent
		err = json.Unmarshal(body, &agent)
		if err != nil {
			log.Error("failed unmarshal body")
			http.Error(w, "failed unmarshal body", http.StatusInternalServerError)
			return
		}
		agent.LastHearBeat = time.Now()
		err = app.SetAgent(ctx, log, agent, rep)
		if err != nil {
			log.Error("failed set agent  status: %v", err)
			http.Error(w, "failed set agent  status", http.StatusInternalServerError)
			return
		}
		log.Info("agent status decoded", slog.Any("req", agent))
	}
}
func PostAgentsHandler(ctx context.Context, log *slog.Logger, rep orch.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// if r.Method != http.MethodGet {
		// 	log.Error("bad get agents method type")
		// 	http.Error(w, "bad get agents method type", http.StatusBadRequest)
		// 	return
		// }
		// Устанавливаем заголовки CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")             // Разрешить запросы от всех доменов
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS") // Разрешить GET и OPTIONS запросы
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // Разрешить заголовок Content-Type

		// Проверяем метод запроса
		if r.Method == "OPTIONS" {
			// Отправляем пустой ответ для предварительного запроса OPTIONS
			w.WriteHeader(http.StatusOK)
		}

		agents, err := app.AllAgents(ctx, log, rep)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, fmt.Sprintf("failed get agents: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "Application/json")
		if err := json.NewEncoder(w).Encode(agents); err != nil {
			log.Error(err.Error())
			http.Error(w, fmt.Sprintf("failed encode agents to json: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

func AgentsInitializeHandler(agCount int, ctx context.Context, log *slog.Logger, rep orch.Repository, client *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			log.Error("bad init agents method type")
			http.Error(w, "bad init agents method type", http.StatusBadRequest)
			return
		}
		err := app.Initialize(agCount, ctx, log, rep, client)
		if err != nil {
			log.Error("bad init agents method type")
			http.Error(w, "bad init agents method type", http.StatusBadRequest)
			return
		}
	}
}
