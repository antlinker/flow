package flow

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
)
