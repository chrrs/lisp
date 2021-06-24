package lisp

import (
	"errors"
	"fmt"
	"strconv"
)

type Node interface {
	fmt.Stringer
	Evaluate() Node
}

type ExpressionType uint8

const (
	SExpression ExpressionType = iota
	QExpression
)

type ExpressionNode struct {
	Type  ExpressionType
	Nodes []Node
}

func (e ExpressionNode) String() string {
	ret := ""

	switch e.Type {
	case SExpression:
		ret += "("
	case QExpression:
		ret += "{"
	}

	for i, node := range e.Nodes {
		if i != 0 {
			ret += " "
		}

		ret += node.String()
	}

	switch e.Type {
	case SExpression:
		ret += ")"
	case QExpression:
		ret += "}"
	}

	return ret
}

type IdentifierNode string

func (i IdentifierNode) String() string {
	return string(i)
}

type NumberNode int

func (v NumberNode) String() string {
	return strconv.Itoa(int(v))
}

type ErrorNode struct {
	Error error
}

func (e ErrorNode) String() string {
	return "runtime error: " + e.Error.Error()
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

	if (*input)[0].Type == WhitespaceToken {
		*input = (*input)[1:]
		return true
	} else {
		return false
	}
}

func findMatchingClose(input []Token, type_ ExpressionType) (int, error) {
	depth := 0

	for i, token := range input {
		switch token.Type {
		case OpenToken:
			depth++
		case CloseToken:
			depth--

			if depth == 0 {
				switch type_ {
				case SExpression:
					if token.Value != ")" {
						return 0, UnexpectedToken(token)
					}
				case QExpression:
					if token.Value != "}" {
						return 0, UnexpectedToken(token)
					}
				}

				return i, nil
			}
		}
	}

	return 0, UnexpectedEOI{}
}

func ParseExpression(input []Token, type_ ExpressionType) (ExpressionNode, error) {
	trimWhitespace(&input)

	ret := ExpressionNode{Type: type_}

	for {
		if len(input) == 0 {
			return ret, nil
		}

		switch input[0].Type {
		case IdentifierToken:
			ret.Nodes = append(ret.Nodes, IdentifierNode(input[0].Value))
			input = input[1:]
		case NumberToken:
			value, _ := strconv.Atoi(input[0].Value)
			ret.Nodes = append(ret.Nodes, NumberNode(value))

			input = input[1:]
		case OpenToken:
			var type_ ExpressionType
			switch input[0].Value {
			case "(":
				type_ = SExpression
			case "{":
				type_ = QExpression
			default:
				return ExpressionNode{}, errors.New("unknown expression type for open bracket " + input[0].Value)
			}

			closeIndex, err := findMatchingClose(input, type_)
			if err != nil {
				return ExpressionNode{}, err
			}

			nestedExpression, err := ParseExpression(input[1:closeIndex], type_)
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
