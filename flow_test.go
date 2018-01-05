package flow_test

import (
	"encoding/json"
	"testing"
	"time"

	"gitee.com/antlinker/flow"
	"gitee.com/antlinker/flow/service/db"
)

func init() {
	flow.Init(&db.Config{
		DSN:          "root:123456@tcp(192.168.33.90:3306)/flows?charset=utf8",
		Trace:        true,
		MaxIdleConns: 100,
		MaxOpenConns: 100,
		MaxLifetime:  time.Hour * 2,
	})

	err := flow.LoadFile("test_data/leave.bpmn")
	if err != nil {
		panic(err)
	}
}

func TestLeaveBzrApprovalPass(t *testing.T) {
	var (
		flowCode = "process_leave_test"
		bzr      = "T002"
	)

	input := map[string]interface{}{
		"day": 1,
		"bzr": bzr,
	}

	// 开始流程
	result, err := flow.StartFlow(flowCode, "node_start", "T001", input)
	if err != nil {
		t.Fatal(err.Error())
	}

	if result.NextNodes[0].CandidateIDs[0] != bzr {
		t.Fatalf("无效的下一级流转：%s", result.String())
	}

	// 查询待办
	todos, err := flow.QueryTodoFlows(flowCode, bzr)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(todos) != 1 {
		bts, _ := json.Marshal(todos)
		t.Fatalf("无效的待办数据:%s", string(bts))
	}

	// 处理流程（通过）
	input["action"] = "pass"
	result, err = flow.HandleFlow(todos[0].RecordID, bzr, input)
	if err != nil {
		t.Fatal(err.Error())
	}

	// 流程结束
	if !result.IsEnd {
		t.Fatalf("无效的处理结果：%s", result.String())
	}
}

func TestLeaveBzrApprovalBack(t *testing.T) {
	var (
		flowCode = "process_leave_test"
		launcher = "T001"
		bzr      = "T002"
	)

	input := map[string]interface{}{
		"day": 1,
		"bzr": bzr,
	}

	// 开始流程
	result, err := flow.StartFlow(flowCode, "node_start", launcher, input)
	if err != nil {
		t.Fatal(err.Error())
	}

	if result.NextNodes[0].CandidateIDs[0] != bzr {
		t.Fatalf("无效的下一级流转：%s", result.String())
	}

	// 查询待办
	todos, err := flow.QueryTodoFlows(flowCode, bzr)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// 处理流程（退回）
	input["action"] = "back"
	result, err = flow.HandleFlow(todos[0].RecordID, bzr, input)
	if err != nil {
		t.Fatal(err.Error())
	}

	if result.IsEnd ||
		result.NextNodes[0].CandidateIDs[0] != launcher {
		t.Fatalf("无效的处理结果：%s", result.String())
	}

	// 查询退回流程
	todos, err = flow.QueryTodoFlows(flowCode, launcher)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// 处理退回流程
	delete(input, "action")
	result, err = flow.HandleFlow(todos[0].RecordID, launcher, input)
	if err != nil {
		t.Fatal(err.Error())
	}

	if result.NextNodes[0].CandidateIDs[0] != bzr {
		t.Fatalf("无效的下一级流转：%s", result.String())
	}

	// 查询待办流程
	todos, err = flow.QueryTodoFlows(flowCode, bzr)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// 处理流程（通过）
	input["action"] = "pass"
	result, err = flow.HandleFlow(todos[0].RecordID, bzr, input)
	if err != nil {
		t.Fatal(err.Error())
	}

	// 流程结束
	if !result.IsEnd {
		t.Fatalf("无效的处理结果：%s", result.String())
	}
}
