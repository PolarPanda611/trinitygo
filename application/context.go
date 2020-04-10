package application

import (
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
}

// ContextImpl Context impl
type ContextImpl struct {
	app      Application
	runtime  map[string]string
	db       *gorm.DB
	dbTxOpen bool
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
	}
}

// SafeRollback safe rollback
func (c *ContextImpl) SafeRollback() {
	if c.dbTxOpen {
		c.db.Rollback()
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

// NewContext new contedt
func NewContext(app Application) Context {
	return &ContextImpl{
		app:      app,
		dbTxOpen: false,
	}
}
