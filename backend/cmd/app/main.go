package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	server "github.com/xleshka/distributedcalc/backend/http-server/handler/add"
	"github.com/xleshka/distributedcalc/backend/internal/config"
	orchestrator "github.com/xleshka/distributedcalc/backend/internal/orchestrator/db"
	"github.com/xleshka/distributedcalc/backend/pkg/postgresql"
)

func main() {
	logg := setupLogger()
	cfg := config.GetConfig(*logg)
	str := cfg.Host /*просто*/
	logg.Info(str)  /*чтобы конфиг использовать*/
	logg.Debug("logger debg mode enabled")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// cache := cache.NewCache()

	postgreSQLClient, err := postgresql.NewClient(context.TODO(), cfg.StorageConfig, logg)
	if err != nil {
		log.Fatalf("%v", err)
	}
	repository := orchestrator.NewRepository(postgreSQLClient, logg)

	mux := http.NewServeMux()
	mux.Handle("/add/", middleware.Recoverer(middleware.Logger(server.AddExpressionHandler(ctx, logg, repository))))
	mux.Handle("/", middleware.Recoverer(middleware.Logger(server.PostExpression(ctx, logg, repository))))
	log.Fatal(http.ListenAndServe(":8080", mux))
}
func setupLogger() *slog.Logger {
	logg := slog.New(slog.NewJSONHandler(
		os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return logg
}
