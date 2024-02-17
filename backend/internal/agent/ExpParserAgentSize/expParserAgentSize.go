package expparser

import (
	"errors"
	"sync"
)

type Node struct {
	Right    *Node
	Left     *Node
	Operator string
	Value    float64
}

func CalcExpression(node *Node, subExprassion *map[int]string, counter *int) error {
	// defer wgArg.Done()
	if node == nil {
		return errors.New("nil node")
	}
	// ch := make(chan float64, 2)
	wg := &sync.WaitGroup{}
	if node.Left != nil {
		wg.Add(1)
		// err := CalcExpression(node.Left, counter, wg,ch)
		// if err != nil {
		// 	return err
		// }
	}
	if node.Right != nil {
		wg.Add(1)
		// err := CalcExpression(node.Right, counter, wg, ch)
		// if err != nil {
		// 	return err
		// }
	}
	if node.Left == nil && node.Right == nil {
		// chArg <- node.Value
		return nil
	}
	wg.Wait()
	if node.Operator != "" {
		// val,val1 := <-ch, <-ch
		// chArg <- val+val1
		return nil
	}
	return nil
}
