package lisp

import "errors"

func (e ExpressionNode) EvalAsSExpr() Node {
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

	switch op {
	case "+":
		return Add(args)
	case "-":
		return Sub(args)
	case "*":
		return Mul(args)
	case "/":
		return Div(args)
	default:
		return ErrorNode{errors.New("unknown identifier " + string(op))}
	}
}

func (e ExpressionNode) Evaluate() Node {
	if e.Type == QExpression {
		return e
	}

	return e.EvalAsSExpr()
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
