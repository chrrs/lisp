package lisp

import (
	"errors"
	"fmt"
	"math"
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

	switch v := args[0].(type) {
	case ExpressionNode:
		if v.Type != QExpression {
			return ErrorNode{IncorrectType{"Q-Expression or String", args[0].TypeString()}}
		}

		if len(v.Nodes) == 0 {
			return ErrorNode{errors.New("cannot take head of empty list")}
		}

		return ExpressionNode{QExpression, v.Nodes[:1]}
	case StringNode:
		if len(v) == 0 {
			return ErrorNode{errors.New("cannot take tail of empty string")}
		}

		return v[:1]
	default:
		return ErrorNode{IncorrectType{"Q-Expression or String", args[0].TypeString()}}
	}
}

func Tail(_ *Environment, args []Node) Node {
	if len(args) > 1 {
		return ErrorNode{errors.New(fmt.Sprintf("expected 1 argument, got %v", len(args)))}
	}

	switch v := args[0].(type) {
	case ExpressionNode:
		if v.Type != QExpression {
			return ErrorNode{IncorrectType{"Q-Expression or String", args[0].TypeString()}}
		}

		if len(v.Nodes) == 0 {
			return ErrorNode{errors.New("cannot take tail of empty list")}
		}

		return ExpressionNode{QExpression, v.Nodes[1:]}
	case StringNode:
		if len(v) == 0 {
			return ErrorNode{errors.New("cannot take tail of empty string")}
		}

		return v[1:]
	default:
		return ErrorNode{IncorrectType{"Q-Expression or String", args[0].TypeString()}}
	}
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
	switch v := args[0].(type) {
	case ExpressionNode:
		nodes := make([]Node, 0)

		for _, n := range args {
			expr, ok := n.(ExpressionNode)
			if !ok || expr.Type != QExpression {
				return ErrorNode{IncorrectType{"Q-Expression or String", n.TypeString()}}
			}
			nodes = append(nodes, expr.Nodes...)
		}

		return ExpressionNode{QExpression, nodes}
	case StringNode:
		ret := StringNode("")

		for _, n := range args {
			str, ok := n.(StringNode)
			if !ok {
				return ErrorNode{IncorrectType{"Q-Expression or String", n.TypeString()}}
			}
			ret += str
		}

		return ret
	default:
		return ErrorNode{IncorrectType{"Q-Expression or String", v.TypeString()}}
	}
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

func Fn(_ *Environment, args []Node) Node {
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

func Equal(env *Environment, args []Node) Node {
	if len(args) < 2 {
		return ErrorNode{fmt.Errorf("expected 2 or more arguments, got %v", len(args))}
	}

	first := args[0]

	for _, arg := range args[1:] {
		switch first.(type) {
		case NumberNode:
			n, ok := arg.(NumberNode)
			if !ok {
				return NumberNode(0)
			}

			if first != n {
				return NumberNode(0)
			}
		case StringNode:
			str, ok := arg.(StringNode)
			if !ok {
				return NumberNode(0)
			}

			if first != str {
				return NumberNode(0)
			}
		case ExpressionNode:
			first := first.(ExpressionNode)
			expr, ok := arg.(ExpressionNode)
			if !ok {
				return NumberNode(0)
			}

			if first.Type != expr.Type || len(first.Nodes) != len(expr.Nodes) {
				return NumberNode(0)
			}

			for i, node := range first.Nodes {
				if Equal(env, []Node{expr.Nodes[i], node}) == NumberNode(0) {
					return NumberNode(0)
				}
			}
		default:
			return ErrorNode{fmt.Errorf("unimplemented equality for %v", first.TypeString())}
		}
	}

	return NumberNode(1)
}

func If(env *Environment, args []Node) Node {
	if len(args) != 3 {
		return ErrorNode{fmt.Errorf("expected 3 arguments, got %v", len(args))}
	}

	condition, ok := args[0].(NumberNode)
	if !ok {
		return ErrorNode{IncorrectType{"Number", args[0].TypeString()}}
	}

	yes, ok := args[1].(ExpressionNode)
	if !ok || yes.Type != QExpression {
		return ErrorNode{IncorrectType{"Q-Expression", args[1].TypeString()}}
	}

	no, ok := args[2].(ExpressionNode)
	if !ok || no.Type != QExpression {
		return ErrorNode{IncorrectType{"Q-Expression", args[2].TypeString()}}
	}

	if condition != NumberNode(0) {
		return yes.EvalAsSExpr(env)
	} else {
		return no.EvalAsSExpr(env)
	}
}

func Mod(_ *Environment, args []Node) Node {
	if len(args) != 2 {
		return ErrorNode{fmt.Errorf("expected 1 argument, got %v", len(args))}
	}

	n, ok := args[0].(NumberNode)
	if !ok {
		return ErrorNode{IncorrectType{"Number", args[0].TypeString()}}
	}

	mod, ok := args[1].(NumberNode)
	if !ok {
		return ErrorNode{IncorrectType{"Number", args[1].TypeString()}}
	}

	return NumberNode(math.Mod(float64(n), float64(mod)))
}

func Less(_ *Environment, args []Node) Node {
	if len(args) != 2 {
		return ErrorNode{fmt.Errorf("expected 1 argument, got %v", len(args))}
	}

	n1, ok := args[0].(NumberNode)
	if !ok {
		return ErrorNode{IncorrectType{"Number", args[0].TypeString()}}
	}

	n2, ok := args[1].(NumberNode)
	if !ok {
		return ErrorNode{IncorrectType{"Number", args[1].TypeString()}}
	}

	if n1 < n2 {
		return NumberNode(1)
	} else {
		return NumberNode(0)
	}
}

func LessEqual(_ *Environment, args []Node) Node {
	if len(args) != 2 {
		return ErrorNode{fmt.Errorf("expected 1 argument, got %v", len(args))}
	}

	n1, ok := args[0].(NumberNode)
	if !ok {
		return ErrorNode{IncorrectType{"Number", args[0].TypeString()}}
	}

	n2, ok := args[1].(NumberNode)
	if !ok {
		return ErrorNode{IncorrectType{"Number", args[1].TypeString()}}
	}

	if n1 <= n2 {
		return NumberNode(1)
	} else {
		return NumberNode(0)
	}
}

func More(_ *Environment, args []Node) Node {
	if len(args) != 2 {
		return ErrorNode{fmt.Errorf("expected 1 argument, got %v", len(args))}
	}

	n1, ok := args[0].(NumberNode)
	if !ok {
		return ErrorNode{IncorrectType{"Number", args[0].TypeString()}}
	}

	n2, ok := args[1].(NumberNode)
	if !ok {
		return ErrorNode{IncorrectType{"Number", args[1].TypeString()}}
	}

	if n1 > n2 {
		return NumberNode(1)
	} else {
		return NumberNode(0)
	}
}

func MoreEqual(_ *Environment, args []Node) Node {
	if len(args) != 2 {
		return ErrorNode{fmt.Errorf("expected 1 argument, got %v", len(args))}
	}

	n1, ok := args[0].(NumberNode)
	if !ok {
		return ErrorNode{IncorrectType{"Number", args[0].TypeString()}}
	}

	n2, ok := args[1].(NumberNode)
	if !ok {
		return ErrorNode{IncorrectType{"Number", args[1].TypeString()}}
	}

	if n1 >= n2 {
		return NumberNode(1)
	} else {
		return NumberNode(0)
	}
}
