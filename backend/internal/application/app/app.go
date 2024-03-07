package application

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strings"
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
)

func Initialize(agentCount int, ctx context.Context, logg *slog.Logger, rep orch.Repository, client *http.Client) error {
	agentCache = cache.NewCache()
	operaionsCache = cache.NewCache()
	expressionCache = cache.NewCache()

	_, err := AllAgents(ctx, logg, rep)
	if err != nil {
		log.Fatalf("%v", err)
		return err
	}

	_, err = AllOperations(ctx, logg, rep)
	if err != nil {
		log.Fatalf("%v", err)
		return err
	}

	err = InitAllExpressions(ctx, logg, rep)
	if err != nil {
		log.Fatalf("%v", err)
		return err
	}

	return nil
}
func CORSPOST(w http.ResponseWriter, r *http.Request) {

}
func CORSGET(w http.ResponseWriter, r *http.Request) {

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
func ValidOperation(oper string) []string {
	return strings.Fields(oper)
}

// func Cacl(Expression string, log *slog.Logger) string {
// 	out := make(chan float64)
// 	CaclSubExprassion(Expression, out, wg, log)
// 	close(out)
// 	return fmt.Sprintf("%0.1f", <-out)
// }
// func CaclSubExprassion(Expression string, out chan<- float64, wg *sync.WaitGroup, log *slog.Logger) {
// 	defer wg.Done()

//		members := strings.Fields(Expression)
//		log.Info("CaclSubExprassion")
//		switch members[1] {
//		case "+":
//			data, _ := operaionsCache.Get("+")
//			plusOper := data.(orch.Operation)
//			timer := time.NewTimer(time.Duration(plusOper.ExecutionTimeByMilliseconds * int(time.Millisecond)))
//			ExecutionTimer(*timer)
//			left, err := strconv.ParseFloat(members[0], 64)
//			if err != nil {
//				log.Error("failed parse: %v", err)
//				return
//			}
//			right, err := strconv.ParseFloat(members[0], 64)
//			if err != nil {
//				log.Error("failed parse: %v", err)
//				return
//			}
//			out <- left + right
//		case "-":
//			data, _ := operaionsCache.Get("-")
//			minusOper := data.(orch.Operation)
//			timer := time.NewTimer(time.Duration(minusOper.ExecutionTimeByMilliseconds * int(time.Millisecond)))
//			ExecutionTimer(*timer)
//			left, err := strconv.ParseFloat(members[0], 64)
//			if err != nil {
//				log.Error("failed parse: %v", err)
//				return
//			}
//			right, err := strconv.ParseFloat(members[0], 64)
//			if err != nil {
//				log.Error("failed parse: %v", err)
//				return
//			}
//			out <- left - right
//		case "*":
//			data, _ := operaionsCache.Get("*")
//			miltiplyOper := data.(orch.Operation)
//			timer := time.NewTimer(time.Duration(miltiplyOper.ExecutionTimeByMilliseconds * int(time.Millisecond)))
//			ExecutionTimer(*timer)
//			left, err := strconv.ParseFloat(members[0], 64)
//			if err != nil {
//				log.Error("failed parse: %v", err)
//				return
//			}
//			right, err := strconv.ParseFloat(members[0], 64)
//			if err != nil {
//				log.Error("failed parse: %v", err)
//				return
//			}
//			out <- left * right
//		case "/":
//			data, _ := operaionsCache.Get("/")
//			divideOper := data.(orch.Operation)
//			timer := time.NewTimer(time.Duration(divideOper.ExecutionTimeByMilliseconds * int(time.Millisecond)))
//			ExecutionTimer(*timer)
//			left, err := strconv.ParseFloat(members[0], 64)
//			if err != nil {
//				log.Error("failed parse: %v", err)
//				return
//			}
//			right, err := strconv.ParseFloat(members[0], 64)
//			if err != nil {
//				log.Error("failed parse: %v", err)
//				return
//			}
//			out <- left / right
//		}
//	}
//
//	func ExecutionTimer(v time.Timer) {
//		<-v.C
//	}
func InitAllExpressions(ctx context.Context, log *slog.Logger, rep orch.Repository) error {
	rep.GetAllExpressions(ctx, expressionCache)
	for _, expr := range expressionCache.Data {
		expression, _ := expr.(orch.Expression)
		if expression.Status == "wait" {
			go CalcExpression(ctx, log, expression, rep)
		}
	}
	return nil
}
func AllExpressions(ctx context.Context, log *slog.Logger, rep orch.Repository) ([]orch.Expression, error) {
	expressions := make([]orch.Expression, 0)
	for _, expr := range expressionCache.Data {
		expression, _ := expr.(orch.Expression)
		expressions = append(expressions, expression)
		if expression.Status == "wait" {
			go CalcExpression(ctx, log, expression, rep)
		}
	}
	return expressions, nil
}
func SetExpression(ctx context.Context, express *orch.Expression, log *slog.Logger, rep orch.Repository) error {
	express.CompletedTime = time.Now()
	err := rep.SetExpression(ctx, *express, expressionCache)
	if err != nil {
		return err
	}
	return nil
}
func AddExpression(ctx context.Context, express *orch.Expression, log *slog.Logger, rep orch.Repository) (*orch.Expression, error) {

	expression, exist := rep.CheckExists(ctx, express.Expression)
	if exist {
		if expression.Status == "wait" {
			go CalcExpression(ctx, log, *express, rep)
		}
		return &expression, nil
	}
	err := rep.Add(ctx, express, expressionCache)
	if err != nil {
		log.Error("failed add expression to bd")
		return &orch.Expression{}, err
	}
	log.Info(express.Id, express.Expression, express.Status)
	if express.Status == "wait" {
		go CalcExpression(ctx, log, *express, rep)
	}
	log.Info(express.Id, express.Expression, express.Status)
	return express, nil
}
func CalcExpression(ctx context.Context, log *slog.Logger, expression orch.Expression, rep orch.Repository) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("failed calculate exprassion: %v", err)
			expression.Status = "wait"
			err := SetExpression(ctx, &expression, log, rep)
			if err != nil {
				log.Error("failed calculate exprassion: %v", err)
				return
			}
		}
	}()
	expression.Status = "proccess"

	err := SetExpression(ctx, &expression, log, rep)
	if err != nil {
		log.Error("failed calculate exprassion: %v", err)
		return
	}
	res, err := expparser.ValidatedPostOrder(ctx, log, expression.Expression, operaionsCache)

	if err != nil {
		expression.Status = "error"
		err := SetExpression(ctx, &expression, log, rep)
		if err != nil {
			return
		}
		log.Error("failed calculate exprassion: %v", err)
		return
	}
	expression.Expression += fmt.Sprintf(" = %0.7f", res)
	expression.Status = "calculated"
	err = SetExpression(ctx, &expression, log, rep)
	if err != nil {
		return
	}

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

	for _, operation := range operaionsCache.Data {
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

func SetAgent(ctx context.Context, log *slog.Logger, ag agent.Agent, rep orch.Repository) error {
	if ag.Status == "Error" {
		go DeleteAgent(ctx, log, ag.Address, rep)
	}
	err := rep.SetAgent(ctx, ag, agentCache)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func AddAgent(ctx context.Context, log *slog.Logger, ag agent.Agent, rep orch.Repository) (string, error) {

	agnt, found := rep.ChechIfExist(ctx, ag)
	if found {
		return agnt.ID, nil
	}
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
func timerStop(v time.Timer, ctx context.Context, address string, rep orch.Repository) error {
	<-v.C
	err := rep.DeleteAgent(ctx, address, agentCache)
	if err != nil {
		return err
	}
	return nil
}

func DeleteAgent(ctx context.Context, log *slog.Logger, address string, rep orch.Repository) {
	timer := time.NewTimer(20 * time.Second)

	err := timerStop(*timer, ctx, address, rep)
	if err != nil {
		log.Error(err.Error())
	}
}
