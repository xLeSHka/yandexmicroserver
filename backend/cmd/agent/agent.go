package main

import (
	"log/slog"
	"os"

	"github.com/xleshka/distributedcalc/backend/internal/config"
)

func main() {
	logg := setupLogger()
	cfg := config.GetConfig(*logg)
	str := cfg.HTTPServer.Host
	logg.Info(str)
	logg.Debug("logger debg mode enabled")
	// ctx := context.Background()

	// mux := http.NewServeMux()
	// mux.Handle("/add", middleware.Recoverer(middleware.Logger(server.GetExpressionHandler(ctx, logg, repository))))
	// mux.Handle("/", middleware.Recoverer(middleware.Logger(server.PostExpressionsHandler(ctx, logg, repository))))
	// mux.Handle("/operations", middleware.Recoverer(middleware.Logger(server.PostOperationsHandler(ctx, logg, repository))))
	// mux.Handle("/agents", middleware.Recoverer(middleware.Logger(server.PostAgentsHandler(ctx, logg, repository))))
	// log.Fatal(http.ListenAndServe(":"+cfg.HTTPServer.AgentPort, mux))
}
func setupLogger() *slog.Logger {
	logg := slog.New(slog.NewJSONHandler(
		os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return logg
}
