package flow

import (
	"context"
	"encoding/json"

	"ant-flow/expression"
)

// Execer 表达式执行器
type Execer interface {
	// 执行表达式返回布尔类型的值
	ExecReturnBool(ctx context.Context, exp, params []byte) (bool, error)

	// 执行表达式返回字符串切片类型的值
	ExecReturnStringSlice(ctx context.Context, exp, params []byte) ([]string, error)
}

// NewQLangExecer 创建基于qlang的表达式执行器
func NewQLangExecer() Execer {
	return &execer{}
}

type execer struct {
}

func (*execer) ExecReturnBool(ctx context.Context, exp, params []byte) (bool, error) {
	var m map[string]interface{}
	err := json.Unmarshal(params, &m)
	if err != nil {
		return false, err
	}
	return expression.ExecParamBool(ctx, string(exp), m)
}

func (*execer) ExecReturnStringSlice(ctx context.Context, exp, params []byte) ([]string, error) {
	var m map[string]interface{}
	err := json.Unmarshal(params, &m)
	if err != nil {
		return nil, err
	}
	return expression.ExecParamSliceStr(ctx, string(exp), m)
}
