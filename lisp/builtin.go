package lisp

import (
	"errors"
	"fmt"
)

func Add(args []Node) Node {
	sum := 0
	for _, n := range args {
		n, ok := n.(NumberNode)
		if !ok {
			return ErrorNode{errors.New("cannot operate on non-number")}
		}
		sum += int(n)
	}
	return NumberNode(sum)
}

func Sub(args []Node) Node {
	if len(args) == 1 {
		n, ok := args[0].(NumberNode)
		if !ok {
			return ErrorNode{errors.New("cannot negate a non-number")}
		}
		return -n
	}

	sum, ok := args[0].(NumberNode)
	if !ok {
		return ErrorNode{errors.New("cannot operate on non-number")}
	}

	for _, n := range args[1:] {
		n, ok := n.(NumberNode)
		if !ok {
			return ErrorNode{errors.New("cannot operate on non-number")}
		}
		sum -= n
	}

	return sum
}

func Mul(args []Node) Node {
	sum, ok := args[0].(NumberNode)
	if !ok {
		return ErrorNode{errors.New("cannot operate on non-number")}
	}

	for _, n := range args[1:] {
		n, ok := n.(NumberNode)
		if !ok {
			return ErrorNode{errors.New("cannot operate on non-number")}
		}
		sum *= n
	}

	return sum
}

func Div(args []Node) Node {
	sum, ok := args[0].(NumberNode)
	if !ok {
		return ErrorNode{errors.New("cannot operate on non-number")}
	}

	for _, n := range args[1:] {
		n, ok := n.(NumberNode)
		if !ok {
			return ErrorNode{errors.New("cannot operate on non-number")}
		}
		sum /= n
	}

	return sum
}

func Head(args []Node) Node {
	if len(args) > 1 {
		return ErrorNode{errors.New(fmt.Sprintf("expected 1 argument, got %v", len(args)))}
	}

	expr, ok := args[0].(ExpressionNode)
	if !ok || expr.Type != QExpression {
		return ErrorNode{errors.New(fmt.Sprintf("cannot operate on non-qexpression"))}
	}

	if len(expr.Nodes) == 0 {
		return ErrorNode{errors.New("cannot take head of empty list")}
	}

	return expr.Nodes[0]
}

func Tail(args []Node) Node {
	if len(args) > 1 {
		return ErrorNode{errors.New(fmt.Sprintf("expected 1 argument, got %v", len(args)))}
	}

	expr, ok := args[0].(ExpressionNode)
	if !ok || expr.Type != QExpression {
		return ErrorNode{errors.New(fmt.Sprintf("cannot operate on non-qexpression"))}
	}

	if len(expr.Nodes) == 0 {
		return ErrorNode{errors.New("cannot take tail of empty list")}
	}

	return ExpressionNode{QExpression, expr.Nodes[1:]}
}

func List(args []Node) Node {
	return ExpressionNode{QExpression, args}
}

func Eval(args []Node) Node {
	if len(args) > 1 {
		return ErrorNode{errors.New(fmt.Sprintf("expected 1 argument, got %v", len(args)))}
	}

	expr, ok := args[0].(ExpressionNode)
	if !ok {
		return ErrorNode{errors.New(fmt.Sprintf("cannot operate on non-qexpression"))}
	}

	return expr.EvalAsSExpr()
}

func Join(args []Node) Node {
	nodes := make([]Node, 0)
	for _, n := range args {
		expr, ok := n.(ExpressionNode)
		if !ok || expr.Type != QExpression {
			return ErrorNode{errors.New("cannot operate on non-qexpression")}
		}
		nodes = append(nodes, expr.Nodes...)
	}
	return ExpressionNode{QExpression, nodes}
}
