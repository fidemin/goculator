package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLexer(t *testing.T) {
	assert := assert.New(t)
	var testdata = []struct {
		input  string
		result []Token
	}{
		{
			"32+21-1 /13.2 *23",
			[]Token{
				Token{TokenTypeNUM, "32"},
				Token{TokenTypePLUS, "+"},
				Token{TokenTypeNUM, "21"},
				Token{TokenTypeMINUS, "-"},
				Token{TokenTypeNUM, "1"},
				Token{TokenTypeDIV, "/"},
				Token{TokenTypeNUM, "13.2"},
				Token{TokenTypeMULTI, "*"},
				Token{TokenTypeNUM, "23"},
			},
		},
		{
			"32+(21-1.11)",
			[]Token{
				Token{TokenTypeNUM, "32"},
				Token{TokenTypePLUS, "+"},
				Token{TokenTypeLPARAN, "("},
				Token{TokenTypeNUM, "21"},
				Token{TokenTypeMINUS, "-"},
				Token{TokenTypeNUM, "1.11"},
				Token{TokenTypeRPARAN, ")"},
			},
		},
		{
			"",
			[]Token{},
		},
	}

	for _, data := range testdata {
		lexer := NewLexer(data.input)

		tokens := make([]Token, 0)
		for lexer.Scan() {
			token := lexer.Token()
			tokens = append(tokens, token)
		}

		if err := lexer.Err(); err != nil {
			assert.Fail(err.Error())
			return
		}

		for i, token := range data.result {
			assert.Equal(token.Type, tokens[i].Type)
			assert.Equal(token.Value, tokens[i].Value)
		}
	}
}
