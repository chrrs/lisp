package lisp

import (
	"fmt"
	"strconv"
)

type Node interface {
	fmt.Stringer
}

type ExpressionNode struct {
	Nodes     []Node
}

func (e ExpressionNode) String() string {
	ret := "("

	for i, node := range e.Nodes {
		if i != 0 {
			ret += " "
		}

		ret += node.String()
	}

	ret += ")"
	return ret
}

type IdentifierNode struct {
	Name string
}

func (i IdentifierNode) String() string {
	return i.Name
}

type ValueNode struct {
	Value int
}

func (v ValueNode) String() string {
	return strconv.Itoa(v.Value)
}

type UnexpectedToken Token

func (t UnexpectedToken) Error() string {
	return fmt.Sprintf("unexpected token in input: %v", t.Value)
}

type UnexpectedEOI struct{}

func (_ UnexpectedEOI) Error() string {
	return "unexpected end of input"
}

func trimWhitespace(input *[]Token) bool {
	if len(*input) == 0 {
		return false
	}

	if (*input)[0].Type == Whitespace {
		*input = (*input)[1:]
		return true
	} else {
		return false
	}
}

func findMatchingClose(input []Token) int {
	depth := 0

	for i, token := range input {
		switch token.Type {
		case Open:
			depth++
		case Close:
			depth--

			if depth == 0 {
				return i
			}
		}
	}

	return -1
}

func ParseExpression(input []Token) (ExpressionNode, error) {
	trimWhitespace(&input)

	ret := ExpressionNode{}

	for {
		if len(input) == 0 {
			return ret, nil
		}

		switch input[0].Type {
		case Identifier:
			ret.Nodes = append(ret.Nodes, IdentifierNode{input[0].Value})
			input = input[1:]
		case Number:
			value, _ := strconv.Atoi(input[0].Value)
			ret.Nodes = append(ret.Nodes, ValueNode{value})

			input = input[1:]
		case Open:
			closeIndex := findMatchingClose(input)
			if closeIndex == -1 {
				return ExpressionNode{}, UnexpectedEOI{}
			}

			nestedExpression, err := ParseExpression(input[1:closeIndex])
			if err != nil {
				return ExpressionNode{}, err
			}

			ret.Nodes = append(ret.Nodes, nestedExpression)

			input = input[closeIndex+1:]
		default:
			return ExpressionNode{}, UnexpectedToken(input[0])
		}

		if !trimWhitespace(&input) {
			if len(input) != 0 {
				return ExpressionNode{}, UnexpectedToken(input[0])
			}

			return ret, nil
		}
	}
}
