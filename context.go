package flow

import (
	"gitee.com/antlinker/flow/schema"
)

// Context 定义流程上下文
type Context struct {
	input        []byte                // 输入数据
	output       []byte                // 输出数据
	flowInstance *schema.FlowInstances // 流程实例
	engine       *Engine               // 流程引擎
}
