package expparser

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/xleshka/distributedcalc/backend/internal/agent"
)

type Node struct {
	Left      *Node
	Right     *Node
	Operation string
	Value     float64
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
func CalcExpression(node *Node, counter *int, waitg *sync.WaitGroup, log *slog.Logger, client *http.Client, agentChan chan agent.Agent) {
	defer waitg.Done()
	if node == nil {
		return
	}
	wg := &sync.WaitGroup{}
	if node.Left != nil {
		wg.Add(1)
		go CalcExpression(node.Left, counter, wg, log, client, agentChan)
	}
	if node.Right != nil {
		wg.Add(1)
		go CalcExpression(node.Right, counter, wg, log, client, agentChan)
	}
	if node.Left == nil && node.Right == nil {
		return
	}
	wg.Wait()
	if node.Operation != "" {
		subExpression := fmt.Sprintf("%0.1f %s %0.1f", node.Left.Value, node.Operation, node.Right.Value)
		// app.pos
		var err error
		for {
			node.Value, err = PostTask(log, subExpression, client, agentChan)
			if err == nil {
				break
			}
		}
		return
	}
}

func PostTask(log *slog.Logger, exprassion string, client *http.Client, agentChan chan agent.Agent) (float64, error) {
	data := []byte(exprassion)
	r := bytes.NewReader(data)
	var ag agent.Agent
	for {
		ag := <-agentChan
		if ag.Status != "Error" {
			agentChan <- ag
			break
		}
		agentChan <- ag
	}

	resp, err := http.Post(ag.Address, "", r)
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

func ValidatedPostOrder(ctx context.Context, log *slog.Logger, expression string, client *http.Client, agentChan chan agent.Agent) (float64, error) {
	node, err := ExpressionParser(expression)
	if err != nil {
		return 0, err
	}
	var counter int
	wg := &sync.WaitGroup{}
	wg.Add(1)
	CalcExpression(node, &counter, wg, log, client, agentChan)
	wg.Wait()
	return node.Value, nil
}
