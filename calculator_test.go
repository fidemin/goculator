package main

import (
	"github.com/stretchr/testify/assert"
	"math"
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
			"32+(21-var_1k)-1",
			[]Token{
				Token{TokenTypeNUM, "32"},
				Token{TokenTypePLUS, "+"},
				Token{TokenTypeLPARAN, "("},
				Token{TokenTypeNUM, "21"},
				Token{TokenTypeMINUS, "-"},
				Token{TokenTypeVAR, "var_1k"},
				Token{TokenTypeRPARAN, ")"},
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

func TestInterpreter(t *testing.T) {
	assert := assert.New(t)
	var testdata = []struct {
		input  string
		result float64
	}{
		{
			"32+21.1-21",
			32.1,
		},
		{

			"2.1-2*4/2+1",
			-0.9,
		},
		{
			"2.1/(2.1+4.2)",
			0.33,
		},
		{
			"2.1/0",
			math.Inf(0),
		},
		{
			"0/0",
			math.NaN(),
		},
		{
			"2.1/(var1 + var2)",
			0.33,
		},
		{
			"",
			0,
		},
	}

	context := NewDefaultContext(
		map[string]float64{
			"var1": 2.1,
			"var2": 4.2,
		},
	)

	for _, data := range testdata {
		inter := NewInterpreter(data.input)
		inter.SetContext(context)

		result, err := inter.Interpret()

		if err != nil {
			assert.Fail(err.Error())
			return
		}

		if math.IsNaN(data.result) {
			assert.True(math.IsNaN(result))
		} else {
			assert.InDelta(data.result, result, 0.01)
		}

	}
}
