package expparser

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/xleshka/distributedcalc/backend/internal/application/cache"
	orch "github.com/xleshka/distributedcalc/backend/internal/orchestrator"
)

type Node struct {
	Left      *Node
	Right     *Node
	Operation string
	Value     float64
}

type SubExpression struct {
	Expression string `json:"expression"`
}

func ExpressionParser(s string) (*Node, error) {
	var (
		members   = strings.Fields(s)
		stack     []*Node
		operators []string
	)
	for _, member := range members {
		switch member {
		case "+", "-", "*", "/":
			for len(operators) > 0 && precedence(operators[len(operators)-1]) >= precedence(member) {
				popMember(&stack, &operators)
			}
			operators = append(operators, member)
		default:
			val, err := strconv.ParseFloat(member, 64)
			if err != nil {
				return nil, err
			}
			stack = append(stack, &Node{Value: val})
		}
	}
	for len(operators) > 0 {
		popMember(&stack, &operators)
	}
	if len(stack) != 1 {
		return nil, errors.New("expression failed parse")
	}
	return stack[0], nil
}

/*порядок подвыражения и чистит стек от использованных значений*/
func precedence(s string) int {
	switch s {
	case "*", "/":
		return 2
	case "+", "-":
		return 1
	default:
		return 0
	}
}

/*формирует ноду*/
func popMember(stack *[]*Node, operators *[]string) {
	operator := (*operators)[len(*operators)-1]
	*operators = (*operators)[:len(*operators)-1]

	right := (*stack)[len(*stack)-1]
	*stack = (*stack)[:len(*stack)-1]

	left := (*stack)[len(*stack)-1]
	*stack = (*stack)[:len(*stack)-1]

	node := &Node{Right: right, Left: left, Operation: operator}
	*stack = append(*stack, node)
}
func CalcExpression(node *Node, log *slog.Logger, operaionsCache *cache.Cache, wg *sync.WaitGroup, cancelCtx context.CancelFunc) {
	defer wg.Done()
	if node == nil {
		return
	}
	wg1 := &sync.WaitGroup{}
	if node.Left != nil {
		wg1.Add(1)
		go CalcExpression(node.Left, log, operaionsCache, wg1, cancelCtx)
	}
	if node.Right != nil {
		wg1.Add(1)
		go CalcExpression(node.Right, log, operaionsCache, wg1, cancelCtx)
	}
	if node.Left == nil && node.Right == nil {
		return
	}
	wg1.Wait()
	if node.Operation != "" {
		subExpression := fmt.Sprintf("%0.1f %s %0.1f", node.Left.Value, node.Operation, node.Right.Value)
		res, err := Calc(subExpression, log, operaionsCache)
		if err != nil {
			cancelCtx()
		}
		node.Value = res
		return
	}
}

func Calc(Expression string, log *slog.Logger, operaionsCache *cache.Cache) (float64, error) {
	res, err := CalcSubExpression(Expression, log, operaionsCache)
	if err != nil {
		return 0, err
	}
	return res, nil
}
func CalcSubExpression(Expression string, log *slog.Logger, operaionsCache *cache.Cache) (float64, error) {
	members := strings.Fields(Expression)
	switch members[1] {
	case "+":
		data, _ := operaionsCache.Get("+")
		plusOper := data.(orch.Operation)
		time.Sleep(time.Duration(plusOper.ExecutionTimeByMilliseconds * int(time.Millisecond)))
		left, err := strconv.ParseFloat(members[0], 64)
		if err != nil {
			log.Error("failed parse: %v", err)
			return 0, err
		}
		right, err := strconv.ParseFloat(members[2], 64)
		if err != nil {
			log.Error("failed parse: %v", err)
			return 0, err
		}

		res := left + right
		log.Info("plus out: ", Expression, res)

		return res, nil
	case "-":
		data, _ := operaionsCache.Get("-")
		minusOper := data.(orch.Operation)
		time.Sleep(time.Duration(minusOper.ExecutionTimeByMilliseconds * int(time.Millisecond)))
		left, err := strconv.ParseFloat(members[0], 64)
		if err != nil {
			log.Error("failed parse: %v", err)
			return 0, err
		}
		right, err := strconv.ParseFloat(members[2], 64)
		if err != nil {
			log.Error("failed parse: %v", err)
			return 0, err
		}
		res := left - right
		log.Info("minus out: ", Expression, res)

		return res, nil
	case "*":
		data, _ := operaionsCache.Get("*")
		miltiplyOper := data.(orch.Operation)
		time.Sleep(time.Duration(miltiplyOper.ExecutionTimeByMilliseconds * int(time.Millisecond)))
		left, err := strconv.ParseFloat(members[0], 64)
		if err != nil {
			log.Error("failed parse: %v", err)
			return 0, err
		}
		right, err := strconv.ParseFloat(members[2], 64)
		if err != nil {
			log.Error("failed parse: %v", err)
			return 0, err
		}
		res := left * right
		log.Info("multiply out: ", Expression, res)

		return res, nil
	case "/":
		data, _ := operaionsCache.Get("/")
		divideOper := data.(orch.Operation)
		time.Sleep(time.Duration(divideOper.ExecutionTimeByMilliseconds * int(time.Millisecond)))
		left, err := strconv.ParseFloat(members[0], 64)
		if err != nil {
			log.Error("failed parse: %v", err)
			return 0, err
		}
		right, err := strconv.ParseFloat(members[2], 64)
		if err != nil {
			log.Error("failed parse: %v", err)
			return 0, err
		} else if right == 0 {
			log.Error("divide to null")
			return 0, errors.New("divide to null")
		}

		res := left / right
		log.Info("divide out ", Expression, res)

		return res, nil
	}
	return 0, nil
}

// func PostTask(log *slog.Logger, exprassion string, client *http.Client, agentCache *cache.Cache, errCh chan error) float64 {
// 	expr := SubExpression{Expression: exprassion}

// 	data, err := json.Marshal(expr)
// 	if err != nil {
// 		log.Error("failed marshal data to post task: %v", err)
// 		errCh <- err
// 		return 0
// 	}
// 	r := rand.New(rand.NewSource(time.Now().UnixNano()))
// 	rnd := r.Intn(len(agentCache.Data))
// 	agentAddres := 3030 + rnd

// 	req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:%d", agentAddres), bytes.NewBuffer(data))
// 	if err != nil {
// 		log.Error("failed req POSTTASK: %v", err)
// 		errCh <- err
// 		return 0
// 	}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		log.Error("failed post subExprassion to agent: %v", err)
// 		errCh <- err
// 		return 0
// 	}
// 	body, err := io.ReadAll(resp.Body)
// 	defer resp.Body.Close()
// 	if err != nil {
// 		log.Error("failed read resp body: %v", err)
// 		errCh <- err
// 		return 0
// 	}
// 	res, err := strconv.ParseFloat(string(body), 64)
// 	if err != nil {
// 		log.Error("failed parse resp body from bytes to float: %v", err)
// 		errCh <- err
// 		return 0
// 	}
// 	return res
// }

func ValidatedPostOrder(ctx context.Context, log *slog.Logger, expression string, operaionsCache *cache.Cache) (float64, error) {
	node, err := ExpressionParser(expression)
	if err != nil {
		return 0, err
	}

	ctx, cancelCtx := context.WithCancel(ctx)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	t := time.Now()
	CalcExpression(node, log, operaionsCache, wg, cancelCtx)
	select {
	case <-ctx.Done():
		return 0, errors.New("Error calculate")
	default:
		wg.Wait()
		log.Info("Calc duration ", time.Since(t))
		return node.Value, nil
	}
}
