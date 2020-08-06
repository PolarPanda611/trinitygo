package application

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/PolarPanda611/trinitygo/httputil"
	"github.com/PolarPanda611/trinitygo/sd"
	"github.com/PolarPanda611/trinitygo/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var _ Context = new(ContextImpl)

// Context record all thing inside one request
type Context interface {
	NewHTTPServiceRequest(serviceName string, method httputil.RequestMethod, path string, body []byte) (int, interface{}, error)
	Application() Application
	setRuntime(map[string]string)
	Runtime() map[string]string
	DB() *gorm.DB
	DBTx() *gorm.DB
	SafeCommit()
	SafeRollback()
	DBTxIsOpen() bool
	setGinCTX(c *gin.Context)
	GinCtx() *gin.Context
	setDB(*gorm.DB)
	cleanRuntime()
	GetCurrentUser() (interface{}, error)
	HTTPStatus(code int)
	httpResponseUnauthorizedErr(error)
	HTTPResponseInternalErr(error)
	httpResponseErr(error)
	HTTPResponseOk(interface{}, error)
	httpResponseCreated(interface{}, error)
	httpResponseDeleted(interface{}, error)
	HTTPResponse(interface{}, error)
}

// ContextImpl Context impl
type ContextImpl struct {
	app      Application
	runtime  map[string]string
	db       *gorm.DB
	dbTxOpen bool
	// http
	c               *gin.Context
	funcToGetWhoAmI func(app Application, c *gin.Context, db *gorm.DB) (interface{}, error)
	status          int
}

// Application get app
func (c *ContextImpl) Application() Application {
	return c.app
}

// Runtime get runtime info
func (c *ContextImpl) Runtime() map[string]string {
	newRuntime := make(map[string]string, len(c.runtime))
	for k, v := range c.runtime {
		newRuntime[k] = v
	}
	return newRuntime
}

// setRuntime set runtime info
func (c *ContextImpl) setRuntime(runtime map[string]string) {
	c.runtime = runtime
}

// setRuntime set runtime info
func (c *ContextImpl) setGinCTX(ginCtx *gin.Context) {
	c.c = ginCtx
}

// GinCtx get gin ctx
func (c *ContextImpl) GinCtx() *gin.Context {
	return c.c
}

// DB get db instance
func (c *ContextImpl) DB() *gorm.DB {
	return c.db
}

// DBTx get db instance
func (c *ContextImpl) DBTx() *gorm.DB {
	if c.dbTxOpen {
		return c.db
	}
	c.db = c.db.Begin()
	c.dbTxOpen = true
	return c.db
}

// SafeCommit safe commit
func (c *ContextImpl) SafeCommit() {
	if c.dbTxOpen {
		c.db.Commit()
		c.dbTxOpen = false
	}
}

// SafeRollback safe rollback
func (c *ContextImpl) SafeRollback() {
	if c.dbTxOpen {
		c.db.Rollback()
		c.dbTxOpen = false
	}
}

// DBTxIsOpen return is transaction is open
func (c *ContextImpl) DBTxIsOpen() bool {
	return c.dbTxOpen
}

// setDB set db
func (c *ContextImpl) setDB(db *gorm.DB) {
	c.db = db
}

// CleanRuntime  clean runtime info
func (c *ContextImpl) cleanRuntime() {
	c.c = nil
	c.runtime = nil
	c.db = nil
	c.status = 0
	c.dbTxOpen = false
}

// HTTPStatus set http status
func (c *ContextImpl) HTTPStatus(code int) {
	c.status = code
}

// httpResponseUnauthorizedErr handle http 400 response
func (c *ContextImpl) httpResponseUnauthorizedErr(err error) {
	c.HTTPStatus(http.StatusUnauthorized)
	c.httpResponseErr(err)
}

// HTTPResponseInternalErr handle http 400 response
func (c *ContextImpl) HTTPResponseInternalErr(err error) {
	c.HTTPStatus(http.StatusBadRequest)
	c.httpResponseErr(err)
}

// httpResponseErr handle http 400 response
func (c *ContextImpl) httpResponseErr(err error) {
	if c.c == nil {
		panic("gin context not set , please if the func is used in http request handle")
	}
	if c.status == 0 {
		// default status 400
		c.status = 400
	}
	err = util.HTTPErrEncoder(err)
	if err != nil {
		if c.app.Conf().GetAtomicRequest() {
			c.SafeRollback()
		}
		_, resErr := util.HTTPErrDecoder(err)
		c.c.Error(fmt.Errorf("%v", err))
		c.c.Error(fmt.Errorf("%v", string(debug.Stack())))

		if c.app.ResponseFactory() != nil {
			c.c.JSON(c.status, c.app.ResponseFactory()(c.status, resErr, c.runtime))
		} else {
			if resErr != nil {
				c.c.AbortWithStatusJSON(c.status, resErr)
			} else {
				c.c.AbortWithStatus(c.status)
			}
		}
		return
	}
}

// HTTPResponseOk handle http http.StatusOK 200 response
func (c *ContextImpl) HTTPResponseOk(res interface{}, err error) {
	c.HTTPStatus(http.StatusOK)
	c.HTTPResponse(res, err)
}

// httpResponseCreated handle http http.StatusCreated 201 response
func (c *ContextImpl) httpResponseCreated(res interface{}, err error) {
	c.HTTPStatus(http.StatusCreated)
	c.HTTPResponse(res, err)
}

// httpResponseDeleted handle http http.httpResponseDeleted 204 response
func (c *ContextImpl) httpResponseDeleted(res interface{}, err error) {
	c.HTTPStatus(http.StatusNoContent)
	c.HTTPResponse(res, err)
}

// HTTPResponse handle response
func (c *ContextImpl) HTTPResponse(res interface{}, err error) {
	if c.c == nil {
		panic("gin context not set , please if the func is used in http request handle")
	}
	if err != nil {
		c.HTTPResponseInternalErr(err)
		return
	}
	if c.status == 0 {
		c.status = http.StatusOK
	}
	if c.app.Conf().GetAtomicRequest() {
		c.SafeCommit()
	}
	if c.app.ResponseFactory() != nil {
		c.c.JSON(c.status, c.app.ResponseFactory()(c.status, res, c.runtime))
	} else {
		if res != nil {
			c.c.JSON(c.status, res)
		} else {
			c.c.Status(c.status)
		}

	}

	return
}

// GetCurrentUser get current User
func (c *ContextImpl) GetCurrentUser() (interface{}, error) {
	if c.c == nil {
		panic("gin context not set , please check ")
	}
	if c.app == nil {
		panic("app not set , please check ")
	}
	if c.db == nil {
		panic("db not set , please check ")
	}
	return c.funcToGetWhoAmI(c.app, c.c, c.db)
}

//NewHTTPServiceRequest new http service request
func (c *ContextImpl) NewHTTPServiceRequest(serviceName string, method httputil.RequestMethod, path string, body []byte) (int, interface{}, error) {
	client, err := sd.NewEtcdHTTPClient(c.app.Conf().GetServiceDiscoveryAddress(), c.app.Conf().GetServiceDiscoveryPort(), serviceName, c.app.Conf().GetAppIdleTimeout())
	if err != nil {
		return 0, nil, err
	}
	header := map[string]string{
		"Authorization": c.c.GetHeader("Authorization"),
	}
	for k, v := range c.Runtime() {
		if kValue := c.c.GetString(k); kValue != v {
			header[k] = kValue
		} else {
			header[k] = v
		}
	}
	c.app.Logger().Info(fmt.Sprintf("http call %v , %v , ", method, path))
	return client.Request(method, path, body, header)
}

// NewContext new contedt
func NewContext(app Application, funcToGetWhoAmI func(app Application, c *gin.Context, db *gorm.DB) (interface{}, error)) Context {
	return &ContextImpl{
		app:             app,
		dbTxOpen:        false,
		funcToGetWhoAmI: funcToGetWhoAmI,
	}
}

// MockContext used for unittest , don't use in production
func MockContext(app Application, db *gorm.DB, c *gin.Context, runtime map[string]string) Context {
	return &ContextImpl{
		app:     app,
		db:      db,
		c:       c,
		runtime: runtime,
	}
}
