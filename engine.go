package flow

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"ant-flow/bll"
	"ant-flow/register"
	"ant-flow/schema"
	"ant-flow/service/db"
	"ant-flow/util"
	"github.com/facebookgo/inject"
	"github.com/pkg/errors"
)

// Engine 流程引擎
type Engine struct {
	flowBll *bll.Flow
	parser  Parser
	execer  Execer
}

// Init 初始化流程引擎
func (e *Engine) Init(db *db.DB, parser Parser, execer Execer) (*Engine, error) {

	var (
		g       inject.Graph
		flowBll bll.Flow
	)

	err := g.Provide(&inject.Object{Value: db},
		&inject.Object{Value: &flowBll})
	if err != nil {
		return e, err
	}

	err = g.Populate()
	if err != nil {
		return e, err
	}

	register.FlowDBMap(db)
	err = db.CreateTablesIfNotExists()
	if err != nil {
		return e, err
	}

	e.flowBll = &flowBll
	e.parser = parser
	e.execer = execer
	return e, nil
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

func (e *Engine) parseFormOperating(formOperating *schema.FormOperating, flow *schema.Flow, node *schema.Node, formResult *NodeFormResult) {
	if formResult.ID == "" && len(formResult.Fields) == 0 {
		return
	}

	formExists := false
	for _, f := range formOperating.FormGroup {
		if f.Code == formResult.ID {
			node.FormID = f.RecordID
			formExists = true
			break
		}
	}
	if formExists || len(formResult.Fields) == 0 {
		return
	}

	form := &schema.Form{
		RecordID: util.UUID(),
		FlowID:   flow.RecordID,
		Code:     formResult.ID,
		TypeCode: "META",
		Created:  flow.Created,
	}

	if form.Code == "" {
		form.Code = util.UUID()
	}

	meta, _ := json.Marshal(formResult.Fields)
	form.Data = string(meta)

	for _, ff := range formResult.Fields {
		field := &schema.FormField{
			RecordID:     util.UUID(),
			FormID:       form.RecordID,
			Code:         ff.ID,
			Label:        ff.Label,
			TypeCode:     ff.Type,
			DefaultValue: ff.DefaultValue,
			Created:      flow.Created,
		}

		for _, item := range ff.Values {
			formOperating.FieldOptionGroup = append(formOperating.FieldOptionGroup, &schema.FieldOption{
				RecordID:  util.UUID(),
				FieldID:   field.RecordID,
				ValueID:   item.ID,
				ValueName: item.Name,
				Created:   flow.Created,
			})
		}

		for _, item := range ff.Properties {
			formOperating.FieldPropertyGroup = append(formOperating.FieldPropertyGroup, &schema.FieldProperty{
				RecordID: util.UUID(),
				FieldID:  field.RecordID,
				Code:     item.ID,
				Value:    item.Value,
				Created:  flow.Created,
			})
		}

		for _, item := range ff.Validations {
			formOperating.FieldValidationGroup = append(formOperating.FieldValidationGroup, &schema.FieldValidation{
				RecordID:         util.UUID(),
				FieldID:          field.RecordID,
				ConstraintName:   item.Name,
				ConstraintConfig: item.Config,
				Created:          flow.Created,
			})
		}

		formOperating.FormFieldGroup = append(formOperating.FormFieldGroup, field)
	}

	formOperating.FormGroup = append(formOperating.FormGroup, form)
	node.FormID = form.RecordID
}

// 创建节点操作
func (e *Engine) parseOperating(flow *schema.Flow, nodeResults []*NodeResult) (*schema.NodeOperating, *schema.FormOperating) {
	nodeOperating := &schema.NodeOperating{
		NodeGroup: make([]*schema.Node, len(nodeResults)),
	}
	formOperating := &schema.FormOperating{}

	for i, n := range nodeResults {
		node := &schema.Node{
			RecordID: util.UUID(),
			FlowID:   flow.RecordID,
			Code:     n.NodeID,
			Name:     n.NodeName,
			TypeCode: n.NodeType.String(),
			OrderNum: strconv.FormatInt(int64(i+10), 10),
			Created:  flow.Created,
		}

		if n.FormResult != nil {
			e.parseFormOperating(formOperating, flow, node, n.FormResult)
		}

		for _, exp := range n.CandidateExpressions {
			nodeOperating.AssignmentGroup = append(nodeOperating.AssignmentGroup, &schema.NodeAssignment{
				RecordID:   util.UUID(),
				NodeID:     node.RecordID,
				Expression: exp,
				Created:    flow.Created,
			})
		}

		nodeOperating.NodeGroup[i] = node
	}

	var getNodeRecordID = func(nodeCode string) string {
		for _, n := range nodeOperating.NodeGroup {
			if n.Code == nodeCode {
				return n.RecordID
			}
		}
		return ""
	}

	for _, n := range nodeResults {
		for _, r := range n.Routers {
			nodeOperating.RouterGroup = append(nodeOperating.RouterGroup, &schema.NodeRouter{
				RecordID:     util.UUID(),
				SourceNodeID: getNodeRecordID(n.NodeID),
				TargetNodeID: getNodeRecordID(r.TargetNodeID),
				Expression:   r.Expression,
				Explain:      r.Explain,
				Created:      flow.Created,
			})
		}
	}

	return nodeOperating, formOperating
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

	nodeOperating, formOperating := e.parseOperating(flow, result.Nodes)
	return e.flowBll.CreateFlow(flow, nodeOperating, formOperating)
}

// HandleResult 处理结果
type HandleResult struct {
	IsEnd        bool                 `json:"is_end"`        // 是否结束
	NextNodes    []*NextNode          `json:"next_nodes"`    // 下一处理节点
	FlowInstance *schema.FlowInstance `json:"flow_instance"` // 流程实例
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
	result.FlowInstance = nr.GetFlowInstance()

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
	nodeInstance, err := e.flowBll.GetNodeInstance(nodeInstanceID)
	if err != nil {
		return nil, err
	} else if nodeInstance.Status != 1 {
		return nil, nil
	}
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

// StopFlowInstance 停止流程实例
func (e *Engine) StopFlowInstance(flowInstanceID string, allowStop func(*schema.FlowInstance) bool) error {
	flowInstance, err := e.flowBll.GetFlowInstance(flowInstanceID)
	if err != nil {
		return err
	}

	if allowStop != nil && !allowStop(flowInstance) {
		return errors.New("不允许停止流程")
	}

	return e.flowBll.StopFlowInstance(flowInstanceID)
}

// QueryTodoFlows 查询流程待办数据
// flowCode 流程编号
// userID 待办人
func (e *Engine) QueryTodoFlows(flowCode, userID string) ([]*schema.FlowTodoResult, error) {
	return e.flowBll.QueryTodo(flowCode, userID)
}

// QueryFlowHistory 查询流程历史数据
// flowInstanceID 流程实例内码
func (e *Engine) QueryFlowHistory(flowInstanceID string) ([]*schema.FlowHistoryResult, error) {
	return e.flowBll.QueryHistory(flowInstanceID)
}

// QueryDoneFlowIDs 查询已办理的流程实例ID列表
func (e *Engine) QueryDoneFlowIDs(flowCode, userID string) ([]string, error) {
	return e.flowBll.QueryDoneIDs(flowCode, userID)
}

// QueryNodeCandidates 查询节点实例的候选人ID列表
func (e *Engine) QueryNodeCandidates(nodeInstanceID string) ([]string, error) {
	candidates, err := e.flowBll.QueryNodeCandidates(nodeInstanceID)
	if err != nil {
		return nil, err
	}

	ids := make([]string, len(candidates))

	for i, c := range candidates {
		ids[i] = c.CandidateID
	}

	return ids, nil
}

// GetNodeInstance 获取节点实例
func (e *Engine) GetNodeInstance(nodeInstanceID string) (*schema.NodeInstance, error) {
	return e.flowBll.GetNodeInstance(nodeInstanceID)
}
