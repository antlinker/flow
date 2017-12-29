package bll

import (
	"gitee.com/antlinker/flow/schema"
)

// Flow 流程管理
type Flow struct {
	*Bll
}

// CheckFlowCode 检查流程编号是否存在
func (a *Flow) CheckFlowCode(code string) (bool, error) {
	return a.Models.Flow.CheckFlowCode(code)
}

// CreateFlowBasic 创建流程基础数据
func (a *Flow) CreateFlowBasic(flow *schema.Flows, nodes []*schema.FlowNodes, routers []*schema.NodeRouters, assignments []*schema.NodeAssignments) error {
	return a.Models.Flow.CreateFlowBasic(flow, nodes, routers, assignments)
}
