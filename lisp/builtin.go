package lisp

import (
	"errors"
	"fmt"
	"os"
)

type Builtin func(*Environment, []Node) Node

type IncorrectType struct {
	Expected string
	Actual   string
}

func (i IncorrectType) Error() string {
	return fmt.Sprintf("expected %v, got %v", i.Expected, i.Actual)
}

func Add(_ *Environment, args []Node) Node {
	sum := 0.0
	for _, node := range args {
		n, ok := node.(NumberNode)
		if !ok {
			return ErrorNode{IncorrectType{"Number", node.TypeString()}}
		}
		sum += float64(n)
	}
	return NumberNode(sum)
}

func Sub(_ *Environment, args []Node) Node {
	if len(args) == 1 {
		n, ok := args[0].(NumberNode)
		if !ok {
			return ErrorNode{IncorrectType{"Number", args[0].TypeString()}}
		}
		return -n
	}

	sum, ok := args[0].(NumberNode)
	if !ok {
		return ErrorNode{IncorrectType{"Number", args[0].TypeString()}}
	}

	for _, node := range args[1:] {
		n, ok := node.(NumberNode)
		if !ok {
			return ErrorNode{IncorrectType{"Number", node.TypeString()}}
		}
		sum -= n
	}

	return sum
}

func Mul(_ *Environment, args []Node) Node {
	sum, ok := args[0].(NumberNode)
	if !ok {
		return ErrorNode{IncorrectType{"Number", args[0].TypeString()}}
	}

	for _, node := range args[1:] {
		n, ok := node.(NumberNode)
		if !ok {
			return ErrorNode{IncorrectType{"Number", node.TypeString()}}
		}
		sum *= n
	}

	return sum
}

func Div(_ *Environment, args []Node) Node {
	sum, ok := args[0].(NumberNode)
	if !ok {
		return ErrorNode{IncorrectType{"Number", args[0].TypeString()}}
	}

	for _, node := range args[1:] {
		n, ok := node.(NumberNode)
		if !ok {
			return ErrorNode{IncorrectType{"Number", node.TypeString()}}
		}
		sum /= n
	}

	return sum
}

func Head(_ *Environment, args []Node) Node {
	if len(args) > 1 {
		return ErrorNode{fmt.Errorf("expected 1 argument, got %v", len(args))}
	}

	expr, ok := args[0].(ExpressionNode)
	if !ok || expr.Type != QExpression {
		return ErrorNode{IncorrectType{"Q-Expression", args[0].TypeString()}}
	}

	if len(expr.Nodes) == 0 {
		return ErrorNode{errors.New("cannot take head of empty list")}
	}

	return ExpressionNode{QExpression, expr.Nodes[:1]}
}

func Tail(_ *Environment, args []Node) Node {
	if len(args) > 1 {
		return ErrorNode{errors.New(fmt.Sprintf("expected 1 argument, got %v", len(args)))}
	}

	expr, ok := args[0].(ExpressionNode)
	if !ok || expr.Type != QExpression {
		return ErrorNode{IncorrectType{"Q-Expression", args[0].TypeString()}}
	}

	if len(expr.Nodes) == 0 {
		return ErrorNode{errors.New("cannot take tail of empty list")}
	}

	return ExpressionNode{QExpression, expr.Nodes[1:]}
}

func List(_ *Environment, args []Node) Node {
	return ExpressionNode{QExpression, args}
}

func Eval(env *Environment, args []Node) Node {
	if len(args) > 1 {
		return ErrorNode{errors.New(fmt.Sprintf("expected 1 argument, got %v", len(args)))}
	}

	expr, ok := args[0].(ExpressionNode)
	if ok {
		return expr.EvalAsSExpr(env)
	} else {
		return args[0].Evaluate(env)
	}
}

func Join(_ *Environment, args []Node) Node {
	nodes := make([]Node, 0)
	for _, n := range args {
		expr, ok := n.(ExpressionNode)
		if !ok || expr.Type != QExpression {
			return ErrorNode{IncorrectType{"Q-Expression", n.TypeString()}}
		}
		nodes = append(nodes, expr.Nodes...)
	}
	return ExpressionNode{QExpression, nodes}
}

func val(env *Environment, args []Node, global bool) Node {
	expr, ok := args[0].(ExpressionNode)
	if !ok || expr.Type != QExpression {
		return ErrorNode{IncorrectType{"Q-Expression", args[0].TypeString()}}
	}

	for _, node := range expr.Nodes {
		_, ok := node.(IdentifierNode)
		if !ok {
			return ErrorNode{IncorrectType{"Identifier", node.TypeString()}}
		}
	}

	args = args[1:]
	if len(args) != len(expr.Nodes) {
		return ErrorNode{fmt.Errorf("expected %v arguments, got %v", len(expr.Nodes)+1, len(args)+1)}
	}

	for i := range args {
		if global {
			env.Def(expr.Nodes[i].(IdentifierNode), args[i])
		} else {
			env.Put(expr.Nodes[i].(IdentifierNode), args[i])
		}
	}

	return ExpressionNode{Type: SExpression}
}

func Let(env *Environment, args []Node) Node {
	return val(env, args, false)
}

func Def(env *Environment, args []Node) Node {
	return val(env, args, true)
}

func Fn(env *Environment, args []Node) Node {
	if len(args) != 2 {
		return ErrorNode{fmt.Errorf("expected 2 arguments, got %v", len(args))}
	}

	for _, n := range args {
		expr, ok := n.(ExpressionNode)
		if !ok || expr.Type != QExpression {
			return ErrorNode{IncorrectType{"Q-Expression", n.TypeString()}}
		}
	}

	nodes := args[0].(ExpressionNode).Nodes
	formals := make([]IdentifierNode, 0, len(nodes))
	for _, node := range nodes {
		i, ok := node.(IdentifierNode)
		if !ok {
			return ErrorNode{IncorrectType{"Identifier", node.TypeString()}}
		}

		formals = append(formals, i)
	}

	subEnv := NewEnvironment(nil)

	return FunctionNode{
		Environment: &subEnv,
		Formals:     formals,
		Body:        args[1].(ExpressionNode),
	}
}

func Import(env *Environment, args []Node) Node {
	if len(args) != 1 {
		return ErrorNode{fmt.Errorf("expected 1 argument, got %v", len(args))}
	}

	path, ok := args[0].(StringNode)
	if !ok {
		return ErrorNode{IncorrectType{"String", args[0].TypeString()}}
	}

	file, err := os.ReadFile(string(path) + ".clsp")
	if err == nil {
		_, err := Evaluate(env, string(file), true)
		if err != nil {
			return ErrorNode{err}
		}
	}

	return ExpressionNode{Type: SExpression}
}
