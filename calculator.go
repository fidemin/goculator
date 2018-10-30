package main

import (
	"errors"
	"fmt"
	"strconv"
)

type Token struct {
	Type  TokenType
	Value string
}

type TokenType string

const (
	TokenTypeNUM    TokenType = "NUM"
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

	if l.isSpace(l.currentChar) {
		l.skipSpace()
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

func (l *Lexer) isEOF() bool {
	return l.length <= l.pos
}

func (l *Lexer) skipSpace() {
	for !l.isEOF() && l.isSpace(l.currentChar) {
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

func (l *Lexer) advance() {
	l.pos++
	if !l.isEOF() {
		l.currentChar = l.text[l.pos : l.pos+1]
	}
}

func (l *Lexer) isSpace(str string) bool {
	return str == " "
}

func (l *Lexer) isIntOrDot() bool {
	if l.currentChar == "." {
		return true
	}

	if _, err := strconv.Atoi(l.currentChar); err != nil {
		return false
	}
	return true
}
