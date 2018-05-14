package expression

import (
	"context"
	"database/sql"
	"time"

	"qlang.io/cl/qlang"
	//"qlang.io/cl/qlang"
)

type dbkey struct{}

// CreateExpContextByDB 创建含有DB的ctx
func CreateExpContextByDB(ctx context.Context, db *sql.DB) ExpContext {
	if ctx == nil {
		panic("ctx不能为nil")
	}
	ectx, ok := ctx.(*expContext)
	if ok {

		ectx.ctx = context.WithValue(ectx.ctx, dbkey{}, db)
		return ectx
	}

	return &expContext{
		ctx:        context.WithValue(ctx, dbkey{}, db),
		ql:         qlang.New(),
		predefined: predefined{data: make([]pairs, 0, 4)},
	}
}

// FromExpContextForDB 从ctx中获取*sql.DB
func FromExpContextForDB(ctx context.Context) *sql.DB {
	db := ctx.Value(dbkey{})
	if db == nil {
		return nil
	}
	dbs, ok := db.(*sql.DB)
	if ok {
		return dbs
	}
	return nil
}

// CreateExpContext 创建一个ExpContext
// 实现了context.Context接口
func CreateExpContext(ctx context.Context) ExpContext {
	if ctx == nil {
		panic("ctx不能为nil")
	}
	ectx, ok := ctx.(*expContext)
	if ok {
		return ectx
	}
	return &expContext{
		ctx:        ctx,
		ql:         qlang.New(),
		predefined: predefined{data: make([]pairs, 0, 4)},
	}
}
func qlangFromContext(ctx ExpContext) *qlang.Qlang {
	ql, ok := ctx.(*expContext)
	if ok {
		return ql.ql
	}
	return nil
}

type pairs struct {
	Key   string
	Value string
}

type expContext struct {
	predefined
	ctx context.Context
	ql  *qlang.Qlang
	err error
}

func (c *expContext) Var(key string) interface{} {
	return c.ql.Var(key)
}

func (c *expContext) AddVar(key string, value interface{}) {
	c.ql.SetVar(key, value)
}

// Deadline context.Context 接口实现
func (c *expContext) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

// Done context.Context 接口实现
func (c *expContext) Done() <-chan struct{} {
	return c.ctx.Done()
}

// Err context.Context 接口实现
func (c *expContext) Err() error {
	if c.err != nil {
		return c.err
	}
	return c.Err()
}

// Value context.Context 接口实现
func (c *expContext) Value(s interface{}) interface{} {
	return c.ctx.Value(s)
}
