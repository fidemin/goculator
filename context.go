package goculator

import (
	"errors"
	"fmt"
)

//Stringer is implemented by any value that has a Value method, which returns value of variable name.
type Context interface {
	Value(string) (float64, error)
}

// DefaultContext is simple Context which represents variable, value relations using map[string]float64.
type DefaultContext struct {
	keyValueMap map[string]float64
}

// NewDefaultContext is new DefaultContext with keyValues which has key as variable name and value as variable value.
func NewDefaultContext(keyValues map[string]float64) *DefaultContext {
	c := new(DefaultContext)
	c.keyValueMap = keyValues
	return c
}

// Value returns float64 value from key.
func (c *DefaultContext) Value(key string) (float64, error) {
	value, ok := c.keyValueMap[key]
	if !ok {
		return 0, errors.New(fmt.Sprintf("no value for key '%s'", key))
	}
	return value, nil
}
