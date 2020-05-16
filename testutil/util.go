package testutil

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// PlayTest interface for play test
type PlayTest interface {
	Match(args ...interface{})
}

// Play play func for impl
func Play(t *testing.T, impl interface{}, funcName string, args ...interface{}) PlayTest {
	f, ok := reflect.TypeOf(impl).MethodByName(funcName)
	if !ok {
		t.Errorf("%v func not exist ", funcName)
	}

	var inParam []reflect.Value
	inParam = append(inParam, reflect.ValueOf(impl))
	for _, v := range args {
		inParam = append(inParam, reflect.ValueOf(v))
	}
	res := f.Func.Call(inParam)

	return &playTestImpl{
		t:           t,
		actualValue: res,
	}
}

// playTestImpl play test
type playTestImpl struct {
	t           *testing.T
	actualValue []reflect.Value
	expectValue []reflect.Value
}

func (p *playTestImpl) Match(args ...interface{}) {
	if len(p.actualValue) != len(args) {
		p.t.Errorf("actual value length is %v , expected value length is %v  , not matched", len(p.actualValue), len(args))
		return
	}
	for k, v := range p.actualValue {
		if !v.IsZero() {
			if !assert.Equal(p.t, args[k], v.Interface(), "index ", k, "not matched") {
				p.t.Fail()
			}
		} else {
			if !assert.Equal(p.t, args[k], nil, "index ", k, "not matched") {
				p.t.Fail()
			}
		}

	}
}
