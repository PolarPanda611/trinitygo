package httputil

import "fmt"

var _ error = ResponseData{}

// ResponseData response data
type ResponseData struct {
	Status  int
	Result  interface{}
	Err     interface{}
	Runtime map[string]string
}

func (r ResponseData) Error() string {
	return fmt.Sprintf("%v", r.Err)
}

// RequestMethod Supported Request Method
type RequestMethod string

const (
	// GET http get
	GET RequestMethod = "GET"
	// HEAD http head
	HEAD RequestMethod = "HEAD"
	// POST http POST
	POST RequestMethod = "POST"
	// PUT http PUT
	PUT RequestMethod = "PUT"
	// PATCH http PATCH
	PATCH RequestMethod = "PATCH"
	// DELETE http DELETE
	DELETE RequestMethod = "DELETE"
	// OPTIONS http OPTIONS
	OPTIONS RequestMethod = "OPTIONS"
	// TRACE http TRACE
	TRACE RequestMethod = "TRACE"
)
