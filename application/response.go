package application

import (
	"github.com/PolarPanda611/trinitygo/httputils"
)

// RequestMap request map to register request
type RequestMap struct {
	Method     httputils.RequestMethod
	SubPath    string
	FuncName   string
	Validators []Validator
}

// NewRequestMapping request mapping
// @funcName if funcname is "" , trinitygoo will use the default http method name
// to find the method
// e.g : http method "GET" ==> find method "GET"
func NewRequestMapping(method httputils.RequestMethod, path string, funcName string, validators ...Validator) *RequestMap {
	return &RequestMap{
		Method:     method,
		SubPath:    path,
		FuncName:   funcName,
		Validators: validators,
	}
}

// Validator to validator request
type Validator func(Context)
