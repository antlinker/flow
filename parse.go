package flow

import (
	"context"
)

// Parser 流程数据解析器
type Parser interface {
	// 解析流程定义数据
	Parse(ctx context.Context, data []byte) (*ParseResult, error)
}

// ParseResult 流程数据
type ParseResult struct {
	FlowID      string        // 流程ID
	FlowName    string        // 流程名称
	FlowVersion int64         // 流程版本号
	FlowStatus  int           // 流程状态(1:可用 2:不可用)
	Nodes       []*NodeResult // 节点数据
}

// NodeResult 节点数据
type NodeResult struct {
	NodeID               string            // 节点ID
	NodeName             string            // 节点名称
	NodeType             NodeType          // 节点类型
	Routers              []*RouterResult   // 节点路由
	Properties           []*PropertyResult // 节点属性
	CandidateExpressions []string          // 候选人表达式
	FormResult           *NodeFormResult   // 节点表单
}

// RouterResult 节点路由数据
type RouterResult struct {
	TargetNodeID string // 目标节点ID
	Explain      string // 说明
	Expression   string // 条件表达式
}

// PropertyResult 节点属性
type PropertyResult struct {
	Name  string // 属性名称
	Value string // 属性值
}

// NodeFormResult 节点表单
type NodeFormResult struct {
	ID     string             // 表单ID
	Fields []*FormFieldResult // 表单字段
}

// FormFieldResult 表单字段
type FormFieldResult struct {
	ID           string             // 字段ID
	Type         string             // 字段类型
	Label        string             // 字段标签
	DefaultValue string             // 默认值
	Values       []*FieldOption     // 枚举类型
	Validations  []*FieldValidation // 字段验证
	Properties   []*FieldProperty   // 字段属性
}

// FieldValidation 字段验证
type FieldValidation struct {
	Name   string // 约束名
	Config string // 约束配置
}

// FieldProperty 字段属性
type FieldProperty struct {
	ID    string // 属性ID
	Value string // 属性值
}

// FieldOption 枚举选项
type FieldOption struct {
	ID   string // 选项值ID
	Name string // 选项值名称
}
