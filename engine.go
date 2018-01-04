package flow

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gitee.com/antlinker/flow/bll"
	"gitee.com/antlinker/flow/schema"
	"gitee.com/antlinker/flow/util"
	"github.com/pkg/errors"
)

// Engine 流程引擎
type Engine struct {
	flowBll *bll.Flow
	parser  Parser
	execer  Execer
}

func (e *Engine) parseFile(name string) ([]byte, error) {
	fullName, err := filepath.Abs(name)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(fullName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// LoadFile 加载文件数据
func (e *Engine) LoadFile(name string) error {
	data, err := e.parseFile(name)
	if err != nil {
		return err
	}

	result, err := e.parser.Parse(context.Background(), data)
	if err != nil {
		return err
	}

	// 检查流程编号是否存在，如果存在则不处理
	exists, err := e.flowBll.CheckFlowCode(result.FlowID)
	if err != nil {
		return err
	} else if exists {
		return nil
	}

	flow := &schema.Flows{
		RecordID: util.UUID(),
		Code:     result.FlowID,
		Name:     result.FlowName,
		Version:  result.FlowVersion,
		XML:      string(data),
		Created:  time.Now().Unix(),
	}

	var (
		nodes       = make([]*schema.FlowNodes, len(result.Nodes))
		nodeRouters []*schema.NodeRouters
		nodeAssigns []*schema.NodeAssignments
	)

	for i, n := range result.Nodes {
		node := &schema.FlowNodes{
			RecordID: util.UUID(),
			FlowID:   flow.RecordID,
			Code:     n.NodeID,
			Name:     n.NodeName,
			TypeCode: n.NodeType.String(),
			OrderNum: strconv.FormatInt(int64(i+10), 10),
			Created:  flow.Created,
		}

		for _, r := range n.Routers {
			nodeRouters = append(nodeRouters, &schema.NodeRouters{
				RecordID:     util.UUID(),
				SourceNodeID: node.RecordID,
				TargetNodeID: r.TargetNodeID,
				Expression:   r.Expression,
				Explain:      r.Explain,
				Created:      flow.Created,
			})
		}

		for _, exp := range n.CandidateExpressions {
			nodeAssigns = append(nodeAssigns, &schema.NodeAssignments{
				RecordID:   util.UUID(),
				NodeID:     node.RecordID,
				Expression: exp,
				Created:    flow.Created,
			})
		}

		nodes[i] = node
	}

	return e.flowBll.CreateFlowBasic(flow, nodes, nodeRouters, nodeAssigns)
}

// HandleResult 处理结果
type HandleResult struct {
	IsEnd     bool        // 是否结束
	NextNodes []*NextNode // 下一处理节点
}

// NextNode 下一节点
type NextNode struct {
	Node         *schema.FlowNodes // 节点信息
	CandidateIDs []string          // 节点候选人
}

func (e *Engine) nextFlowHandle(nodeInstanceID, userID string, inputData []byte) (*HandleResult, error) {
	var result HandleResult

	var onNextNode = OnNextNodeOption(func(node *schema.FlowNodes, nodeInstance *schema.NodeInstances, nodeCandidates []*schema.NodeCandidates) {
		var cids []string
		for _, nc := range nodeCandidates {
			cids = append(cids, nc.CandidateID)
		}

		result.NextNodes = append(result.NextNodes, &NextNode{
			Node:         node,
			CandidateIDs: cids,
		})
	})

	var onFlowEnd = OnFlowEndOption(func(_ *schema.FlowInstances) {
		result.IsEnd = true
	})

	nr, err := new(NodeRouter).Init(e, nodeInstanceID, inputData, onNextNode, onFlowEnd)
	if err != nil {
		return nil, err
	}

	err = nr.Next(userID)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// StartFlow 启动流程
// flowCode 流程编号
// nodeCode 开始节点编号
// userID 发起人
// inputData 输入数据
func (e *Engine) StartFlow(flowCode, nodeCode, userID string, inputData []byte) (*HandleResult, error) {
	nodeInstance, err := e.flowBll.LaunchFlowInstance(flowCode, nodeCode, userID, inputData)
	if err != nil {
		return nil, err
	} else if nodeInstance == nil {
		return nil, errors.New("未找到流程信息")
	}

	return e.nextFlowHandle(nodeInstance.RecordID, userID, inputData)
}

// HandleFlow 处理流程节点
// nodeInstanceID 节点实例内码
// userID 处理人
// inputData 输入数据
func (e *Engine) HandleFlow(nodeInstanceID, userID string, inputData []byte) (*HandleResult, error) {
	return e.nextFlowHandle(nodeInstanceID, userID, inputData)
}

// QueryTodoFlows 查询待办流程数据
// flowCode 流程编号
// userID 待办人
func (e *Engine) QueryTodoFlows(flowCode, userID string) ([]*schema.NodeInstances, error) {
	return nil, nil
}
