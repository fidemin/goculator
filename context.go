package goculator

import (
	"errors"
	"fmt"
)

type Context interface {
	Value(string) (float64, error)
}

type DefaultContext struct {
	keyValueMap map[string]float64
}

func NewDefaultContext(keyValues map[string]float64) *DefaultContext {
	c := new(DefaultContext)
	c.keyValueMap = keyValues
	return c
}

func (c *DefaultContext) Value(key string) (float64, error) {
	value, ok := c.keyValueMap[key]
	if !ok {
		return 0, errors.New(fmt.Sprintf("no value for key '%s'", key))
	}
	return value, nil
}
