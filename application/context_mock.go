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
func (c *ContextMock) NewHTTPServiceRequest(serviceName string, method httputil.RequestMethod, path string, body []byte) (interface{}, error) {
	args := c.Called(serviceName, method, path, body)
	return args.Get(0), args.Error(1)
}

// Application application mock
func (c *ContextMock) Application() Application {
	args := c.Called()
	return args.Get(0).(Application)
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

// HTTPResponseUnauthorizedErr mock
func (c *ContextMock) HTTPResponseUnauthorizedErr(error) {}

// HTTPResponseInternalErr mock
func (c *ContextMock) HTTPResponseInternalErr(error) {}

// HTTPResponseErr mock
func (c *ContextMock) HTTPResponseErr(int, error) {}

// HTTPResponseOk mock
func (c *ContextMock) HTTPResponseOk(interface{}, error) {}

// HTTPResponseCreated mock
func (c *ContextMock) HTTPResponseCreated(interface{}, error) {}

// HTTPResponse mock
func (c *ContextMock) HTTPResponse(int, interface{}, error) {}
