package expression

import (
	"context"

	_ "qlang.io/cl/qlang"
)

// ExpContext 表达式上下文　，实现　context.Context接口
type ExpContext interface {
	context.Context
	NameSpacer
}

// NameSpacer 命名域
type NameSpacer interface {
	Var(key string) interface{}
	SetVar(key string, value interface{})
	SetJsonStr(key string, value string)
}

// Execer 表达式执行器
type Execer interface {
	// 执行表达式
	Exec(ctx context.Context, exp, params []byte) (interface{}, error)

	// 执行表达式返回布尔类型的值
	ExecToBool(ctx context.Context, exp, params []byte) (bool, error)

	// 执行表达式返回字符串切片类型的值
	ExecToStringSlice(ctx context.Context, exp, params []byte) ([]string, error)
}

// Import 导入模块扩展
func Import(name string, table map[string]interface{}) {
	qlang.Import(name, table)
}
