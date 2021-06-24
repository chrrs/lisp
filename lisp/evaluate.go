package lisp

import "errors"

func (e ExpressionNode) Evaluate() Node {
	if len(e.Nodes) == 0 {
		return e
	}

	nodes := make([]Node, len(e.Nodes))
	for i, node := range e.Nodes {
		_, ok := node.(ErrorNode)
		if ok {
			return node
		}

		nodes[i] = node.Evaluate()
	}

	if len(nodes) == 1 {
		return nodes[0]
	}

	opNode := nodes[0]
	args := nodes[1:]

	op, ok := opNode.(IdentifierNode)
	if !ok {
		return ErrorNode{errors.New("s-expressions should start with an identifier")}
	}

	for _, arg := range args {
		_, ok := arg.(NumberNode)
		if !ok {
			return ErrorNode{errors.New("cannot operate on a non-number")}
		}
	}

	if len(args) == 1 && op == "-" {
		return -args[0].(NumberNode)
	}

	ret := args[0].(NumberNode)

	for _, arg := range args[1:] {
		switch op {
		case "+":
			ret += arg.(NumberNode)
		case "-":
			ret -= arg.(NumberNode)
		case "*":
			ret *= arg.(NumberNode)
		case "/":
			if arg.(NumberNode) == 0 {
				return ErrorNode{errors.New("divide by zero")}
			}

			ret /= arg.(NumberNode)
		}
	}

	return ret
}

func (i IdentifierNode) Evaluate() Node {
	return i
}

func (v NumberNode) Evaluate() Node {
	return v
}

func (e ErrorNode) Evaluate() Node {
	return e
}
