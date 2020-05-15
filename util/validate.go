package util

import (
	"reflect"
)

// ValueStrategy interface
type ValueStrategy interface {
	Load(args ...interface{})
	HasValue() bool
	ZeroIndex() int
}

// ValueValidationImpl value validation imple
type ValueValidationImpl struct {
	args      []interface{}
	zeroIndex int
}

// Load load args
func (v *ValueValidationImpl) Load(args ...interface{}) {
	v.args = args
}

// IfHasNilValue  check all args if has value  , if all has value , return true , elase return flase
func (v *ValueValidationImpl) IfHasNilValue() bool {
	var HasNilValue = false
loop:
	for i, arg := range v.args {
		// if has value
		if reflect.ValueOf(arg).IsZero() {
			HasNilValue = true
			v.zeroIndex = i
			break loop
		}
	}
	return HasNilValue
}

// ZeroIndex  return first zero index
func (v *ValueValidationImpl) ZeroIndex() int {
	return v.zeroIndex
}
