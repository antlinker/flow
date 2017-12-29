package flow

import (
	"gitee.com/antlinker/flow/bll"
)

// HookHandle 定义钩子处理函数
type HookHandle func([]byte) error

// Engine 流程引擎
type Engine struct {
	blls *bll.All
}
