package expression

import (
	"context"

	"qlang.io/cl/qlang"
)

// ExpContext 表达式上下文　，实现　context.Context接口
type ExpContext interface {
	context.Context
	Predefined
	// 获取变量key的当前值
	Var(key string) interface{}
	// 为表达式中的变量赋值
	AddVar(key string, value interface{})
}

type Predefined interface {
	// 预定义变量
	// 将添加的数据转为脚本中的语句添加到脚本中
	// key为变量名
	// value 为脚本值的原型
	// 		值类型为string 则在脚本中 key=value key类型由值字符串内容决定。
	// 可以使用该方法初始化变量，变量类型可以是各种类型，只要符合脚本规则
	PredefinedVar(key string, value string)

	// 预定义变量
	// key为变量名
	// value 将转为json字符串 在调用 PredefinedVar(key string, value string)
	// 在脚本中是一个map 可以通过key[键]访问 也可以通过key.键访问，和访问json格式方式一样
	PredefinedJson(key string, value interface{}) error
}

// Execer 表达式执行器
type Execer interface {
	Predefined
	// 执行表达式
	// ctx 上下文
	// exp 为执行表达式
	Exec(ctx ExpContext, exp string) (*OutData, error)
	ImportAlias(model, alias string)
	Import(model string)
}

// Import 导入模块扩展
// 同一模块只能被导入一次，多次导入会导致panic
func Import(name string, table map[string]interface{}) {
	qlang.Import(name, table)
}
