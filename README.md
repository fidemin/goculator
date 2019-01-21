# goculator

Arithmetic calculator for Golang

## Simple Use

It is very simple to use. 

``New`` function accepts arithmatic formular string, and returns ``Calculator`` instance. ``Go`` method calculates the formula and returns the result as float64.

### Example
```go
package main

import (
	"log"

	"github.com/yhmin84/goculator"
)

func main() {
    input := "(32 + 34) / 11"

	// returns Calculator instance
    calc := goculator.New(input)
	// calculate
    result, err := calc.Go()

    if err != nil {
        fmt.Println(err.Error())
        return
    }

    fmt.Println(result)

}
```

The result will be printed as below.

```
6
```

## Variable Context

Another feature is variable context interface, ``Context``, which has only one method ``Value``.

```go
type Context interface {
	Value(string) (float64, error)
}
```

The argument of the method is a variable and it returns the float64 value which matches the variable. The ``Calculator`` instance has SetContext method which accepts the interface as an argument. If the variable exists in the formula, it will be replaced with the result of ``Value`` method.

### Example
```go
package main 

import (
    "errors"
    "fmt"

    "github.com/yhmin84/goculator"
)

type MockContext struct{}

func (c *MockContext) Value(variable string) (float64, error) {
    keyValue := map[string]float64{
        "var": 34,
    }

    result, ok := keyValue[key]

    if !ok {
        return float64(0), errors.New("no value for a variable")
    }

    return result, nil
}

func main() {
    input := "(32 + var) / 11"
    calc := goculator.New(input)
    calc.SetContext(new(MockContext))
    result, err := calc.Go()

    if err != nil {
        fmt.Println(err.Error())
        return
    }

    fmt.Println(result)

}
```

``var`` in ``"(32 + var) / 11"`` is replaced with 34 and result will be printed as below.

```
6
```

For convenience, ``DefaultContext`` struct is given in package.

### Example
```go
package main

import (
    "fmt"

    "github.com/yhmin84/goculator"
)

func main() {
    keyValue := map[string]float64{
        "var": 34,
    }

    context := goculator.NewDefaultContext(keyValue)

    input := "(32 + var) / 11"
    calc := goculator.New(input)
    calc.SetContext(context)
    result, err := calc.Go()

    if err != nil {
        fmt.Println(err.Error())
        return
    }

    fmt.Println(result)

}
```

The result will be ``6``.
