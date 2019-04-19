package goculator

import (
	"errors"
	"fmt"
	"strconv"
)

type Calculator struct {
	input   string
	Lexer   *Lexer
	context Context
}

func New(input string) *Calculator {
	interpret := new(Calculator)
	interpret.input = input
	Lexer := NewLexer(input)
	interpret.Lexer = Lexer
	interpret.Lexer.Scan()
	return interpret
}

func (c *Calculator) Bind(context Context) {
	c.context = context
}

func (c *Calculator) Go() (float64, error) {
	return c.expr()
}

func (c *Calculator) eat(TokenType TokenType) error {
	if c.currentToken().Type != TokenType {
		return errors.New(
			fmt.Sprintf(
				"expected token type %s is not matching currunt token type %s",
				TokenType,
				c.Lexer.Token().Type,
			),
		)
	}
	c.Lexer.Scan()
	return c.Lexer.Err()
}

func (c *Calculator) currentToken() Token {
	return c.Lexer.Token()
}

func (c *Calculator) value(key string) (float64, error) {
	if c.context == nil {
		return 0, errors.New("no context given for variable")
	}

	return c.context.Value(key)
}

func (c *Calculator) term() (float64, error) {

	token := c.currentToken()

	// For parantheses case
	if token.Type == TokenTypeLPARAN {
		if err := c.eat(TokenTypeLPARAN); err != nil {
			return 0, err
		}
		result, err := c.expr()
		if err != nil {
			return 0, err
		}
		if err := c.eat(TokenTypeRPARAN); err != nil {
			return 0, err
		}
		return result, nil
	}

	// For variable case
	if token.Type == TokenTypeVAR {
		if err := c.eat(TokenTypeVAR); err != nil {
			return 0, err
		}

		value, err := c.value(token.Value)
		if err != nil {
			return 0, err
		}
		return value, nil
	}

	// For number case
	if err := c.eat(TokenTypeNUM); err != nil {
		return 0, err
	}

	result, err := strconv.ParseFloat(token.Value, 64)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (c *Calculator) factor() (float64, error) {
	result, err := c.term()
	if err != nil {
		return 0, err
	}

	for c.isCurrentTokenMultiOrDiv() {
		op := c.currentToken()
		switch op.Type {
		case TokenTypeMULTI:
			if err := c.eat(TokenTypeMULTI); err != nil {
				return 0, err
			}
		case TokenTypeDIV:
			if err := c.eat(TokenTypeDIV); err != nil {
				return 0, err
			}
		}

		num, err := c.term()

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

func (c *Calculator) expr() (float64, error) {
	if c.currentToken().Type == TokenTypeEOF {
		return 0, nil
	}

	result, err := c.factor()
	if err != nil {
		return 0, err
	}

	for c.isCurrentTokenPlusOrMinus() {
		op := c.currentToken()
		switch op.Type {
		case TokenTypePLUS:
			if err := c.eat(TokenTypePLUS); err != nil {
				return 0, err
			}
		case TokenTypeMINUS:
			if err := c.eat(TokenTypeMINUS); err != nil {
				return 0, err
			}
		}

		num, err := c.factor()

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

func (c *Calculator) isCurrentTokenPlusOrMinus() bool {
	cTokenType := c.currentToken().Type
	if cTokenType == TokenTypePLUS || cTokenType == TokenTypeMINUS {
		return true
	}
	return false
}

func (c *Calculator) isCurrentTokenMultiOrDiv() bool {
	cTokenType := c.currentToken().Type
	if cTokenType == TokenTypeMULTI || cTokenType == TokenTypeDIV {
		return true
	}
	return false
}
