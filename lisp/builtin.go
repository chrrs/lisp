package lisp

import "errors"

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
