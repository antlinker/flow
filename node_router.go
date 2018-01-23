package flow

import (
	"context"
	"encoding/json"

	"gitee.com/antlinker/flow/schema"
	"github.com/pkg/errors"
)

// 定义错误
var (
	ErrNotFound = errors.New("未找到流程相关的信息")
)

type (
	// NextNodeHandle 定义下一节点处理函数
	NextNodeHandle func(*schema.Node, *schema.NodeInstance, []*schema.NodeCandidate)
	// EndHandle 定义流程结束处理函数
	EndHandle func(*schema.FlowInstance)
)

var defaultNodeRouterOptions = &nodeRouterOptions{
	autoStart: true,
}

type nodeRouterOptions struct {
	autoStart  bool
	onNextNode NextNodeHandle
	onFlowEnd  EndHandle
}

// NodeRouterOption 节点路由配置
type NodeRouterOption func(*nodeRouterOptions)

// AutoStartOption 自动开始流程配置
func AutoStartOption(autoStart bool) NodeRouterOption {
	return func(o *nodeRouterOptions) {
		o.autoStart = autoStart
	}
}

// OnNextNodeOption 注册下一节点处理事件配置
func OnNextNodeOption(fn NextNodeHandle) NodeRouterOption {
	return func(o *nodeRouterOptions) {
		o.onNextNode = fn
	}
}

// OnFlowEndOption 注册流程结束事件
func OnFlowEndOption(fn EndHandle) NodeRouterOption {
	return func(o *nodeRouterOptions) {
		o.onFlowEnd = fn
	}
}

// NodeRouter 节点路由
type NodeRouter struct {
	node         *schema.Node
	flowInstance *schema.FlowInstance
	nodeInstance *schema.NodeInstance
	inputData    []byte
	engine       *Engine
	opts         *nodeRouterOptions
	parent       *NodeRouter
	stop         bool
}

// Init 初始化节点路由
func (n *NodeRouter) Init(engine *Engine, nodeInstanceID string, inputData []byte, options ...NodeRouterOption) (*NodeRouter, error) {
	opts := defaultNodeRouterOptions
	for _, opt := range options {
		opt(opts)
	}

	n.opts = opts
	n.inputData = inputData
	n.engine = engine

	nodeInstance, err := n.engine.flowBll.GetNodeInstance(nodeInstanceID)
	if err != nil {
		return nil, err
	} else if nodeInstance == nil {
		return nil, ErrNotFound
	}
	n.nodeInstance = nodeInstance

	flowInstance, err := n.engine.flowBll.GetFlowInstance(nodeInstance.FlowInstanceID)
	if err != nil {
		return nil, err
	} else if flowInstance == nil {
		return nil, ErrNotFound
	}
	n.flowInstance = flowInstance

	node, err := n.engine.flowBll.GetNode(nodeInstance.NodeID)
	if err != nil {
		return nil, err
	} else if node == nil {
		return nil, ErrNotFound
	}
	n.node = node

	return n, nil
}

func (n *NodeRouter) next(nodeInstanceID, processor string) (*NodeRouter, error) {
	nextRouter, err := new(NodeRouter).Init(n.engine, nodeInstanceID, n.inputData)
	if err != nil {
		return nil, err
	}
	nextRouter.opts = n.opts
	nextRouter.parent = n

	err = nextRouter.Next(processor)
	if err != nil {
		return nil, err
	}
	return nextRouter, nil
}

// Next 流向下一节点
func (n *NodeRouter) Next(processor string) error {
	nodeType, err := GetNodeTypeByName(n.node.TypeCode)
	if err != nil {
		return err
	}

	if nodeType == UserTask && n.parent != nil {
		pNodeType, err := GetNodeTypeByName(n.parent.node.TypeCode)
		if err != nil {
			return err
		}

		if !(pNodeType == StartEvent && n.parent.opts.autoStart) {
			// 通知下一节点实例事件
			if fn := n.opts.onNextNode; fn != nil {
				candidates, err := n.engine.flowBll.QueryNodeCandidates(n.nodeInstance.RecordID)
				if err != nil {
					return err
				}
				fn(n.node, n.nodeInstance, candidates)
			}
			return nil
		}

	}

	// 完成当前节点
	err = n.engine.flowBll.DoneNodeInstance(n.nodeInstance.RecordID, processor, n.inputData)
	if err != nil {
		return err
	}

	// 如果当前节点是人工任务，检查下一节点是否是并行网关，如果是则检查还未完成的待办事项，如果有则停止流转
	if nodeType == UserTask && n.parent == nil {
		ok, err := n.checkNextNodeType(ParallelGateway)
		if err != nil {
			return err
		} else if ok {
			exists, err := n.engine.flowBll.CheckFlowInstanceTodo(n.flowInstance.RecordID)
			if err != nil {
				return err
			} else if exists {
				return nil
			}
		}
	}

	// 如果是结束时间或终止事件，则停止流转
	if nodeType == EndEvent ||
		nodeType == TerminateEvent {
		isEnd := false

		// 如果是结束事件，则检查还未完成的待办事项，如果没有则结束流程并通知结束事件
		if nodeType == EndEvent {
			exists, err := n.engine.flowBll.CheckFlowInstanceTodo(n.flowInstance.RecordID)
			if err != nil {
				return err
			} else if !exists {
				isEnd = true
			}
		}

		// 如果是终止事件，则结束流程并通知结束事件
		if nodeType == TerminateEvent {
			isEnd = true
		}

		if isEnd {
			// 流程实例结束处理
			err = n.engine.flowBll.DoneFlowInstance(n.flowInstance.RecordID)
			if err != nil {
				return err
			}

			n.stop = true
			if fn := n.opts.onFlowEnd; fn != nil {
				fn(n.flowInstance)
			}
		}
		return nil
	}

	// 增加下一节点
	nodeInstanceIDs, err := n.addNextNodeInstances()
	if err != nil {
		return err
	}

	for _, instanceID := range nodeInstanceIDs {
		nextRouter, err := n.next(instanceID, processor)
		if err != nil {
			return err
		}

		if nextRouter.stop {
			break
		}
	}

	return nil
}

// 增加下一处理节点实例
func (n *NodeRouter) addNextNodeInstances() ([]string, error) {
	routers, err := n.engine.flowBll.QueryNodeRouters(n.node.RecordID)
	if err != nil {
		return nil, err
	} else if len(routers) == 0 {
		return nil, nil
	}

	var nodeInstanceIDs []string
	for _, r := range routers {
		if r.Expression != "" {
			allow, err := n.engine.execer.ExecReturnBool(context.Background(), []byte(r.Expression), n.getExpData())
			if err != nil {
				return nil, err
			} else if !allow {
				continue
			}
		}

		// 查询指派人表达式
		assigns, err := n.engine.flowBll.QueryNodeAssignments(r.TargetNodeID)
		if err != nil {
			return nil, err
		}

		var candidates []string
		for _, assign := range assigns {
			ss, err := n.engine.execer.ExecReturnStringSlice(context.Background(), []byte(assign.Expression), n.getExpData())
			if err != nil {
				return nil, err
			}
			candidates = append(candidates, ss...)
		}

		instanceID, err := n.engine.flowBll.CreateNodeInstance(n.flowInstance.RecordID, r.TargetNodeID, n.inputData, candidates)
		if err != nil {
			return nil, err
		}
		nodeInstanceIDs = append(nodeInstanceIDs, instanceID)
	}
	return nodeInstanceIDs, nil
}

// 检查下一节点类型
func (n *NodeRouter) checkNextNodeType(t NodeType) (bool, error) {
	routers, err := n.engine.flowBll.QueryNodeRouters(n.node.RecordID)
	if err != nil {
		return false, err
	} else if len(routers) == 0 {
		return false, nil
	}

	for _, r := range routers {
		if r.Expression != "" {
			allow, err := n.engine.execer.ExecReturnBool(context.Background(), []byte(r.Expression), n.getExpData())
			if err != nil {
				return false, err
			} else if !allow {
				continue
			}
		}

		node, err := n.engine.flowBll.GetNode(r.TargetNodeID)
		if err != nil {
			return false, err
		} else if node == nil {
			return false, nil
		}

		if node.TypeCode == t.String() {
			return true, nil
		}
	}

	return false, nil
}

// 获取表达式数据
func (n *NodeRouter) getExpData() []byte {
	var input map[string]interface{}
	json.Unmarshal(n.inputData, &input)

	r := map[string]interface{}{
		"input": input,
		"flow":  n.flowInstance,
		"node":  n.nodeInstance,
	}
	b, _ := json.Marshal(r)
	return b
}
