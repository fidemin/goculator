package goculator

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestCalculator(t *testing.T) {
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
		calc := New(data.input)
		calc.SetContext(context)

		result, err := calc.Go()

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
