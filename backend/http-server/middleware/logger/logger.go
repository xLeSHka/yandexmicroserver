package logger

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With(slog.String("component", "middleware/logger"))

		log.Info("logger middleware enabled")

		/*обработчик*/
		fn := func(w http.ResponseWriter, r *http.Request) {
			/*собираем исходную информацию о запросе*/
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_agent", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)
			/* создаем обертку для полученния сведений об ответе*/
			wrw := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			/*момент получения запроса*/
			t1 := time.Now()
			defer func() {
				entry.Info("request completed",
					slog.Int("status", wrw.Status()),
					slog.Int("bytes", wrw.BytesWritten()),
					slog.Time("request_time", t1),
				)
			}()
			next.ServeHTTP(wrw, r)
		}
		return http.HandlerFunc(fn)
	}

}
