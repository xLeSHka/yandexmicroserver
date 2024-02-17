package expparser

import (
	"errors"
	"strconv"
	"strings"
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
		default:
			val, err := strconv.ParseFloat(member, 256)
			if err != nil {
				return nil, err
			}
			stack = append(stack, &Node{Value: val})
		}
	}
	for len(operators) > 0 {
		popMember(&stack, &operators)
	}
	if len(stack) == 1 {
		return nil, errors.New("Expression failed parse")
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

// func ValidatedPostOrder(s string) (float64, error) {
// 	out := make(chan float64)
// 	node, err := ExpressionParser(s)
// 	if err != nil {
// 		return 0, err
// 	}
// 	var counter int
// 	err = EvaluatePostOrder(node, &counter, out)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return <-out, nil
// }
