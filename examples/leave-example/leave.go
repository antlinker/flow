package main

import (
	"encoding/json"
	"flag"

	"github.com/LyricTian/logger"

	"gitee.com/antlinker/flow"
	"gitee.com/antlinker/flow/service/db"
)

var (
	dsn   string
	trace bool
	bpmn  string
	end   chan struct{}
)

func init() {
	flag.StringVar(&dsn, "d", "root:123456@tcp(192.168.33.90:3306)/flows?charset=utf8", "mysql连接串")
	flag.BoolVar(&trace, "t", false, "trace")
	flag.StringVar(&bpmn, "b", "", "bpmn流程")
}

func main() {
	flag.Parse()
	end = make(chan struct{})

	flow.Init(&db.Config{
		DSN:          dsn,
		Trace:        true,
		MaxIdleConns: 100,
		MaxOpenConns: 100,
	})

	err := flow.LoadFile(bpmn)
	if err != nil {
		logger.Fatalf("加载文件发生错误:%s", err.Error())
	}

	input := map[string]interface{}{
		"day": 1,
		"bzr": "S002",
	}

	result, err := flow.StartFlow("process_leave", "node_start", "S001", input)
	if err != nil {
		logger.Fatalf("启动流程发生错误：%s", err.Error())
	}

	br, _ := json.Marshal(result)
	logger.Infof("流程启动成功：%s", string(br))

	todos, err := flow.QueryTodoFlows("process_leave", "S002")
	if err != nil {
		logger.Fatalf("查询流程待办发生错误:%s", err.Error())
	}

	bts, _ := json.Marshal(todos)
	logger.Infof("待办事项：%s", string(bts))

	for _, todo := range todos {
		input := map[string]interface{}{
			"day":    1,
			"bzr":    "S002",
			"action": "pass",
		}
		result, err := flow.HandleFlow(todo.RecordID, "S002", input)
		if err != nil {
			logger.Fatalf("处理流程待办发生错误:%s", err.Error())
		}

		br, _ := json.Marshal(result)
		logger.Infof("处理流程结果：%s", string(br))
	}

	<-end
}
