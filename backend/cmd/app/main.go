package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	server "github.com/xleshka/distributedcalc/backend/http-server/handler/add"
	"github.com/xleshka/distributedcalc/backend/http-server/middleware"
	"github.com/xleshka/distributedcalc/backend/internal/config"
	"github.com/xleshka/distributedcalc/backend/internal/orchestrator/db"
	"github.com/xleshka/distributedcalc/backend/pkg/postgresql"
)

func main() {
	logg := setupLogger()
	cfg := config.GetConfig(*logg)
	str := cfg.HTTPServer.Host
	logg.Info(str)
	logg.Debug("logger debg mode enabled")
	ctx := context.Background()
	client := &http.Client{}

	postgreSQLClient, err := postgresql.NewClient(context.TODO(), cfg.StorageConfig, logg)
	if err != nil {
		log.Fatal(err)
	}
	repository := db.NewRepository(postgreSQLClient, logg)
	agCount, err := strconv.Atoi(cfg.HTTPServer.AgentCount)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/initialize", middleware.RecoveryMiddleware(middleware.LoggingMiddleware(server.AgentsInitializeHandler(agCount, ctx, logg, repository, client), logg)))
	mux.Handle("/add", middleware.RecoveryMiddleware(middleware.LoggingMiddleware(server.GetExpressionHandler(ctx, logg, repository, client), logg)))
	mux.Handle("/expressions", middleware.RecoveryMiddleware(middleware.LoggingMiddleware(server.PostExpressionsHandler(ctx, logg, repository, client), logg)))
	mux.Handle("/operations", middleware.RecoveryMiddleware(middleware.LoggingMiddleware(server.PostOperationsHandler(ctx, logg, repository), logg)))
	mux.Handle("/agents", middleware.RecoveryMiddleware(middleware.LoggingMiddleware(server.PostAgentsHandler(ctx, logg, repository), logg)))
	mux.Handle("/setOperation", middleware.RecoveryMiddleware(middleware.LoggingMiddleware(server.GetOperationHandler(ctx, logg, repository), logg)))
	mux.Handle("/setAgentStatus", middleware.RecoveryMiddleware(middleware.LoggingMiddleware(server.GetAgentStatusHandler(ctx, logg, repository), logg)))
	mux.Handle("/addAgent", middleware.RecoveryMiddleware(middleware.LoggingMiddleware(server.GetAddAgentHandler(ctx, logg, repository), logg)))

	log.Fatal(http.ListenAndServe(":"+cfg.HTTPServer.OrchPort, mux))
}
func setupLogger() *slog.Logger {
	logg := slog.New(slog.NewJSONHandler(
		os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return logg
}
