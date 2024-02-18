package expparser

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/xleshka/distributedcalc/backend/internal/agent"
	"github.com/xleshka/distributedcalc/backend/internal/application/cache"
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
func CalcExpression(node *Node, counter *int, log *slog.Logger, client *http.Client, agentChache *cache.Cache) {

	if node == nil {
		return
	}
	if node.Left != nil {
		go CalcExpression(node.Left, counter, log, client, agentChache)
	}
	if node.Right != nil {
		go CalcExpression(node.Right, counter, log, client, agentChache)
	}
	if node.Left == nil && node.Right == nil {
		return
	}
	if node.Operation != "" {
		subExpression := fmt.Sprintf("%0.1f %s %0.1f", node.Left.Value, node.Operation, node.Right.Value)
		// app.pos
		go func() {
			node.Value, _ = PostTask(log, subExpression, client, agentChache)

		}()
		return
	}
}

func PostTask(log *slog.Logger, exprassion string, client *http.Client, agentCache *cache.Cache) (float64, error) {
	expr := SubExpression{Expression: exprassion}

	data, err := json.Marshal(expr)

	var ag agent.Agent
	br := false
	for {
		for _, r := range agentCache.Data {
			ag = r.(agent.Agent)
			if ag.Status != "Error" && ag.Status != "Busy" {
				br = true
				break
			}
		}
		if br {
			break
		}
	}

	req, err := http.NewRequest("POST", ag.Address, bytes.NewBuffer(data))
	if err != nil {
		log.Error("failed req POSTTASK")
		return 0, err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("failed post subExprassion to agent: %v", err)
		return 0, err
	}
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Error("failed read resp body: %v", err)
		return 0, err
	}
	res, err := strconv.ParseFloat(string(body), 64)
	if err != nil {
		log.Error("failed parse resp body from bytes to float: %v", err)
		return 0, err
	}
	return res, nil
}

func ValidatedPostOrder(ctx context.Context, log *slog.Logger, expression string, client *http.Client, agentCache *cache.Cache) (float64, error) {
	node, err := ExpressionParser(expression)
	if err != nil {
		return 0, err
	}
	var counter int
	CalcExpression(node, &counter, log, client, agentCache)
	return node.Value, nil
}
