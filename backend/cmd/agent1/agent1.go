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
	agentPort, err := strconv.Atoi(cfg.AgentPort1)
	if err != nil {
		log.Fatalf("invalid agent count: %v", err)
	}
	heartBeatUrl := "http://" + cfg.HTTPServer.Host + ":" + cfg.OrchPort + "/setAgentStatus"
	addAgentUrl := "http://" + cfg.HTTPServer.Host + ":" + cfg.OrchPort + "/addAgent"
	port := fmt.Sprintf("%d", agentPort)
	errCh := make(chan struct{})
	defer close(errCh)
	agentCtxWithCancel, cancelCtx := context.WithCancel(context.Background())
	agentAddress := "http://localhost:" + port
	client := &http.Client{}
	ag := agent.Agent{Address: agentAddress, Status: "Ok"}
	ag.LastHearBeat = time.Now()
	ag.ID, err = app.AddAgentReq(agentCtxWithCancel, logg, ag, addAgentUrl, client)

	if err != nil {
		ag.Status = "Error"
		Shutdown(agentCtxWithCancel, cancelCtx, logg, ag, heartBeatUrl, client, errCh)
		return
	}
	go func() {
		for {
			select {
			case <-agentCtxWithCancel.Done():
				ag.Status = "Error"
				Shutdown(agentCtxWithCancel, cancelCtx, logg, ag, heartBeatUrl, client, errCh)
				return
			case <-errCh:
				ag.Status = "Error"
				Shutdown(agentCtxWithCancel, cancelCtx, logg, ag, heartBeatUrl, client, errCh)
				return
			case <-time.After(120 * time.Second):
				ag.LastHearBeat = time.Now()
				app.AgentHeartBeat(agentCtxWithCancel, logg, ag, heartBeatUrl, client, errCh)
			}
		}
	}()
	mux := http.NewServeMux()
	// mux.Handle("/", middleware.RecoveryMiddleware(middleware.LoggingMiddleware(server.GetSubExprassion(agentCtxWithCancel, logg, ag, heartBeatUrl, client), logg)))
	log.Fatal(http.ListenAndServe(":"+port, mux))

	// 	}(fmt.Sprintf("%d", port))
	// }

}
func Shutdown(ctx context.Context, cancelCtx context.CancelFunc, log *slog.Logger, ag agent.Agent, url string, client *http.Client, errCh chan struct{}) {

	app.AgentHeartBeat(ctx, log, ag, url, client, errCh)
	cancelCtx()
}
func setupLogger() *slog.Logger {
	logg := slog.New(slog.NewJSONHandler(
		os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return logg
}
