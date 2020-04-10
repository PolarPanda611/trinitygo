package application

import (
	"github.com/PolarPanda611/trinitygo/httputils"
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
	setDB(*gorm.DB)
	cleanRuntime()
	Response(int, interface{}, error)
	HandleResponse(c *gin.Context)
}

// ContextImpl Context impl
type ContextImpl struct {
	app      Application
	runtime  map[string]string
	db       *gorm.DB
	dbTxOpen bool
	// http
	status int
	res    interface{}
	err    error
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
	c.runtime = nil
	c.db = nil
}

// Response handle response
func (c *ContextImpl) Response(status int, res interface{}, err error) {
	c.status = status
	c.res = res
	c.err = err
}

func (c *ContextImpl) HandleResponse(context *gin.Context) {
	if c.err != nil {
		if c.app.Conf().GetAtomicRequest() {
			c.SafeRollback()
		}
		context.AbortWithStatusJSON(c.status, httputils.ResponseData{
			Status:  c.status,
			Result:  c.err,
			Runtime: c.runtime,
		})
		return
	}
	if c.app.Conf().GetAtomicRequest() {
		c.SafeCommit()
	}
	context.JSON(c.status, httputils.ResponseData{
		Status:  c.status,
		Result:  c.res,
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
