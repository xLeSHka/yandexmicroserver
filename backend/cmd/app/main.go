package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5/middleware"
	server "github.com/xleshka/distributedcalc/backend/http-server/handler/add"
	app "github.com/xleshka/distributedcalc/backend/internal/application/app"
	"github.com/xleshka/distributedcalc/backend/internal/config"
	orchestrator "github.com/xleshka/distributedcalc/backend/internal/orchestrator/db"
	"github.com/xleshka/distributedcalc/backend/pkg/postgresql"
)

func main() {
	logg := setupLogger()
	cfg := config.GetConfig(*logg)
	str := cfg.HTTPServer.Host
	logg.Info(str)
	logg.Debug("logger debg mode enabled")
	ctx := context.Background()

	postgreSQLClient, err := postgresql.NewClient(context.TODO(), cfg.StorageConfig, logg)
	if err != nil {

	}
	repository := orchestrator.NewRepository(postgreSQLClient, logg)
	app.Initialize()
	_, err = app.AllExpressions(ctx, logg, repository)
	if err != nil {
		log.Fatalf("%v", err)
	}

	_, err = app.AllOperations(ctx, logg, repository)
	if err != nil {
		log.Fatalf("%v", err)
	}
	_, err = app.AllAgents(ctx, logg, repository)
	if err != nil {
		log.Fatalf("%v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/add", middleware.Recoverer(middleware.Logger(server.GetExpressionHandler(ctx, logg, repository))))
	mux.Handle("/", middleware.Recoverer(middleware.Logger(server.PostExpressionsHandler(ctx, logg, repository))))
	mux.Handle("/operations", middleware.Recoverer(middleware.Logger(server.PostOperationsHandler(ctx, logg, repository))))
	mux.Handle("/agents", middleware.Recoverer(middleware.Logger(server.PostAgentsHandler(ctx, logg, repository))))
	mux.Handle("/setOperation", middleware.Recoverer(middleware.Logger(server.GetOperationHandler(ctx, logg, repository))))

	log.Fatal(http.ListenAndServe(":"+cfg.HTTPServer.OrchPort, mux))
}
func setupLogger() *slog.Logger {
	logg := slog.New(slog.NewJSONHandler(
		os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return logg
}
