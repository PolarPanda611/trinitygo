package httputils

// ResponseData response data
type ResponseData struct {
	Status  int         // the http response status  to return
	Result  interface{} // the response data  if req success
	Runtime map[string]string
}

type RequestMap struct {
	Method   RequestMethod
	SubPath  string
	FuncName string
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

// NewRequestMapping request mapping
func NewRequestMapping(method RequestMethod, path string, funcName string) *RequestMap {
	return &RequestMap{
		Method:   method,
		SubPath:  path,
		FuncName: funcName,
	}

}
