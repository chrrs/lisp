package lisp

import (
	"fmt"
	"strconv"
	"strings"
)

type Node interface {
	Dump(depth int)
}

type ExpressionNode struct {
	Operation string
	Nodes     []Node
}

func (e ExpressionNode) Dump(depth int) {
	fmt.Println(strings.Repeat("  ", depth), "Expression", e.Operation)
	for _, n := range e.Nodes {
		n.Dump(depth + 1)
	}
}

type ValueNode struct {
	Value int
}

func (v ValueNode) Dump(depth int) {
	fmt.Println(strings.Repeat("  ", depth), "Value", v.Value)
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

func findLast(input []Token, type_ TokenType) int {
	for i := len(input) - 1; i >= 0; i-- {
		if input[i].Type == type_ {
			return i
		}
	}

	return -1
}

func ParseExpression(input []Token) (ExpressionNode, error) {
	trimWhitespace(&input)

	if len(input) == 0 {
		return ExpressionNode{}, UnexpectedEOI{}
	}

	if input[0].Type != Identifier {
		return ExpressionNode{}, UnexpectedToken(input[0])
	}

	ret := ExpressionNode{}
	ret.Operation, input = input[0].Value, input[1:]

	for trimWhitespace(&input) {
		if len(input) == 0 {
			break
		}

		switch input[0].Type {
		case Number:
			value, _ := strconv.Atoi(input[0].Value)
			ret.Nodes = append(ret.Nodes, ValueNode{value})

			input = input[1:]
		case Open:
			closeIndex := findLast(input, Close)
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
	}

	return ret, nil
}
