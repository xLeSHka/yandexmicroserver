package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	server "github.com/xleshka/distributedcalc/backend/http-server/handler/add"
	"github.com/xleshka/distributedcalc/backend/http-server/middleware"
	"github.com/xleshka/distributedcalc/backend/internal/agent"
	app "github.com/xleshka/distributedcalc/backend/internal/application/app"
	"github.com/xleshka/distributedcalc/backend/internal/config"
)

func main() {
	logg := setupLogger()
	cfg := config.GetConfig(*logg)
	str := cfg.HTTPServer.Host
	logg.Info(str)
	logg.Debug("logger debg mode enabled")
	agentPort, err := strconv.Atoi(cfg.AgentPort)
	if err != nil {
		log.Fatalf("invalid agent count: %v", err)
	}
	heartBeatUrl := "http://" + cfg.HTTPServer.Host + ":" + cfg.OrchPort + "/setAgentStatus"
	addAgentUrl := "http://" + cfg.HTTPServer.Host + ":" + cfg.OrchPort + "/addAgent"
	// postgreSQLClient, err := postgresql.NewClient(context.TODO(), cfg.StorageConfig, logg)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// repository := db.NewRepository(postgreSQLClient, logg)

	// mux.Handle("/", middleware.Recoverer(middleware.Logger(server.PostExpressionsHandler(ctx, logg, repository))))
	// mux.Handle("/operations", middleware.Recoverer(middleware.Logger(server.PostOperationsHandler(ctx, logg, repository))))
	// mux.Handle("/agents", middleware.Recoverer(middleware.Logger(server.PostAgentsHandler(ctx, logg, repository))))
	// for i := 0; i < agCount; i++ {
	// 	port := agentPort + i
	// 	go func(port string) {
	port := fmt.Sprintf("%d", agentPort)
	errCh := make(chan struct{})
	defer close(errCh)
	agentCtxWithCancel, cancelCtx := context.WithCancel(context.Background())
	agentAddress := "http://localhost:" + port
	client := &http.Client{}
	ag := agent.Agent{Address: agentAddress, Status: "Ok"}

	ag.ID, err = app.AddAgentReq(agentCtxWithCancel, logg, ag, addAgentUrl, client)

	if err != nil {
		Shutdown(agentCtxWithCancel, cancelCtx, logg, ag, heartBeatUrl, client, errCh)
		return
	}
	go func() {
		select {
		case <-agentCtxWithCancel.Done():
			Shutdown(agentCtxWithCancel, cancelCtx, logg, ag, heartBeatUrl, client, errCh)
			return
		case <-errCh:
			Shutdown(agentCtxWithCancel, cancelCtx, logg, ag, heartBeatUrl, client, errCh)
			return
		default:
			go func() {
				for {
					select {
					case <-time.After(120 * time.Second):
						app.AgentHeartBeat(agentCtxWithCancel, logg, ag, heartBeatUrl, client, errCh)
					}
				}
			}()

		}
	}()
	mux := http.NewServeMux()
	mux.Handle("/", middleware.RecoveryMiddleware(middleware.LoggingMiddleware(server.GetSubExprassion(agentCtxWithCancel, logg), logg)))
	log.Fatal(http.ListenAndServe(":"+port, mux))

	// 	}(fmt.Sprintf("%d", port))
	// }

}
func Shutdown(ctx context.Context, cancelCtx context.CancelFunc, log *slog.Logger, ag agent.Agent, url string, client *http.Client, errCh chan struct{}) {
	ag.Status = "Error"
	app.AgentHeartBeat(ctx, log, ag, url, client, errCh)
	cancelCtx()
}
func setupLogger() *slog.Logger {
	logg := slog.New(slog.NewJSONHandler(
		os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return logg
}
