package expression

import (
	"context"
	"time"

	"qlang.io/cl/qlang"
	//"qlang.io/cl/qlang"
)

func CreateExpContext(ctx context.Context) ExpContext {
	if ctx == nil {
		panic("ctx不能为nil")
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
