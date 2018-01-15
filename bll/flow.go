package bll

import (
	"time"

	"gitee.com/antlinker/flow/schema"
	"gitee.com/antlinker/flow/util"
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
	if flow.Flag == 0 {
		flow.Flag = 1
	}
	return a.Models.Flow.CreateFlowBasic(flow, nodes, routers, assignments)
}

// GetNode 获取流程节点
func (a *Flow) GetNode(recordID string) (*schema.FlowNodes, error) {
	return a.Models.Flow.GetNode(recordID)
}

// GetFlowInstance 获取流程实例
func (a *Flow) GetFlowInstance(recordID string) (*schema.FlowInstances, error) {
	return a.Models.Flow.GetFlowInstance(recordID)
}

// GetNodeInstance 获取流程节点实例
func (a *Flow) GetNodeInstance(recordID string) (*schema.NodeInstances, error) {
	return a.Models.Flow.GetNodeInstance(recordID)
}

// QueryNodeRouters 查询节点路由
func (a *Flow) QueryNodeRouters(sourceNodeID string) ([]*schema.NodeRouters, error) {
	return a.Models.Flow.QueryNodeRouters(sourceNodeID)
}

// QueryNodeAssignments 查询节点指派
func (a *Flow) QueryNodeAssignments(nodeID string) ([]*schema.NodeAssignments, error) {
	return a.Models.Flow.QueryNodeAssignments(nodeID)
}

// CreateNodeInstance 创建节点实例
func (a *Flow) CreateNodeInstance(flowInstanceID, nodeID string, inputData []byte, candidates []string) (string, error) {
	nodeInstance := &schema.NodeInstances{
		RecordID:       util.UUID(),
		FlowInstanceID: flowInstanceID,
		NodeID:         nodeID,
		InputData:      string(inputData),
		Status:         1,
		Created:        time.Now().Unix(),
	}

	var nodeCandidates []*schema.NodeCandidates
	for _, c := range candidates {
		nodeCandidates = append(nodeCandidates, &schema.NodeCandidates{
			RecordID:       util.UUID(),
			NodeInstanceID: nodeInstance.RecordID,
			CandidateID:    c,
			Created:        nodeInstance.Created,
		})
	}

	err := a.Models.Flow.CreateNodeInstance(nodeInstance, nodeCandidates)
	if err != nil {
		return "", err
	}

	return nodeInstance.RecordID, nil
}

// DoneNodeInstance 完成节点实例
func (a *Flow) DoneNodeInstance(nodeInstanceID, processor string, outData []byte) error {
	nodeInstance, err := a.Models.Flow.GetNodeInstance(nodeInstanceID)
	if err != nil {
		return err
	} else if nodeInstance == nil || nodeInstance.Status == 2 {
		return nil
	}

	info := map[string]interface{}{
		"processor":    processor,
		"process_time": time.Now().Unix(),
		"out_data":     string(outData),
		"status":       2,
		"updated":      time.Now().Unix(),
	}
	return a.Models.Flow.UpdateNodeInstance(nodeInstanceID, info)
}

// CheckFlowInstanceTodo 检查流程实例待办事项
func (a *Flow) CheckFlowInstanceTodo(flowInstanceID string) (bool, error) {
	return a.Models.Flow.CheckFlowInstanceTodo(flowInstanceID)
}

// DoneFlowInstance 完成流程实例
func (a *Flow) DoneFlowInstance(flowInstanceID string) error {
	info := map[string]interface{}{
		"status": 9,
	}
	return a.Models.Flow.UpdateFlowInstance(flowInstanceID, info)
}

// LaunchFlowInstance 发起流程实例
func (a *Flow) LaunchFlowInstance(flowCode, nodeCode, launcher string, inputData []byte) (*schema.NodeInstances, error) {
	flow, err := a.Models.Flow.GetFlowByCode(flowCode)
	if err != nil {
		return nil, err
	} else if flow == nil {
		return nil, nil
	}

	node, err := a.Models.Flow.GetNodeByCode(flow.RecordID, nodeCode)
	if err != nil {
		return nil, err
	} else if node == nil {
		return nil, nil
	}

	flowInstance := &schema.FlowInstances{
		RecordID:   util.UUID(),
		FlowID:     flow.RecordID,
		Launcher:   launcher,
		LaunchTime: time.Now().Unix(),
		Status:     1,
		Created:    time.Now().Unix(),
	}

	nodeInstance := &schema.NodeInstances{
		RecordID:       util.UUID(),
		FlowInstanceID: flowInstance.RecordID,
		NodeID:         node.RecordID,
		InputData:      string(inputData),
		Status:         1,
		Created:        flowInstance.Created,
	}

	err = a.Models.Flow.CreateFlowInstance(flowInstance, nodeInstance)
	if err != nil {
		return nil, err
	}

	return nodeInstance, nil
}

// QueryNodeCandidates 查询节点候选人
func (a *Flow) QueryNodeCandidates(nodeInstanceID string) ([]*schema.NodeCandidates, error) {
	return a.Models.Flow.QueryNodeCandidates(nodeInstanceID)
}

// QueryTodoNodeInstances 查询用户的待办节点实例数据
func (a *Flow) QueryTodoNodeInstances(flowCode, userID string) ([]*schema.NodeInstances, error) {
	flow, err := a.Models.Flow.GetFlowByCode(flowCode)
	if err != nil {
		return nil, err
	} else if flow == nil {
		return nil, nil
	}
	return a.Models.Flow.QueryTodoNodeInstances(flow.RecordID, userID)
}
