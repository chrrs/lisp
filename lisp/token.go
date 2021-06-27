package lisp

import (
	"fmt"
	"regexp"
)

type TokenType uint8

func (t TokenType) String() string {
	switch t {
	case WhitespaceToken:
		return "Whitespace"
	case OpenToken:
		return "Open"
	case CloseToken:
		return "Close"
	case NumberToken:
		return "Number"
	case IdentifierToken:
		return "Identifier"
	default:
		return "<unknown>"
	}
}

const (
	WhitespaceToken TokenType = iota
	OpenToken
	CloseToken
	NumberToken
	StringToken
	IdentifierToken
)

var Patterns = map[TokenType]*regexp.Regexp{
	WhitespaceToken: regexp.MustCompile("^\\s+"),
	OpenToken:       regexp.MustCompile("^[({]"),
	CloseToken:      regexp.MustCompile("^[)}]"),
	NumberToken:     regexp.MustCompile("^[+-]?(\\d+\\.?\\d*|\\.\\d+)"),
	StringToken:     regexp.MustCompile("^\"(?:[^\\\\\"]|\\\\.)*\""),
	IdentifierToken: regexp.MustCompile("^[a-zA-Z_+\\-*/\\\\=<>!&][a-zA-Z0-9_+\\-*/\\\\=<>!&]*"),
}

type Token struct {
	Type  TokenType
	Value string
}

func (t Token) String() string {
	return fmt.Sprintf("%v(%v)", t.Type, t.Value)
}

type UnexpectedCharacter uint8

func (t UnexpectedCharacter) Error() string {
	return fmt.Sprintf("unexpected character in input: %c", uint8(t))
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
			return nil, UnexpectedCharacter(toParse[0])
		}
	}

	return tokens, nil
}
