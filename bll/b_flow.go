package bll

import (
	"time"

	"github.com/antlinker/flow/model"
	"github.com/antlinker/flow/schema"
	"github.com/antlinker/flow/util"
)

// Flow 流程管理
type Flow struct {
	FlowModel *model.Flow `inject:""`
}

// GetFlow 获取流程数据
func (a *Flow) GetFlow(recordID string) (*schema.Flow, error) {
	return a.FlowModel.GetFlow(recordID)
}

// GetFlowByCode 根据编号查询流程数据
func (a *Flow) GetFlowByCode(code string) (*schema.Flow, error) {
	return a.FlowModel.GetFlowByCode(code)
}

// CreateFlow 创建流程数据
func (a *Flow) CreateFlow(flow *schema.Flow, nodes *schema.NodeOperating, forms *schema.FormOperating) error {
	if flow.Flag == 0 {
		flow.Flag = 1
	}
	return a.FlowModel.CreateFlow(flow, nodes, forms)
}

// GetNode 获取流程节点
func (a *Flow) GetNode(recordID string) (*schema.Node, error) {
	return a.FlowModel.GetNode(recordID)
}

// GetFlowInstance 获取流程实例
func (a *Flow) GetFlowInstance(recordID string) (*schema.FlowInstance, error) {
	return a.FlowModel.GetFlowInstance(recordID)
}

// GetFlowInstanceByNode 根据节点实例获取流程实例
func (a *Flow) GetFlowInstanceByNode(nodeInstanceID string) (*schema.FlowInstance, error) {
	return a.FlowModel.GetFlowInstanceByNode(nodeInstanceID)
}

// GetNodeInstance 获取流程节点实例
func (a *Flow) GetNodeInstance(recordID string) (*schema.NodeInstance, error) {
	return a.FlowModel.GetNodeInstance(recordID)
}

// QueryNodeRouters 查询节点路由
func (a *Flow) QueryNodeRouters(sourceNodeID string) ([]*schema.NodeRouter, error) {
	return a.FlowModel.QueryNodeRouters(sourceNodeID)
}

// QueryNodeAssignments 查询节点指派
func (a *Flow) QueryNodeAssignments(nodeID string) ([]*schema.NodeAssignment, error) {
	return a.FlowModel.QueryNodeAssignments(nodeID)
}

// CreateNodeInstance 创建节点实例
func (a *Flow) CreateNodeInstance(flowInstanceID, nodeID string, inputData []byte, candidates []string) (string, error) {
	nodeInstance := &schema.NodeInstance{
		RecordID:       util.UUID(),
		FlowInstanceID: flowInstanceID,
		NodeID:         nodeID,
		InputData:      string(inputData),
		Status:         1,
		Created:        time.Now().Unix(),
	}

	var nodeCandidates []*schema.NodeCandidate
	for _, c := range candidates {
		nodeCandidates = append(nodeCandidates, &schema.NodeCandidate{
			RecordID:       util.UUID(),
			NodeInstanceID: nodeInstance.RecordID,
			CandidateID:    c,
			Created:        nodeInstance.Created,
		})
	}

	err := a.FlowModel.CreateNodeInstance(nodeInstance, nodeCandidates)
	if err != nil {
		return "", err
	}

	return nodeInstance.RecordID, nil
}

// DoneNodeInstance 完成节点实例
func (a *Flow) DoneNodeInstance(nodeInstanceID, processor string, outData []byte) error {
	nodeInstance, err := a.FlowModel.GetNodeInstance(nodeInstanceID)
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
	return a.FlowModel.UpdateNodeInstance(nodeInstanceID, info)
}

// CheckFlowInstanceTodo 检查流程实例待办事项
func (a *Flow) CheckFlowInstanceTodo(flowInstanceID string) (bool, error) {
	return a.FlowModel.CheckFlowInstanceTodo(flowInstanceID)
}

// DoneFlowInstance 完成流程实例
func (a *Flow) DoneFlowInstance(flowInstanceID string) error {
	info := map[string]interface{}{
		"status": 9,
	}
	return a.FlowModel.UpdateFlowInstance(flowInstanceID, info)
}

// StopFlowInstance 停止流程实例
func (a *Flow) StopFlowInstance(flowInstanceID string) error {
	info := map[string]interface{}{
		"status": 9,
	}
	return a.FlowModel.UpdateFlowInstance(flowInstanceID, info)
}

// LaunchFlowInstance 发起流程实例
func (a *Flow) LaunchFlowInstance(flowCode, nodeCode, launcher string, inputData []byte) (*schema.NodeInstance, error) {
	flow, err := a.FlowModel.GetFlowByCode(flowCode)
	if err != nil {
		return nil, err
	} else if flow == nil {
		return nil, nil
	}

	node, err := a.FlowModel.GetNodeByCode(flow.RecordID, nodeCode)
	if err != nil {
		return nil, err
	} else if node == nil {
		return nil, nil
	}

	flowInstance := &schema.FlowInstance{
		RecordID:   util.UUID(),
		FlowID:     flow.RecordID,
		Launcher:   launcher,
		LaunchTime: time.Now().Unix(),
		Status:     1,
		Created:    time.Now().Unix(),
	}

	nodeInstance := &schema.NodeInstance{
		RecordID:       util.UUID(),
		FlowInstanceID: flowInstance.RecordID,
		NodeID:         node.RecordID,
		InputData:      string(inputData),
		Status:         1,
		Created:        flowInstance.Created,
	}

	err = a.FlowModel.CreateFlowInstance(flowInstance, nodeInstance)
	if err != nil {
		return nil, err
	}

	return nodeInstance, nil
}

// QueryNodeCandidates 查询节点候选人
func (a *Flow) QueryNodeCandidates(nodeInstanceID string) ([]*schema.NodeCandidate, error) {
	return a.FlowModel.QueryNodeCandidates(nodeInstanceID)
}

// QueryTodo 查询用户的待办节点实例数据
func (a *Flow) QueryTodo(flowCode, userID string) ([]*schema.FlowTodoResult, error) {
	return a.FlowModel.QueryTodo(flowCode, userID)
}

// QueryAllFlowPage 查询流程分页数据
func (a *Flow) QueryAllFlowPage(params schema.FlowQueryParam, pageIndex, pageSize uint) (int64, []*schema.FlowQueryResult, error) {
	return a.FlowModel.QueryAllFlowPage(params, pageIndex, pageSize)
}

// DeleteFlow 删除流程
func (a *Flow) DeleteFlow(flowID string) error {
	return a.FlowModel.DeleteFlow(flowID)
}

// QueryHistory 查询流程实例历史数据
func (a *Flow) QueryHistory(flowInstanceID string) ([]*schema.FlowHistoryResult, error) {
	return a.FlowModel.QueryHistory(flowInstanceID)
}

// QueryDoneIDs 查询已办理的流程实例ID列表
func (a *Flow) QueryDoneIDs(flowCode, userID string) ([]string, error) {
	return a.FlowModel.QueryDoneIDs(flowCode, userID)
}

// QueryGroupFlowPage 查询流程分组分页数据
func (a *Flow) QueryGroupFlowPage(params schema.FlowQueryParam, pageIndex, pageSize uint) (int64, []*schema.FlowQueryResult, error) {
	return a.FlowModel.QueryGroupFlowPage(params, pageIndex, pageSize)
}

// UpdateFlowInfo 更新流程
func (a *Flow) UpdateFlowInfo(recordID string, info map[string]interface{}) error {
	return a.FlowModel.Update(recordID, info)
}

// UpdateFlowStatus 更新流程状态
func (a *Flow) UpdateFlowStatus(recordID string, status int) error {
	info := map[string]interface{}{
		"updated": time.Now().Unix(),
		"status":  status,
	}

	return a.UpdateFlowInfo(recordID, info)
}
