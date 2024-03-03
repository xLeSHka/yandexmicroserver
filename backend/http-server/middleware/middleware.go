package middleware

import (
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func LoggingMiddleware(next http.Handler, log *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Info("logger middleware enabled")

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
	})
}
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic: %v", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
