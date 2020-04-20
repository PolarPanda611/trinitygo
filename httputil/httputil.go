package httputil

// ResponseData response data
type ResponseData struct {
	Status  int         // the http response status  to return
	Result  interface{} `json:"Result,omitempty"` // the response data  if req success
	Error   interface{} `json:"Error,omitempty"`  // the response data  if req success
	Runtime map[string]string
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
