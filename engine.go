package flow

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/antlinker/flow/bll"
	"github.com/antlinker/flow/register"
	"github.com/antlinker/flow/schema"
	"github.com/antlinker/flow/service/db"
	"github.com/antlinker/flow/util"
	"github.com/facebookgo/inject"
	"github.com/pkg/errors"
)

// Logger 定义日志接口
type Logger interface {
	Errorf(format string, args ...interface{})
}

// AutoCallbackHandler 自动执行节点回调处理
type AutoCallbackHandler func(action, flag, userID string, input []byte, result *HandleResult) error

// Engine 流程引擎
type Engine struct {
	flowBll      *bll.Flow
	parser       Parser
	execer       Execer
	logger       Logger
	timingStart  bool
	timingTicker *time.Ticker
	timingWg     *sync.WaitGroup
	getDBContext func(flag string) context.Context
	autoCallback AutoCallbackHandler
}

// Init 初始化流程引擎
func (e *Engine) Init(parser Parser, execer Execer, sqlDB *sql.DB, trace bool) (*Engine, error) {

	var (
		g       inject.Graph
		flowBll bll.Flow
	)

	db := db.NewMySQLWithDB(sqlDB, trace)
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

// SetLogger 设定日志接口
func (e *Engine) SetLogger(logger Logger) {
	e.logger = logger
}

// SetGetDBContext 设定获取DB上下文
func (e *Engine) SetGetDBContext(fn func(flag string) context.Context) {
	e.getDBContext = fn
}

// SetAutoCallback 设定自动节点回调函数
func (e *Engine) SetAutoCallback(callback AutoCallbackHandler) {
	e.autoCallback = callback
}

// FlowBll 流程业务
func (e *Engine) FlowBll() *bll.Flow {
	return e.flowBll
}

func (e *Engine) errorf(format string, args ...interface{}) {
	if e.logger != nil {
		e.logger.Errorf(format, args...)
	}
}

// StartTiming 启动定时器
func (e *Engine) StartTiming(interval time.Duration) {
	if e.timingStart {
		return
	}

	e.timingStart = true
	e.timingWg = new(sync.WaitGroup)
	e.timingTicker = time.NewTicker(interval)

	go func() {
		for range e.timingTicker.C {
			err := e.handleTiming()
			if err != nil {
				e.errorf("处理定时任务发生错误：%v", err)
			}
		}
	}()
}

func (e *Engine) handleTiming() error {
	defer func() {
		if err := recover(); err != nil {
			e.errorf("处理定时任务发生崩溃：%v", err)
		}
	}()
	items, err := e.flowBll.QueryExpiredNodeTiming()
	if err != nil {
		return err
	}

	for _, item := range items {
		err = e.handleExpiredNodeTiming(item)
		if err != nil {
			return err
		}
	}
	return nil
}

// 处理定时节点
func (e *Engine) handleExpiredNodeTiming(item *schema.NodeTiming) error {
	e.timingWg.Add(1)
	defer e.timingWg.Done()

	ni, err := e.flowBll.GetNodeInstance(item.NodeInstanceID)
	if err != nil {
		return err
	} else if ni == nil || ni.Status != 1 {
		return nil
	}

	ctx := context.Background()
	if fn := e.getDBContext; fn != nil {
		ctx = fn(item.Flag)
	}

	if item.Input != "" {
		var v map[string]interface{}
		_ = json.Unmarshal([]byte(item.Input), &v)

		if ni.InputData != "" {
			iv := make(map[string]interface{})
			_ = json.Unmarshal([]byte(ni.InputData), &iv)

			for key, val := range v {
				iv[key] = val
			}
			v = iv
		}

		buf, _ := json.Marshal(v)
		ni.InputData = string(buf)
	}

	result, err := e.HandleFlow(ctx, item.NodeInstanceID, item.Processor, []byte(ni.InputData))
	if err != nil {
		return err
	}

	err = e.flowBll.DeleteNodeTiming(item.NodeInstanceID)
	if err != nil {
		return err
	}

	if fn := e.autoCallback; fn != nil {
		err = fn("", item.Flag, item.Processor, []byte(ni.InputData), result)
		if err != nil {
			return err
		}
	}

	return nil
}

// StopTiming 停止定时器
func (e *Engine) StopTiming() {
	if !e.timingStart {
		return
	}
	e.timingStart = false
	e.timingTicker.Stop()
	e.timingWg.Wait()
}

// 读取XML文件
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
	_, err = e.CreateFlow(data)
	return err
}

func (e *Engine) parseFormOperating(formOperating *schema.FormOperating, flow *schema.Flow, node *schema.Node, formResult *NodeFormResult) {
	if formResult.ID == "" {
		return
	}

	for _, f := range formOperating.FormGroup {
		if f.Code == formResult.ID {
			node.FormID = f.RecordID
			return
		}
	}

	form := &schema.Form{
		RecordID: util.UUID(),
		FlowID:   flow.RecordID,
		Code:     formResult.ID,
		TypeCode: "META",
		Created:  flow.Created,
	}

	// 解析URL类型
	if fields := formResult.Fields; len(fields) == 2 {
		if fields[0].ID == "type_code" &&
			fields[0].DefaultValue == "URL" &&
			fields[1].ID == "data" {
			form.TypeCode = "URL"
			form.Data = fields[1].DefaultValue
			formOperating.FormGroup = append(formOperating.FormGroup, form)
			node.FormID = form.RecordID
			return
		}
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

		// 增加节点属性
		for _, p := range n.Properties {
			nodeOperating.PropertyGroup = append(nodeOperating.PropertyGroup, &schema.NodeProperty{
				RecordID: util.UUID(),
				NodeID:   getNodeRecordID(n.NodeID),
				Name:     p.Name,
				Value:    p.Value,
				Created:  flow.Created,
			})
		}
	}

	return nodeOperating, formOperating
}

// CreateFlow 创建流程数据
func (e *Engine) CreateFlow(data []byte) (string, error) {
	result, err := e.parser.Parse(context.Background(), data)
	if err != nil {
		return "", err
	}

	// 检查流程是否存在，如果存在则检查版本号是否一致，如果不一致则创建新流程
	oldFlow, err := e.flowBll.GetFlowByCode(result.FlowID)
	if err != nil {
		return "", err
	} else if oldFlow != nil {
		if result.FlowVersion <= oldFlow.Version {
			return oldFlow.RecordID, nil
		}
	}

	flow := &schema.Flow{
		RecordID: util.UUID(),
		Code:     result.FlowID,
		Name:     result.FlowName,
		Version:  result.FlowVersion,
		XML:      string(data),
		Status:   result.FlowStatus,
		Created:  time.Now().Unix(),
	}

	nodeOperating, formOperating := e.parseOperating(flow, result.Nodes)

	// 解析节点表单数据
	for _, node := range result.Nodes {
		// 查找表单ID不为空并且不包含表单字段的节点
		if node.FormResult != nil && node.FormResult.ID != "" && len(node.FormResult.Fields) == 0 {
			// 查找表单ID
			var formID string
			for _, form := range formOperating.FormGroup {
				if form.Code == node.FormResult.ID {
					formID = form.RecordID
					break
				}
			}
			if formID != "" {
				for i, ns := range nodeOperating.NodeGroup {
					if ns.Code == node.NodeID {
						nodeOperating.NodeGroup[i].FormID = formID
						break
					}
				}
			}
		}
	}

	err = e.flowBll.CreateFlow(flow, nodeOperating, formOperating)
	if err != nil {
		return "", err
	}
	return flow.RecordID, nil
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
	Node         *schema.Node         // 节点信息
	CandidateIDs []string             // 节点候选人
	NodeInstance *schema.NodeInstance // 节点实例
}

func (e *Engine) nextFlowHandle(ctx context.Context, nodeInstanceID, userID string, inputData []byte) (*HandleResult, error) {
	var result HandleResult

	var onNextNode = OnNextNodeOption(func(node *schema.Node, nodeInstance *schema.NodeInstance, nodeCandidates []*schema.NodeCandidate) {
		var cids []string
		for _, nc := range nodeCandidates {
			cids = append(cids, nc.CandidateID)
		}

		result.NextNodes = append(result.NextNodes, &NextNode{
			Node:         node,
			NodeInstance: nodeInstance,
			CandidateIDs: cids,
		})
	})

	var onFlowEnd = OnFlowEndOption(func(_ *schema.FlowInstance) {
		result.IsEnd = true
	})

	nr, err := new(NodeRouter).Init(ctx, e, nodeInstanceID, inputData, onNextNode, onFlowEnd)
	if err != nil {
		return nil, err
	}

	err = nr.Next(userID)
	if err != nil {
		return nil, err
	}
	result.FlowInstance = nr.GetFlowInstance()

	if !result.IsEnd {
		for _, item := range result.NextNodes {
			prop, verr := e.flowBll.GetNodeProperty(item.Node.RecordID)
			if verr != nil {
				return nil, verr
			}

			// 检查节点是否设定定时器，如果设定则加入定时
			if v := prop["timing"]; v != "" {
				expired, verr := strconv.Atoi(v)
				if verr == nil && expired > 0 {
					nt := &schema.NodeTiming{
						NodeInstanceID: item.NodeInstance.RecordID,
						Processor:      item.CandidateIDs[0],
						Input:          prop["timing_input"],
						ExpiredAt:      time.Now().Add(time.Duration(expired) * time.Minute).Unix(),
						Created:        time.Now().Unix(),
					}

					if v, ok := FromFlagContext(ctx); ok {
						nt.Flag = v
					}

					err = e.flowBll.CreateNodeTiming(nt)
					if err != nil {
						e.errorf("%+v", err)
					}
				}
			}
		}
	}

	return &result, nil
}

// StartFlow 启动流程
// flowCode 流程编号
// nodeCode 开始节点编号
// userID 发起人
// inputData 输入数据
func (e *Engine) StartFlow(ctx context.Context, flowCode, nodeCode, userID string, inputData []byte) (*HandleResult, error) {
	nodeInstance, err := e.flowBll.LaunchFlowInstance(flowCode, nodeCode, userID, inputData)
	if err != nil {
		return nil, err
	} else if nodeInstance == nil {
		return nil, errors.New("未找到流程信息")
	}

	return e.nextFlowHandle(ctx, nodeInstance.RecordID, userID, inputData)
}

// LaunchFlow 发起流程（基于流程ID）
func (e *Engine) LaunchFlow(ctx context.Context, flowID, userID string, inputData []byte) (*HandleResult, error) {
	_, ni, err := e.flowBll.LaunchFlowInstance2(flowID, userID, 1, inputData)
	if err != nil {
		return nil, err
	}
	return e.nextFlowHandle(ctx, ni.RecordID, userID, inputData)
}

// HandleFlow 处理流程节点
// nodeInstanceID 节点实例内码
// userID 处理人
// inputData 输入数据
func (e *Engine) HandleFlow(ctx context.Context, nodeInstanceID, userID string, inputData []byte) (*HandleResult, error) {
	// 检查是否是节点候选人
	exists, err := e.flowBll.CheckNodeCandidate(nodeInstanceID, userID)
	if err != nil {
		return nil, err
	} else if !exists {
		return nil, fmt.Errorf("无效的节点处理人")
	}

	nodeInstance, err := e.flowBll.GetNodeInstance(nodeInstanceID)
	if err != nil {
		return nil, err
	} else if nodeInstance == nil || nodeInstance.Status != 1 {
		return nil, fmt.Errorf("无效的处理节点")
	}

	return e.nextFlowHandle(ctx, nodeInstanceID, userID, inputData)
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
	return e.flowBll.QueryTodo("", flowCode, userID, 100)
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
