/*
 * @Author: Daniel TAN
 * @Description:
 * @Date: 2020-09-01 09:15:45
 * @LastEditTime: 2021-01-22 17:02:48
 * @LastEditors: Daniel TAN
 * @FilePath: /trinitygo/httputil/httputil.go
 */
package httputil

import "fmt"

var _ error = ResponseData{}

// ResponseData response data
type ResponseData struct {
	Status  int               `json:"status"`
	Data    interface{}       `json:"data,omitempty"`
	Err     interface{}       `json:"err,omitempty"`
	Runtime map[string]string `json:"runtime,omitempty"`
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
