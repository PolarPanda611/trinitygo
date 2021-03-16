/*
 * @Author: Daniel TAN
 * @Description:
 * @Date: 2021-03-03 15:46:55
 * @LastEditTime: 2021-03-17 01:47:28
 * @LastEditors: Daniel TAN
 * @FilePath: /trinitygo/application/context_pool.go
 */
package application

import (
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

// ContextPool is the context pool, it's used inside router and the framework by itself.
//
// It's the only one real implementation inside this package because it used widely.
type ContextPool struct {
	pool    *sync.Pool
	newFunc func() Context // we need a field otherwise is not working if we change the return value
}

// New creates and returns a new context pool.
func New(newFunc func() Context) *ContextPool {
	c := &ContextPool{pool: &sync.Pool{}, newFunc: newFunc}
	c.pool.New = func() interface{} { return c.newFunc() }
	return c
}

// Attach changes the pool's return value Context.
func (c *ContextPool) Attach(newFunc func() Context) {
	c.newFunc = newFunc
}

// Acquire returns a Context from pool.
// See Release.
func (c *ContextPool) Acquire(app Application, runtime map[string]string, ginCtx *gin.Context) Context {
	ctx := c.pool.Get().(Context)
	ctx.setGinCTX(ginCtx)
	ctx.setRuntime(runtime)
	newDB := app.DB().New()
	newDB.SetLogger(NewDBLogger(app, runtime))
	userIDInt64, _ := strconv.ParseInt(ctx.GinCtx().GetString("user_id"), 10, 64)
	newDB = newDB.Set("user_id", userIDInt64)
	ctx.SetDB(newDB)
	return ctx
}

// Release puts a Context back to its pull, this function releases its resources.
// See Acquire.
func (c *ContextPool) Release(ctx Context) {
	ctx.cleanRuntime()
	c.pool.Put(ctx)
}

// ReleaseLight will just release the object back to the pool, but the
// clean method is caller's responsibility now
func (c *ContextPool) ReleaseLight(ctx Context) {
	c.pool.Put(ctx)
}
