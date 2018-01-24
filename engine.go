package flow

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gitee.com/antlinker/flow/bll"
	"gitee.com/antlinker/flow/schema"
	"gitee.com/antlinker/flow/service/db"
	"gitee.com/antlinker/flow/util"
	"github.com/pkg/errors"
)

// Engine 流程引擎
type Engine struct {
	flowBll *bll.Flow
	parser  Parser
	execer  Execer
}

// Init 初始化流程引擎
func (e *Engine) Init(db *db.DB, parser Parser, execer Execer) *Engine {
	blls := new(bll.All).Init(db)
	e.flowBll = blls.Flow
	e.parser = parser
	e.execer = execer
	return e
}

// SetParser 设定解析器
func (e *Engine) SetParser(parser Parser) {
	e.parser = parser
}

// SetExecer 设定表达式执行器
func (e *Engine) SetExecer(execer Execer) {
	e.execer = execer
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
	return e.CreateFlow(data)
}

// CreateFlow 创建流程数据
func (e *Engine) CreateFlow(data []byte) error {
	result, err := e.parser.Parse(context.Background(), data)
	if err != nil {
		return err
	}

	// 检查流程是否存在，如果存在则检查版本号是否一致，如果不一致则创建新流程
	oldFlow, err := e.flowBll.GetFlowByCode(result.FlowID)
	if err != nil {
		return err
	} else if oldFlow != nil {
		if result.FlowVersion <= oldFlow.Version {
			return nil
		}
	}

	flow := &schema.Flow{
		RecordID: util.UUID(),
		Code:     result.FlowID,
		Name:     result.FlowName,
		Version:  result.FlowVersion,
		XML:      string(data),
		Created:  time.Now().Unix(),
	}

	var (
		nodes       = make([]*schema.Node, len(result.Nodes))
		nodeRouters []*schema.NodeRouter
		nodeAssigns []*schema.NodeAssignment
	)

	for i, n := range result.Nodes {
		node := &schema.Node{
			RecordID: util.UUID(),
			FlowID:   flow.RecordID,
			Code:     n.NodeID,
			Name:     n.NodeName,
			TypeCode: n.NodeType.String(),
			OrderNum: strconv.FormatInt(int64(i+10), 10),
			Created:  flow.Created,
		}

		for _, exp := range n.CandidateExpressions {
			nodeAssigns = append(nodeAssigns, &schema.NodeAssignment{
				RecordID:   util.UUID(),
				NodeID:     node.RecordID,
				Expression: exp,
				Created:    flow.Created,
			})
		}

		nodes[i] = node
	}

	var getNodeRecordID = func(nodeCode string) string {
		for _, n := range nodes {
			if n.Code == nodeCode {
				return n.RecordID
			}
		}
		return ""
	}

	for _, n := range result.Nodes {
		for _, r := range n.Routers {
			nodeRouters = append(nodeRouters, &schema.NodeRouter{
				RecordID:     util.UUID(),
				SourceNodeID: getNodeRecordID(n.NodeID),
				TargetNodeID: getNodeRecordID(r.TargetNodeID),
				Expression:   r.Expression,
				Explain:      r.Explain,
				Created:      flow.Created,
			})
		}
	}

	return e.flowBll.CreateFlowBasic(flow, nodes, nodeRouters, nodeAssigns)
}

// HandleResult 处理结果
type HandleResult struct {
	IsEnd     bool        // 是否结束
	NextNodes []*NextNode // 下一处理节点
}

func (r *HandleResult) String() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// NextNode 下一节点
type NextNode struct {
	Node         *schema.Node // 节点信息
	CandidateIDs []string     // 节点候选人
}

func (e *Engine) nextFlowHandle(nodeInstanceID, userID string, inputData []byte) (*HandleResult, error) {
	var result HandleResult

	var onNextNode = OnNextNodeOption(func(node *schema.Node, nodeInstance *schema.NodeInstance, nodeCandidates []*schema.NodeCandidate) {
		var cids []string
		for _, nc := range nodeCandidates {
			cids = append(cids, nc.CandidateID)
		}

		result.NextNodes = append(result.NextNodes, &NextNode{
			Node:         node,
			CandidateIDs: cids,
		})
	})

	var onFlowEnd = OnFlowEndOption(func(_ *schema.FlowInstance) {
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

// StopFlow 停止流程
func (e *Engine) StopFlow(nodeInstanceID string, allowStop func(*schema.FlowInstance) bool) error {
	flowInstance, err := e.flowBll.GetFlowInstanceByNode(nodeInstanceID)
	if err != nil {
		return err
	} else if flowInstance == nil {
		return errors.New("流程不存在")
	}

	if allowStop != nil && !allowStop(flowInstance) {
		return errors.New("不允许停止流程")
	}

	return e.flowBll.StopFlowInstance(flowInstance.RecordID)
}

// QueryTodoFlows 查询流程待办数据
// flowCode 流程编号
// userID 待办人
func (e *Engine) QueryTodoFlows(flowCode, userID string) ([]*schema.NodeInstance, error) {
	return e.flowBll.QueryTodoNodeInstances(flowCode, userID)
}
