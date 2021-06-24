package lisp

import "errors"

type Environment map[IdentifierNode]Node

func (env Environment) AddBuiltins() {
	env["+"] = FunctionNode{Add}
	env["-"] = FunctionNode{Sub}
	env["*"] = FunctionNode{Mul}
	env["/"] = FunctionNode{Div}

	env["head"] = FunctionNode{Head}
	env["tail"] = FunctionNode{Tail}
	env["list"] = FunctionNode{List}
	env["eval"] = FunctionNode{Eval}
	env["join"] = FunctionNode{Join}
}

func (e ExpressionNode) EvalAsSExpr(env Environment) Node {
	if len(e.Nodes) == 0 {
		return e
	}

	nodes := make([]Node, len(e.Nodes))
	for i, node := range e.Nodes {
		_, ok := node.(ErrorNode)
		if ok {
			return node
		}

		nodes[i] = node.Evaluate(env)
	}

	if len(nodes) == 1 {
		return nodes[0]
	}

	op := nodes[0]
	args := nodes[1:]

	fun, ok := op.(FunctionNode)
	if !ok {
		return ErrorNode{errors.New("S-Expressions should start with an function")}
	}

	return fun.Builtin(env, args)
}

func (e ExpressionNode) Evaluate(env Environment) Node {
	if e.Type == QExpression {
		return e
	}

	return e.EvalAsSExpr(env)
}

func (i IdentifierNode) Evaluate(env Environment) Node {
	node, ok := env[i]
	if !ok {
		return ErrorNode{errors.New("unknown identifier " + string(i))}
	}

	return node
}

func (v NumberNode) Evaluate(_ Environment) Node {
	return v
}

func (e ErrorNode) Evaluate(_ Environment) Node {
	return e
}

func (f FunctionNode) Evaluate(_ Environment) Node {
	return f
}
