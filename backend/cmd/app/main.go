package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	server "github.com/xleshka/distributedcalc/backend/http-server/handler/add"
	"github.com/xleshka/distributedcalc/backend/http-server/middleware/logger"
	"github.com/xleshka/distributedcalc/backend/internal/application/cache"
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

	ctx := context.Background()
	cache := cache.NewCache()
	ctx = context.WithValue(ctx, "cache", cache)

	postgreSQLClient, err := postgresql.NewClient(context.TODO(), cfg.StorageConfig, logg)
	if err != nil {
		log.Fatalf("%v", err)
	}
	r := orchestrator.NewRepository(postgreSQLClient, logg)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(logger.New(logg))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Post("/Add", server.AddExpressionHandler(ctx, logg))
}
func setupLogger() *slog.Logger {
	logg := slog.New(slog.NewJSONHandler(
		os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return logg
}
