package application

import (
	"fmt"

	"github.com/PolarPanda611/trinitygo/httputil"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
)

var _ Context = new(ContextMock)

// ContextMock mock impl for Context
type ContextMock struct {
	mock.Mock
}

// NewHTTPServiceRequest Mock
func (c *ContextMock) NewHTTPServiceRequest(serviceName string, method httputil.RequestMethod, path string, body []byte) (int, interface{}, error) {
	args := c.Called(serviceName, method, path, body)
	return args.Int(0), args.Get(1), args.Error(2)
}

// Application application mock
func (c *ContextMock) Application() Application {
	args := c.Called()
	return args.Get(0).(Application)
}

// AutoFreeOn set auto free on
func (c *ContextMock) AutoFreeOn() {
}

// AutoFreeOff set auto free off
func (c *ContextMock) AutoFreeOff() {
}

func (c *ContextMock) AutoFree() bool {
	args := c.Called()
	return args.Bool(0)
}
func (c *ContextMock) SetIsConfigured() {}
func (c *ContextMock) IsConfigured() bool {
	args := c.Called()
	return args.Bool(0)
}

func (c *ContextMock) setRuntime(map[string]string) {}

// Runtime runtime mock
func (c *ContextMock) Runtime() map[string]string {
	args := c.Called()
	return args.Get(0).(map[string]string)
}

// DB mock
func (c *ContextMock) DB() *gorm.DB {
	args := c.Called()
	return args.Get(0).(*gorm.DB)
}

// DBTx mock
func (c *ContextMock) DBTx() *gorm.DB {
	args := c.Called()
	return args.Get(0).(*gorm.DB)
}

// SafeCommit mock
func (c *ContextMock) SafeCommit() {}

// SafeRollback mock
func (c *ContextMock) SafeRollback() {}

// DBTxIsOpen mock
func (c *ContextMock) DBTxIsOpen() bool {
	args := c.Called()
	return args.Get(0).(bool)
}
func (c *ContextMock) setGinCTX(gCtx *gin.Context) {
}

// GinCtx mock
func (c *ContextMock) GinCtx() *gin.Context {
	args := c.Called()
	return args.Get(0).(*gin.Context)
}
func (c *ContextMock) setDB(*gorm.DB) {}

func (c *ContextMock) cleanRuntime() {
	fmt.Println("123")
}

// GetCurrentUser get current user
func (c *ContextMock) GetCurrentUser() (interface{}, error) {
	args := c.Called()
	return args.Get(0), args.Error(1)
}

// HTTPStatus mock
func (c *ContextMock) HTTPStatus(code int) {}

// httpResponseUnauthorizedErr mock
func (c *ContextMock) httpResponseUnauthorizedErr(error) {}

// HTTPResponseInternalErr mock
func (c *ContextMock) HTTPResponseInternalErr(error) {}

// httpResponseErr mock
func (c *ContextMock) httpResponseErr(error) {}

// HTTPResponseOk mock
func (c *ContextMock) HTTPResponseOk(interface{}, error) {}

// httpResponseCreated mock
func (c *ContextMock) httpResponseCreated(interface{}, error) {}

// HTTPResponse mock
func (c *ContextMock) HTTPResponse(interface{}, error) {}

// httpResponseDeleted mock
func (c *ContextMock) httpResponseDeleted(res interface{}, err error) {}
