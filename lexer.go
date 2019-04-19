package goculator

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var charToTokenType = map[string]TokenType{
	"+": TokenTypePLUS,
	"-": TokenTypeMINUS,
	"*": TokenTypeMULTI,
	"/": TokenTypeDIV,
	"(": TokenTypeLPARAN,
	")": TokenTypeRPARAN,
}

// Lexer scans input text to Token.
type Lexer struct {
	text        string
	length      int
	current     Token
	currentChar string
	pos         int
	err         error
}

// Lexer returns new Lexer with input text
func NewLexer(text string) *Lexer {
	lexer := new(Lexer)
	lexer.text = text
	lexer.length = len(text)
	lexer.current = Token{}

	return lexer
}

// Token returns Token of currently parsed characters
func (l *Lexer) Token() Token {
	return l.current
}

// Err returns error during lexical analysis.
// If Scan method returns false, error should be checked using Err method.
func (l *Lexer) Err() error {
	return l.err
}

// Scan returns true if converting of current characters to Token is successful.
// Scan returns false if EOF is encountered or error happens while scanning.
func (l *Lexer) Scan() bool {
	// l.text could be ''. EOF check should be done first.
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
		l.current = Token{charToTokenType[l.currentChar], l.currentChar}
		l.advance()
		return true
	}

	// unacceptable charater is included in text
	l.err = errors.New(fmt.Sprintf("'%s' is not acceptable character for Lexer", l.currentChar))

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
