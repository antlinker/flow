package flow

import (
	"context"
)

// Parser 流程数据解析器
type Parser interface {
	// 解析流程定义数据
	Parse(ctx context.Context, data []byte) (*ParseResult, error)
}

// ParseResult 流程数据解析结果
type ParseResult struct {
	FlowID      string        // 流程ID
	FlowName    string        // 流程名称
	FlowVersion int           // 流程版本号
	Nodes       []*NodeResult // 节点数据
}

// NodeResult 节点数据解析结果
type NodeResult struct {
	NodeID               string          // 节点ID
	NodeName             string          // 节点名称
	NodeType             NodeType        // 节点类型
	Routers              []*RouterResult // 节点路由
	CandidateExpressions []string        // 候选人表达式
}

// RouterResult 节点路由数据解析结果
type RouterResult struct {
	TargetNodeID string // 目标节点ID
	Explain      string // 说明
	Expression   string // 条件表达式
}
