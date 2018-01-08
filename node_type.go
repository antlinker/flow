package flow

import "errors"

// NodeType 节点类型
type NodeType string

func (n NodeType) String() string {
	return string(n)
}

const (
	// StartEvent 开始事件
	StartEvent NodeType = "startEvent"
	// EndEvent 结束事件
	EndEvent NodeType = "endEvent"
	// TerminateEvent 终止事件
	TerminateEvent NodeType = "terminateEvent"
	// UserTask 人工任务
	UserTask NodeType = "userTask"
	// ExclusiveGateway 排他网关
	ExclusiveGateway NodeType = "exclusiveGateway"
	// ParallelGateway 并行网关
	ParallelGateway NodeType = "parallelGateway"
	// Unknown 未知类型
	Unknown NodeType = "Unknown"
)

// GetNodeTypeByName 转换节点类型
func GetNodeTypeByName(s string) (NodeType, error) {
	switch s {
	case "startEvent":
		return StartEvent, nil
	case "endEvent":
		return EndEvent, nil
	case "terminateEvent":
		return TerminateEvent, nil
	case "userTask":
		return UserTask, nil
	case "exclusiveGateway":
		return ExclusiveGateway, nil
	case "parallelGateway":
		return ParallelGateway, nil
	}
	return Unknown, errors.New(s + "不支持的类型")
}
