package goculator

import (
	"errors"
	"fmt"
	"strconv"
)

type Calculator struct {
	input   string
	lexer   *Lexer
	context Context
}

func NewCalculator(input string) *Calculator {
	interpret := new(Calculator)
	interpret.input = input
	lexer := NewLexer(input)
	interpret.lexer = lexer
	interpret.lexer.Scan()
	return interpret
}

func (t *Calculator) SetContext(c Context) {
	t.context = c
}

func (t *Calculator) Go() (float64, error) {
	return t.expr()
}

func (t *Calculator) eat(tokenType TokenType) error {
	if t.currentToken().Type != tokenType {
		return errors.New(
			fmt.Sprintf(
				"expected token type %s is not matching currunt token type %s",
				tokenType,
				t.lexer.Token().Type,
			),
		)
	}
	t.lexer.Scan()
	return t.lexer.Err()
}

func (t *Calculator) currentToken() Token {
	return t.lexer.Token()
}

func (t *Calculator) value(key string) (float64, error) {
	if t.context == nil {
		return 0, errors.New("no context given for variable")
	}

	return t.context.Value(key)
}

func (t *Calculator) term() (float64, error) {

	token := t.currentToken()

	// For parantheses case
	if token.Type == TokenTypeLPARAN {
		if err := t.eat(TokenTypeLPARAN); err != nil {
			return 0, err
		}
		result, err := t.expr()
		if err != nil {
			return 0, err
		}
		if err := t.eat(TokenTypeRPARAN); err != nil {
			return 0, err
		}
		return result, nil
	}

	// For variable case
	if token.Type == TokenTypeVAR {
		if err := t.eat(TokenTypeVAR); err != nil {
			return 0, err
		}

		value, err := t.value(token.Value)
		if err != nil {
			return 0, err
		}
		return value, nil
	}

	// For number case
	if err := t.eat(TokenTypeNUM); err != nil {
		return 0, err
	}

	result, err := strconv.ParseFloat(token.Value, 64)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (t *Calculator) factor() (float64, error) {
	result, err := t.term()
	if err != nil {
		return 0, err
	}

	for t.currentToken().IsMultiDiv() {
		op := t.currentToken()
		switch op.Type {
		case TokenTypeMULTI:
			if err := t.eat(TokenTypeMULTI); err != nil {
				return 0, err
			}
		case TokenTypeDIV:
			if err := t.eat(TokenTypeDIV); err != nil {
				return 0, err
			}
		}

		num, err := t.term()

		if err != nil {
			return 0, err
		}

		switch op.Type {
		case TokenTypeMULTI:
			result = result * num
		case TokenTypeDIV:
			result = result / num
		}
	}

	return result, nil
}

func (t *Calculator) expr() (float64, error) {
	if t.currentToken().Type == TokenTypeEOF {
		return 0, nil
	}

	result, err := t.factor()
	if err != nil {
		return 0, err
	}

	for t.currentToken().IsPlusMinus() {
		op := t.currentToken()
		switch op.Type {
		case TokenTypePLUS:
			if err := t.eat(TokenTypePLUS); err != nil {
				return 0, err
			}
		case TokenTypeMINUS:
			if err := t.eat(TokenTypeMINUS); err != nil {
				return 0, err
			}
		}

		num, err := t.factor()

		if err != nil {
			return 0, err
		}

		switch op.Type {
		case TokenTypePLUS:
			result = result + num
		case TokenTypeMINUS:
			result = result - num
		}
	}

	return result, nil
}
