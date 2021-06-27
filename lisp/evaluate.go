package lisp

import (
	"errors"
	"fmt"
)

type Environment struct {
	Parent *Environment
	values map[IdentifierNode]Node
}

func NewEnvironment(parent *Environment) Environment {
	return Environment{parent, make(map[IdentifierNode]Node)}
}

func (env *Environment) Get(id IdentifierNode) Node {
	node, ok := env.values[id]
	if ok {
		return node
	}

	if env.Parent != nil {
		return env.Parent.Get(id)
	}

	return nil
}

func (env *Environment) Put(id IdentifierNode, value Node) {
	env.values[id] = value
}

func (env *Environment) Def(id IdentifierNode, value Node) {
	if env.Parent != nil {
		env.Parent.Def(id, value)
		return
	}

	env.Put(id, value)
}

func (env *Environment) AddBuiltins() {
	env.Def("+", FunctionNode{Builtin: Add})
	env.Def("-", FunctionNode{Builtin: Sub})
	env.Def("*", FunctionNode{Builtin: Mul})
	env.Def("/", FunctionNode{Builtin: Div})

	env.Def("=", FunctionNode{Builtin: Equal})

	env.Def("import", FunctionNode{Builtin: Import})
	env.Def("head", FunctionNode{Builtin: Head})
	env.Def("tail", FunctionNode{Builtin: Tail})
	env.Def("list", FunctionNode{Builtin: List})
	env.Def("eval", FunctionNode{Builtin: Eval})
	env.Def("join", FunctionNode{Builtin: Join})
	env.Def("def", FunctionNode{Builtin: Def})
	env.Def("let", FunctionNode{Builtin: Let})
	env.Def("fn", FunctionNode{Builtin: Fn})
}

func (e ExpressionNode) EvalAsSExpr(env *Environment) Node {
	if len(e.Nodes) == 0 {
		return e
	}

	nodes := make([]Node, len(e.Nodes))
	for i, node := range e.Nodes {
		evaluated := node.Evaluate(env)

		_, ok := evaluated.(ErrorNode)
		if ok {
			return evaluated
		}

		nodes[i] = evaluated
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

	return fun.call(env, args)
}

func (e ExpressionNode) Evaluate(env *Environment) Node {
	if e.Type == QExpression {
		return e
	}

	return e.EvalAsSExpr(env)
}

func (i IdentifierNode) Evaluate(env *Environment) Node {
	node := env.Get(i)
	if node == nil {
		return ErrorNode{errors.New("unknown identifier " + string(i))}
	}

	return node
}

func (v NumberNode) Evaluate(_ *Environment) Node {
	return v
}

func (s StringNode) Evaluate(_ *Environment) Node {
	return s
}

func (e ErrorNode) Evaluate(_ *Environment) Node {
	return e
}

func (f FunctionNode) call(env *Environment, args []Node) Node {
	if f.Builtin != nil {
		return f.Builtin(env, args)
	}

	formals := f.Formals

	for i, arg := range args {
		if len(formals) == 0 {
			return ErrorNode{fmt.Errorf("expected %v arguments, got %v", len(f.Formals), len(args))}
		}

		ident := formals[0]
		formals = formals[1:]

		if ident == "&" {
			if len(formals) != 1 {
				return ErrorNode{fmt.Errorf("expected 1 variadic argument, got %v", len(formals))}
			}

			ident = formals[0]
			formals = formals[:0]
			f.Environment.Put(ident, ExpressionNode{QExpression, args[i:]})

			break
		}

		f.Environment.Put(ident, arg)
	}

	if len(formals) == 2 && formals[0] == "&" {
		f.Environment.Put(formals[1], ExpressionNode{QExpression, make([]Node, 0)})
		formals = formals[:0]
	}

	if len(formals) == 0 {
		f.Environment.Parent = env
		return f.Body.EvalAsSExpr(f.Environment)
	}

	return FunctionNode{
		Builtin:     nil,
		Environment: f.Environment,
		Formals:     formals,
		Body:        f.Body,
	}

}

func (f FunctionNode) Evaluate(_ *Environment) Node {
	return f
}

func Evaluate(env *Environment, input string, multi bool) (Node, error) {
	tokens, err := Tokenize(input)
	if err != nil {
		return nil, fmt.Errorf("tokenization error: %v", err)
	}

	if multi {
		expression, err := ParseExpression(tokens, QExpression)
		if err != nil {
			return nil, fmt.Errorf("parsing error: %v", err)
		}

		for _, node := range expression.Nodes {
			out := node.Evaluate(env)
			err, ok := out.(ErrorNode)
			if ok {
				return nil, err.Error
			}
		}

		return ExpressionNode{Type: SExpression, Nodes: make([]Node, 0)}, nil
	} else {
		expression, err := ParseExpression(tokens, SExpression)
		if err != nil {
			return nil, fmt.Errorf("parsing error: %v", err)
		}

		return expression.Evaluate(env), nil
	}
}
