package lisp

import (
	"fmt"
	"regexp"
)

type TokenType uint8

func (t TokenType) String() string {
	switch t {
	case Whitespace:
		return "Whitespace"
	case Open:
		return "Open"
	case Close:
		return "Close"
	case Number:
		return "Number"
	case Identifier:
		return "Identifier"
	default:
		return "<unknown>"
	}
}

const (
	Whitespace TokenType = iota
	Open
	Close
	Number
	Identifier
)

var Patterns = map[TokenType]*regexp.Regexp{
	Whitespace: regexp.MustCompile("^\\s+"),
	Open: regexp.MustCompile("^\\("),
	Close: regexp.MustCompile("^\\)"),
	Number: regexp.MustCompile("^\\d+"),
	Identifier: regexp.MustCompile("^[+\\-/*]"),
}

type Token struct {
	Type  TokenType
	Value string
}

type UnexpectedToken uint8

func (t UnexpectedToken) Error() string {
	return fmt.Sprintf("unexpected token in input: %c", uint8(t))
}

func Tokenize(input string) ([]Token, error) {
	toParse := input
	tokens := make([]Token, 0)

	for len(toParse) > 0 {
		parsed := false

		for type_, pattern := range Patterns {
			matches := pattern.FindStringSubmatch(toParse)
			if len(matches) > 0 {
				parsed = true
				tokens = append(tokens, Token{type_, matches[0]})
				toParse = toParse[len(matches[0]):]
				break
			}
		}

		if !parsed {
			return nil, UnexpectedToken(toParse[0])
		}
	}

	return tokens, nil
}