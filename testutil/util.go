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
// not support unexported func
func Play(t *testing.T, impl interface{}, funcName string, args ...interface{}) PlayTest {
	f, ok := reflect.TypeOf(impl).MethodByName(funcName)
	if !ok {
		t.Errorf("%v func not exist ", funcName)
		t.FailNow()
	}

	var inParam []reflect.Value
	inParam = append(inParam, reflect.ValueOf(impl))
	for _, v := range args {
		inParam = append(inParam, reflect.ValueOf(v))
	}
	res := f.Func.Call(inParam)
	var originValue []interface{}
	for _, v := range res {
		if v.IsZero() {
			originValue = append(originValue, nil)
		} else {
			originValue = append(originValue, v.Interface())
		}

	}
	return &playTestImpl{
		t:           t,
		originValue: originValue,
		actualValue: res,
	}
}

// PlayUnexported play func for impl for un exported func
func PlayUnexported(t *testing.T, args ...interface{}) PlayTest {
	actualValue := make([]reflect.Value, len(args))
	for k, v := range args {
		actualValue[k] = reflect.ValueOf(v)
	}
	return &playTestImpl{
		t:           t,
		originValue: args,
		actualValue: actualValue,
	}
}

// playTestImpl play test
type playTestImpl struct {
	t           *testing.T
	originValue []interface{}
	actualValue []reflect.Value
	expectValue []reflect.Value
}

func (p *playTestImpl) Match(args ...interface{}) {
	if len(p.actualValue) != len(args) {
		p.t.Errorf("actual value length is %v , expected value length is %v  , not matched", len(p.actualValue), len(args))
		p.t.FailNow()
	}
	for k, v := range p.actualValue {
		if args[k] == nil {
			if p.originValue[k] != nil {
				if !reflect.ValueOf(p.originValue[k]).IsNil() {
					if !assert.Equal(p.t, args[k], p.originValue[k], "index ", k, "not matched") {
						p.t.FailNow()
					}
				}
			} else {
				if !assert.Equal(p.t, args[k], p.originValue[k], "index ", k, "not matched") {
					p.t.FailNow()
				}
			}
		} else {
			if !assert.Equal(p.t, args[k], v.Interface(), "index ", k, "not matched") {
				p.t.FailNow()
			}
		}
	}
}
