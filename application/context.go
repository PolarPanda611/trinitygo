package application

import (
	"net/http"

	"github.com/PolarPanda611/trinitygo/httputil"
	"github.com/PolarPanda611/trinitygo/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Context record all thing inside one request
type Context interface {
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
	HTTPResponseUnauthorizedErr(error)
	HTTPResponseInternalErr(error)
	HTTPResponseErr(int, error)
	HTTPResponseOk(interface{}, error)
	HTTPResponseCreated(interface{}, error)
	HTTPResponse(int, interface{}, error)
}

// ContextImpl Context impl
type ContextImpl struct {
	app      Application
	runtime  map[string]string
	db       *gorm.DB
	dbTxOpen bool
	// http
	c *gin.Context
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
		c.dbTxOpen = true
	}
}

// SafeRollback safe rollback
func (c *ContextImpl) SafeRollback() {
	if c.dbTxOpen {
		c.db.Rollback()
		c.dbTxOpen = false
	}
}

// DBTxIsOpen return is transactioon is open
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
}

// HTTPResponseUnauthorizedErr handle http 400 response
func (c *ContextImpl) HTTPResponseUnauthorizedErr(err error) {
	c.HTTPResponseErr(401, err)
}

// HTTPResponseInternalErr handle http 400 response
func (c *ContextImpl) HTTPResponseInternalErr(err error) {
	c.HTTPResponseErr(400, err)
}

// HTTPResponseErr handle http 400 response
func (c *ContextImpl) HTTPResponseErr(status int, err error) {
	if c.c == nil {
		panic("gin context not set , please if the func is used in http request handle")
	}
	err = util.HTTPErrEncoder(err)
	if err != nil {
		if c.app.Conf().GetAtomicRequest() {
			c.SafeRollback()
		}
		_, resErr := util.HTTPErrDecoder(err)
		c.c.AbortWithStatusJSON(status, httputil.ResponseData{
			Status:  status,
			Error:   resErr,
			Runtime: c.runtime,
		})
		return
	}
}

// HTTPResponseOk handle http http.StatusOK 200 response
func (c *ContextImpl) HTTPResponseOk(res interface{}, err error) {
	c.HTTPResponse(http.StatusOK, res, err)
}

// HTTPResponseCreated handle http http.StatusCreated 201 response
func (c *ContextImpl) HTTPResponseCreated(res interface{}, err error) {
	c.HTTPResponse(http.StatusCreated, res, err)
}

// HTTPResponse handle respoonse
func (c *ContextImpl) HTTPResponse(status int, res interface{}, err error) {
	if c.c == nil {
		panic("gin context not set , please if the func is used in http request handle")
	}
	if err != nil {
		c.HTTPResponseInternalErr(err)
		return
	}
	if c.app.Conf().GetAtomicRequest() {
		c.SafeCommit()
	}
	c.c.JSON(status, httputil.ResponseData{
		Status:  status,
		Result:  res,
		Runtime: c.runtime,
	})
	return
}

// NewContext new contedt
func NewContext(app Application) Context {
	return &ContextImpl{
		app:      app,
		dbTxOpen: false,
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
