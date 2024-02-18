package application

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"sync"
	"time"
	"unicode"

	"github.com/xleshka/distributedcalc/backend/internal/agent"
	expparser "github.com/xleshka/distributedcalc/backend/internal/application/ExpParser"
	"github.com/xleshka/distributedcalc/backend/internal/application/cache"
	orch "github.com/xleshka/distributedcalc/backend/internal/orchestrator"
)

var (
	agentCache      *cache.Cache
	operaionsCache  *cache.Cache
	expressionCache *cache.Cache
	agentChan       chan agent.Agent
)

func Initialize(ahentCount int, ctx context.Context, logg *slog.Logger, rep orch.Repository) error {
	agentCache = cache.NewCache()
	operaionsCache = cache.NewCache()
	expressionCache = cache.NewCache()
	agentChan := make(chan agent.Agent, ahentCount)
	_, err := AllOperations(ctx, logg, rep)
	if err != nil {
		log.Fatalf("%v", err)
		return err
	}

	_, err = AllOperations(ctx, logg, rep)
	if err != nil {
		log.Fatalf("%v", err)
		return err
	}
	_, err = AllAgents(ctx, logg, rep)
	if err != nil {
		log.Fatalf("%v", err)
		return err
	}
	for _, a := range agentCache.GetAll() {
		ag := a.(agent.Agent)
		if ag.Status == "Error" {
			continue
		}
		agentChan <- ag
	}
	if len(agentChan) == 0 {
		log.Fatal("Nil agents running")
		return errors.New("nil agents running")
	}
	return nil
}
func ValidExpression(expression string) bool {

	for _, r := range expression {
		/*проверяем число ли символ, если нет, то проверяем это один из валидных символов или нет*/
		switch {
		case unicode.Is(unicode.Digit, r):
			continue
		case r == ' ':
			continue
		case r == '+':
			continue
		case r == '-':
			continue
		case r == '*':
			continue
		case r == '/':
			continue
		default:
			return false
		}
	}
	return true
}

func AllExpressions(ctx context.Context, log *slog.Logger, rep orch.Repository, client *http.Client) ([]orch.Expression, error) {
	expressions := make([]orch.Expression, 0)
	if len(expressionCache.GetAll()) == 0 {
		rep.GetAllExpressions(ctx, expressionCache)
	}
	for _, expr := range expressionCache.GetAll() {
		expression, _ := expr.(orch.Expression)
		expressions = append(expressions, expression)
	}
	return expressions, nil
}
func SetExpression(ctx context.Context, express *orch.Expression, log *slog.Logger, rep orch.Repository) error {
	express.CompletedTime = time.Now()
	err := rep.SetExpression(ctx, express, expressionCache)
	if err != nil {
		return err
	}
	return nil
}
func AddExpression(ctx context.Context, express *orch.Expression, log *slog.Logger, rep orch.Repository, client *http.Client) (*orch.Expression, error) {

	expression, exist := rep.CheckExists(ctx, express.Expression)
	if exist {
		if expression.Status == "wait" {
			CalcExpression(ctx, log, express, rep, client, agentChan)
		}
		return &expression, nil
	}
	err := rep.Add(ctx, express, expressionCache)
	if err != nil {
		log.Error("failed add expression to bd")
		return &orch.Expression{}, err
	}
	if express.Status == "wait" {
		CalcExpression(ctx, log, express, rep, client, agentChan)
	}

	return express, nil
}
func CalcExpression(ctx context.Context, log *slog.Logger, expression *orch.Expression, rep orch.Repository, client *http.Client, agentChan chan agent.Agent) {
	go func() {
		expression.Status = "proccess"
		err := SetExpression(ctx, expression, log, rep)
		if err != nil {
			expression.Status = "error"
			SetExpression(ctx, expression, log, rep)
			log.Error("failed calculate exprassion: %v", err)
			return
		}
		res, err := expparser.ValidatedPostOrder(ctx, log, expression.Expression, client, agentChan)
		mu := sync.Mutex{}
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			expression.Status = "error"
			SetExpression(ctx, expression, log, rep)
			log.Error("failed calculate exprassion: %v", err)
			return
		}
		expression.Expression += fmt.Sprintf(" = %0.1f", res)
		expression.Status = "calculated"
		SetExpression(ctx, expression, log, rep)
	}()
}
func AllOperations(ctx context.Context, log *slog.Logger, rep orch.Repository) ([]orch.Operation, error) {
	operations := make([]orch.Operation, 0)
	if len(operaionsCache.Data) == 0 {
		rep.GetAllOperations(ctx, operaionsCache)
		if len(operaionsCache.Data) < 4 {
			oprtrs := []string{
				"+", "-", "*", "/",
			}
			for _, o := range oprtrs {
				op := orch.Operation{Operation: o,
					ExecutionTimeByMilliseconds: 2000}
				err := AddOperation(ctx, log, op, rep)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	for _, operation := range operaionsCache.GetAll() {
		op, _ := operation.(orch.Operation)
		operations = append(operations, op)
	}
	log.Info(fmt.Sprintf("operations: %v", operations))
	return operations, nil
}
func AddOperation(ctx context.Context, log *slog.Logger, operation orch.Operation, rep orch.Repository) error {
	fmt.Println(operation)
	err := rep.AddOperation(ctx, &operation, operaionsCache)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}
func SetOperation(ctx context.Context, log *slog.Logger, operaion orch.Operation, rep orch.Repository) error {

	err := rep.SetExecutionTime(ctx, &operaion, operaionsCache)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func AllAgents(ctx context.Context, log *slog.Logger, rep orch.Repository) ([]agent.Agent, error) {
	agents := make([]agent.Agent, 0)
	if len(agentCache.Data) == 0 {
		rep.GetAllAgents(ctx, agentCache)
	}
	for _, ag := range agentCache.Data {
		a, _ := ag.(agent.Agent)
		if a.Status == "Error" {
			go DeleteAgent(ctx, log, a.ID, rep)
		}
		agents = append(agents, a)
	}
	return agents, nil
}

func SetAgent(ctx context.Context, log *slog.Logger, id, status string, rep orch.Repository) error {
	if status == "Error" {
		go DeleteAgent(ctx, log, id, rep)
		return nil
	}
	err := rep.SetAgent(ctx, id, status, agentCache)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func AddAgent(ctx context.Context, log *slog.Logger, ag agent.Agent, rep orch.Repository) (string, error) {
	err := rep.AddAgent(ctx, &ag, agentCache)
	if err != nil {
		return "", err
	}
	return ag.ID, nil
}
func AddAgentReq(ctx context.Context, log *slog.Logger, ag agent.Agent, url string, client *http.Client) (string, error) {

	data, err := json.Marshal(ag)
	if err != nil {
		log.Error("failed marshal agent: %s, error: %v", ag.Address, err)
		return "", err
	}
	r := bytes.NewReader(data)
	resp, err := http.Post(url, "application/json", r)
	if err != nil {
		log.Error("failed req: %v", err)
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Error("failed read resp body: %v", err)
		return "", err
	}
	log.Info(fmt.Sprintf("add agent: %s", string(body)))
	return string(body), nil
}
func AgentHeartBeat(ctx context.Context, log *slog.Logger, ag agent.Agent, url string, client *http.Client, errCh chan struct{}) {
	log.Info("AgentHeartBeat")
	data, err := json.Marshal(ag)
	if err != nil {
		log.Error("failed marshal agent: %s, error: %v", ag.Address, err)
		errCh <- struct{}{}
		return
	}
	r := bytes.NewReader(data)
	resp, err := http.Post(url, "application/json", r)
	if err != nil {
		log.Error("failed req: %v", err)
		errCh <- struct{}{}
		return
	}
	log.Info(fmt.Sprintf("heart beat: %s", resp.Status))
}
func timerStop(v time.Timer, ctx context.Context, id string, rep orch.Repository) error {
	<-v.C
	err := rep.DeleteAgent(ctx, id, agentCache)
	if err != nil {
		return err
	}
	return nil
}
func DeleteAgent(ctx context.Context, log *slog.Logger, id string, rep orch.Repository) {
	timer := time.NewTimer(60 * time.Second)

	err := timerStop(*timer, ctx, id, rep)
	if err != nil {
		log.Error(err.Error())
	}
}
