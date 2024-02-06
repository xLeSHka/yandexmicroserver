package main

import (
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	server "github.com/xleshka/distributedcalc/backend/http-server/handler/add"
	"github.com/xleshka/distributedcalc/backend/http-server/middleware/logger"
	"github.com/xleshka/distributedcalc/backend/internal/config"
)

func main() {
	log := setupLogger()
	cfg := config.GetConfig(*log)
	str := cfg.Host /*просто*/
	log.Info(str)   /*чтобы конфиг использовать*/
	log.Debug("logger debg mode enabled")

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Post("/Add", server.AddExpressionHandler(log))
}
func setupLogger() *slog.Logger {
	log := slog.New(slog.NewJSONHandler(
		os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return log
}
