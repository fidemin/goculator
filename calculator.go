package goculator

import (
	"errors"
	"fmt"
	"strconv"
)

type Calculator struct {
	input   string
	lexer   *lexer
	context Context
}

func New(input string) *Calculator {
	interpret := new(Calculator)
	interpret.input = input
	lexer := newLexer(input)
	interpret.lexer = lexer
	interpret.lexer.Scan()
	return interpret
}

func (t *Calculator) Bind(c Context) {
	t.context = c
}

func (t *Calculator) Go() (float64, error) {
	return t.expr()
}

func (t *Calculator) eat(tokenType tokenType) error {
	if t.currenttoken().Type != tokenType {
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

func (t *Calculator) currenttoken() token {
	return t.lexer.Token()
}

func (t *Calculator) value(key string) (float64, error) {
	if t.context == nil {
		return 0, errors.New("no context given for variable")
	}

	return t.context.Value(key)
}

func (t *Calculator) term() (float64, error) {

	token := t.currenttoken()

	// For parantheses case
	if token.Type == tokenTypeLPARAN {
		if err := t.eat(tokenTypeLPARAN); err != nil {
			return 0, err
		}
		result, err := t.expr()
		if err != nil {
			return 0, err
		}
		if err := t.eat(tokenTypeRPARAN); err != nil {
			return 0, err
		}
		return result, nil
	}

	// For variable case
	if token.Type == tokenTypeVAR {
		if err := t.eat(tokenTypeVAR); err != nil {
			return 0, err
		}

		value, err := t.value(token.Value)
		if err != nil {
			return 0, err
		}
		return value, nil
	}

	// For number case
	if err := t.eat(tokenTypeNUM); err != nil {
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

	for t.currenttoken().IsMultiDiv() {
		op := t.currenttoken()
		switch op.Type {
		case tokenTypeMULTI:
			if err := t.eat(tokenTypeMULTI); err != nil {
				return 0, err
			}
		case tokenTypeDIV:
			if err := t.eat(tokenTypeDIV); err != nil {
				return 0, err
			}
		}

		num, err := t.term()

		if err != nil {
			return 0, err
		}

		switch op.Type {
		case tokenTypeMULTI:
			result = result * num
		case tokenTypeDIV:
			result = result / num
		}
	}

	return result, nil
}

func (t *Calculator) expr() (float64, error) {
	if t.currenttoken().Type == tokenTypeEOF {
		return 0, nil
	}

	result, err := t.factor()
	if err != nil {
		return 0, err
	}

	for t.currenttoken().IsPlusMinus() {
		op := t.currenttoken()
		switch op.Type {
		case tokenTypePLUS:
			if err := t.eat(tokenTypePLUS); err != nil {
				return 0, err
			}
		case tokenTypeMINUS:
			if err := t.eat(tokenTypeMINUS); err != nil {
				return 0, err
			}
		}

		num, err := t.factor()

		if err != nil {
			return 0, err
		}

		switch op.Type {
		case tokenTypePLUS:
			result = result + num
		case tokenTypeMINUS:
			result = result - num
		}
	}

	return result, nil
}
