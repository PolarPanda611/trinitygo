package httputils

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
	GET     RequestMethod = "GET"
	HEAD    RequestMethod = "HEAD"
	POST    RequestMethod = "POST"
	PUT     RequestMethod = "PUT"
	PATCH   RequestMethod = "PATCH"
	DELETE  RequestMethod = "DELETE"
	OPTIONS RequestMethod = "OPTIONS"
	TRACE   RequestMethod = "TRACE"
)
