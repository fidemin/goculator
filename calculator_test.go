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
			"3+2-1",
			[]Token{
				Token{TokenTypeNUM, "3"},
				Token{TokenTypePLUS, "+"},
				Token{TokenTypeNUM, "2"},
				Token{TokenTypeMINUS, "-"},
				Token{TokenTypeNUM, "1"},
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
