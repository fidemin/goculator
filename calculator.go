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
	TokenTypeNUM   TokenType = "NUM"
	TokenTypePLUS  TokenType = "PLUS"
	TokenTypeMINUS TokenType = "MINUS"
	TokenTypeEOF   TokenType = "EOF"
	TokenTypeNONE  TokenType = "NONE"
)

type Lexer struct {
	text    string
	length  int
	current Token
	pos     int
	err     error
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
	if l.length <= l.pos {
		l.current = Token{TokenTypeEOF, ""}
		return false
	}

	next := l.text[l.pos : l.pos+1]
	if l.isStrInt(next) {
		l.current = Token{TokenTypeNUM, next}
		l.pos++
		return true
	}

	if next == "+" {
		l.current = Token{TokenTypePLUS, "+"}
		l.pos++
		return true
	}

	if next == "-" {
		l.current = Token{TokenTypeMINUS, "-"}
		l.pos++
		return true
	}

	l.err = errors.New(fmt.Sprintf("'%s' is not acceptable string for lexer", next))

	return false
}

func (l *Lexer) isStrInt(str string) bool {
	if _, err := strconv.Atoi(str); err != nil {
		return false
	}
	return true
}
