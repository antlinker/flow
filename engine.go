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
)

// HookHandle 定义钩子处理函数
type HookHandle func([]byte) error

// Engine 流程引擎
type Engine struct {
	flowBll *bll.Flow
	parser  Parser
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
