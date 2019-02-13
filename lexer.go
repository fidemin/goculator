package goculator

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type token struct {
	Type  tokenType
	Value string
}

func (t token) IsPlusMinus() bool {
	switch t.Type {
	case tokenTypePLUS, tokenTypeMINUS:
		return true
	}
	return false
}

func (t token) IsMultiDiv() bool {
	switch t.Type {
	case tokenTypeMULTI, tokenTypeDIV:
		return true
	}
	return false
}

type tokenType string

const (
	tokenTypeNUM    tokenType = "NUM"
	tokenTypeVAR    tokenType = "VAR"
	tokenTypePLUS   tokenType = "PLUS"
	tokenTypeMINUS  tokenType = "MINUS"
	tokenTypeMULTI  tokenType = "MULTI"
	tokenTypeDIV    tokenType = "DIV"
	tokenTypeEOF    tokenType = "EOF"
	tokenTypeNONE   tokenType = "NONE"
	tokenTypeLPARAN tokenType = "LPARAN"
	tokenTypeRPARAN tokenType = "RPARAN"
)

var tokenTypeMap = map[string]tokenType{
	"+": tokenTypePLUS,
	"-": tokenTypeMINUS,
	"*": tokenTypeMULTI,
	"/": tokenTypeDIV,
	"(": tokenTypeLPARAN,
	")": tokenTypeRPARAN,
}

type lexer struct {
	text        string
	length      int
	current     token
	currentChar string
	pos         int
	err         error
}

func newLexer(text string) *lexer {
	lexer := new(lexer)
	lexer.text = text
	lexer.length = len(text)
	lexer.current = token{}
	return lexer
}

func (l *lexer) Token() token {
	return l.current
}

func (l *lexer) Err() error {
	return l.err
}

func (l *lexer) Scan() bool {
	if l.isEOF() {
		l.current = token{tokenTypeEOF, ""}
		return false
	}

	if l.pos == 0 {
		l.currentChar = l.text[l.pos : l.pos+1]
	}

	if l.isSpace() {
		l.skipSpace()
	}

	if l.isStr() {
		l.current = token{tokenTypeVAR, l.variable()}
		return true
	}

	if l.isIntOrDot() {
		l.current = token{tokenTypeNUM, l.number()}
		return true
	}

	switch l.currentChar {
	case "+", "-", "*", "/", ")", "(":
		l.current = token{tokenTypeMap[l.currentChar], l.currentChar}
		l.advance()
		return true
	}

	l.err = errors.New(fmt.Sprintf("'%s' is not acceptable string for lexer", l.currentChar))

	return false
}

func (l *lexer) skipSpace() {
	for !l.isEOF() && l.isSpace() {
		l.advance()
	}
}

func (l *lexer) number() string {
	number := ""
	for !l.isEOF() && l.isIntOrDot() {
		number += l.currentChar
		l.advance()
	}
	return number
}

func (l *lexer) variable() string {
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

func (l *lexer) advance() {
	l.pos++
	if !l.isEOF() {
		l.currentChar = l.text[l.pos : l.pos+1]
	}
}

func (l *lexer) isEOF() bool {
	return l.length <= l.pos
}

func (l *lexer) isSpace() bool {
	return l.currentChar == " "
}

func (l *lexer) isIntOrDot() bool {
	if l.currentChar == "." {
		return true
	}
	return l.isInt()
}

func (l *lexer) isInt() bool {
	if _, err := strconv.Atoi(l.currentChar); err != nil {
		return false
	}
	return true
}

func (l *lexer) isStr() bool {
	isStr, err := regexp.MatchString("^[a-zA-Z_]$", l.currentChar)
	if err != nil {
		return false
	}
	return isStr
}
