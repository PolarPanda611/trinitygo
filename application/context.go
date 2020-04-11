package application

import (
	"net/http"

	"github.com/PolarPanda611/trinitygo/httputils"
	"github.com/PolarPanda611/trinitygo/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Context record all thing inside one request
type Context interface {
	GetApplication() Application
	setRuntime(map[string]string)
	GetRuntime() map[string]string
	GetDB() *gorm.DB
	GetTXDB() *gorm.DB
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

// GetApplication get app
func (c *ContextImpl) GetApplication() Application {
	return c.app
}

// GetRuntime get runtime info
func (c *ContextImpl) GetRuntime() map[string]string {
	return c.runtime
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

// GetDB get db instance
func (c *ContextImpl) GetDB() *gorm.DB {
	return c.db
}

// GetTXDB get db instance
func (c *ContextImpl) GetTXDB() *gorm.DB {
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

// GetRuntime get runtime info
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
	err = utils.HTTPErrEncoder(err)
	if err != nil {
		if c.app.Conf().GetAtomicRequest() {
			c.SafeRollback()
		}
		_, resErr := utils.HTTPErrDecoder(err)
		c.c.AbortWithStatusJSON(status, httputils.ResponseData{
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
	c.c.JSON(status, httputils.ResponseData{
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
