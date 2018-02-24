package flow

import (
	"encoding/json"
	"net/http"

	"gitee.com/antlinker/flow/expression/sql"
	"gitee.com/antlinker/flow/schema"
	"gitee.com/antlinker/flow/service/db"
)

var (
	engine *Engine
)

// Init 初始化流程配置
func Init(opts ...db.Option) {
	db, err := db.NewMySQL(opts...)
	if err != nil {
		panic(err)
	}
	e, err := new(Engine).Init(db, NewXMLParser(), NewQLangExecer())
	if err != nil {
		panic(err)
	}
	engine = e
	sql.Reg(db.Db)
}

// SetParser 设定解析器
func SetParser(parser Parser) {
	engine.SetParser(parser)
}

// SetExecer 设定表达式执行器
func SetExecer(execer Execer) {
	engine.SetExecer(execer)
}

// LoadFile 加载流程文件数据
func LoadFile(name string) error {
	return engine.LoadFile(name)
}

// StartFlow 启动流程
// flowCode 流程编号
// nodeCode 开始节点编号
// userID 发起人
// input 输入数据
func StartFlow(flowCode, nodeCode, userID string, input interface{}) (*HandleResult, error) {
	inputData, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	return engine.StartFlow(flowCode, nodeCode, userID, inputData)
}

// HandleFlow 处理流程节点
// nodeInstanceID 节点实例内码
// userID 处理人
// input 输入数据
func HandleFlow(nodeInstanceID, userID string, input interface{}) (*HandleResult, error) {
	inputData, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	return engine.HandleFlow(nodeInstanceID, userID, inputData)
}

// StopFlow 停止流程
func StopFlow(nodeInstanceID string, allowStop func(*schema.FlowInstance) bool) error {
	return engine.StopFlow(nodeInstanceID, allowStop)
}

// QueryTodoFlows 查询流程待办数据
// flowCode 流程编号
// userID 待办人
func QueryTodoFlows(flowCode, userID string) ([]*schema.FlowTodoResult, error) {
	return engine.QueryTodoFlows(flowCode, userID)
}

// QueryFlowHistory 查询流程历史数据
// flowInstanceID 流程实例内码
func QueryFlowHistory(flowInstanceID string) ([]*schema.FlowHistoryResult, error) {
	return engine.QueryFlowHistory(flowInstanceID)
}

// QueryDoneFlowIDs 查询已办理的流程实例ID列表
func QueryDoneFlowIDs(flowCode, userID string) ([]string, error) {
	return engine.QueryDoneFlowIDs(flowCode, userID)
}

// StartServer 启动管理服务
func StartServer(opts ...ServerOption) http.Handler {
	srv := new(Server).Init(engine, opts...)
	return srv
}
