package flow

import (
	"context"

	"github.com/antlinker/flow/expression"
)

type (
	expKey  struct{}
	flagKey struct{}
)

// NewExpContext 创建表达式的上下文值
func NewExpContext(ctx context.Context, exp expression.ExpContext) context.Context {
	return context.WithValue(ctx, expKey{}, exp)
}

// FromExpContext 获取表达式的上下文
func FromExpContext(ctx context.Context) (expression.ExpContext, bool) {
	exp, ok := ctx.Value(expKey{}).(expression.ExpContext)
	return exp, ok
}

// NewFlagContext 创建flag的上下文值
func NewFlagContext(ctx context.Context, flag string) context.Context {
	return context.WithValue(ctx, flagKey{}, flag)
}

// FromFlagContext 获取flag的上下文
func FromFlagContext(ctx context.Context) (string, bool) {
	flag, ok := ctx.Value(flagKey{}).(string)
	return flag, ok
}
