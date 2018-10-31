package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type Token struct {
	Type  TokenType
	Value string
}

func (t Token) IsPlusMinus() bool {
	switch t.Type {
	case TokenTypePLUS, TokenTypeMINUS:
		return true
	}
	return false
}

func (t Token) IsMultiDiv() bool {
	switch t.Type {
	case TokenTypeMULTI, TokenTypeDIV:
		return true
	}
	return false
}

type TokenType string

const (
	TokenTypeNUM    TokenType = "NUM"
	TokenTypeVAR    TokenType = "VAR"
	TokenTypePLUS   TokenType = "PLUS"
	TokenTypeMINUS  TokenType = "MINUS"
	TokenTypeMULTI  TokenType = "MULTI"
	TokenTypeDIV    TokenType = "DIV"
	TokenTypeEOF    TokenType = "EOF"
	TokenTypeNONE   TokenType = "NONE"
	TokenTypeLPARAN TokenType = "LPARAN"
	TokenTypeRPARAN TokenType = "RPARAN"
)

var tokenTypeMap = map[string]TokenType{
	"+": TokenTypePLUS,
	"-": TokenTypeMINUS,
	"*": TokenTypeMULTI,
	"/": TokenTypeDIV,
	"(": TokenTypeLPARAN,
	")": TokenTypeRPARAN,
}

type Lexer struct {
	text        string
	length      int
	current     Token
	currentChar string
	pos         int
	err         error
}

func NewLexer(text string) *Lexer {
	lexer := new(Lexer)
	lexer.text = text
	lexer.length = len(text)
	lexer.current = Token{}
	return lexer
}

func (l *Lexer) Token() Token {
	return l.current
}

func (l *Lexer) Err() error {
	return l.err
}

func (l *Lexer) Scan() bool {
	if l.isEOF() {
		l.current = Token{TokenTypeEOF, ""}
		return false
	}

	if l.pos == 0 {
		l.currentChar = l.text[l.pos : l.pos+1]
	}

	if l.isSpace() {
		l.skipSpace()
	}

	if l.isStr() {
		l.current = Token{TokenTypeVAR, l.variable()}
		return true
	}

	if l.isIntOrDot() {
		l.current = Token{TokenTypeNUM, l.number()}
		return true
	}

	switch l.currentChar {
	case "+", "-", "*", "/", ")", "(":
		l.current = Token{tokenTypeMap[l.currentChar], l.currentChar}
		l.advance()
		return true
	}

	l.err = errors.New(fmt.Sprintf("'%s' is not acceptable string for lexer", l.currentChar))

	return false
}

func (l *Lexer) skipSpace() {
	for !l.isEOF() && l.isSpace() {
		l.advance()
	}
}

func (l *Lexer) number() string {
	number := ""
	for !l.isEOF() && l.isIntOrDot() {
		number += l.currentChar
		l.advance()
	}
	return number
}

func (l *Lexer) variable() string {
	variable := ""

	// First character should be alphabet or _
	if !l.isEOF() && l.isStr() {
		variable += l.currentChar
		l.advance()
	}

	for !l.isEOF() && (l.isStr() || l.isInt()) {
		variable += l.currentChar
		l.advance()
	}
	return variable
}

func (l *Lexer) advance() {
	l.pos++
	if !l.isEOF() {
		l.currentChar = l.text[l.pos : l.pos+1]
	}
}

func (l *Lexer) isEOF() bool {
	return l.length <= l.pos
}

func (l *Lexer) isSpace() bool {
	return l.currentChar == " "
}

func (l *Lexer) isIntOrDot() bool {
	if l.currentChar == "." {
		return true
	}
	return l.isInt()
}

func (l *Lexer) isInt() bool {
	if _, err := strconv.Atoi(l.currentChar); err != nil {
		return false
	}
	return true
}

func (l *Lexer) isStr() bool {
	isStr, err := regexp.MatchString("^[a-zA-Z_]$", l.currentChar)
	if err != nil {
		return false
	}
	return isStr
}

type Interpreter struct {
	input   string
	lexer   *Lexer
	context Context
}

func NewInterpreter(input string) *Interpreter {
	interpret := new(Interpreter)
	interpret.input = input
	lexer := NewLexer(input)
	interpret.lexer = lexer
	interpret.lexer.Scan()
	return interpret
}

func (t *Interpreter) SetContext(c Context) {
	t.context = c
}

func (t *Interpreter) Interpret() (float64, error) {
	return t.expr()
}

func (t *Interpreter) eat(tokenType TokenType) error {
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

func (t *Interpreter) currentToken() Token {
	return t.lexer.Token()
}

func (t *Interpreter) value(key string) (float64, error) {
	if t.context == nil {
		return 0, errors.New("no context given for variable")
	}

	return t.context.Value(key)
}

func (t *Interpreter) term() (float64, error) {

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

func (t *Interpreter) factor() (float64, error) {
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

func (t *Interpreter) expr() (float64, error) {
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

type Context interface {
	Value(string) (float64, error)
}

type DefaultContext struct {
	keyValues map[string]float64
}

func NewDefaultContext(keyValues map[string]float64) *DefaultContext {
	c := new(DefaultContext)
	c.keyValues = keyValues
	return c
}

func (c *DefaultContext) Value(key string) (float64, error) {
	value, ok := c.keyValues[key]
	if !ok {
		return 0, errors.New(fmt.Sprintf("no value for key '%s'", key))
	}
	return value, nil
}
