package goculator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLexer(t *testing.T) {
	assert := assert.New(t)
	var testdata = []struct {
		input  string
		result []token
	}{
		{
			"32+21-1 /13.2 *23",
			[]token{
				token{tokenTypeNUM, "32"},
				token{tokenTypePLUS, "+"},
				token{tokenTypeNUM, "21"},
				token{tokenTypeMINUS, "-"},
				token{tokenTypeNUM, "1"},
				token{tokenTypeDIV, "/"},
				token{tokenTypeNUM, "13.2"},
				token{tokenTypeMULTI, "*"},
				token{tokenTypeNUM, "23"},
			},
		},
		{
			"32+(21-1.11)",
			[]token{
				token{tokenTypeNUM, "32"},
				token{tokenTypePLUS, "+"},
				token{tokenTypeLPARAN, "("},
				token{tokenTypeNUM, "21"},
				token{tokenTypeMINUS, "-"},
				token{tokenTypeNUM, "1.11"},
				token{tokenTypeRPARAN, ")"},
			},
		},
		{
			"32+(21-var_1k)-1",
			[]token{
				token{tokenTypeNUM, "32"},
				token{tokenTypePLUS, "+"},
				token{tokenTypeLPARAN, "("},
				token{tokenTypeNUM, "21"},
				token{tokenTypeMINUS, "-"},
				token{tokenTypeVAR, "var_1k"},
				token{tokenTypeRPARAN, ")"},
				token{tokenTypeMINUS, "-"},
				token{tokenTypeNUM, "1"},
			},
		},
		{
			"",
			[]token{},
		},
	}

	for _, data := range testdata {
		lexer := newLexer(data.input)

		tokens := make([]token, 0)
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
