package httputils

// ResponseData response data
type ResponseData struct {
	Status  int         // the http response status  to return
	Result  interface{} // the response data  if req success
	Runtime map[string]string
}
