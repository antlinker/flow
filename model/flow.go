package model

import (
	"database/sql"
	"fmt"

	"gitee.com/antlinker/flow/schema"
	"gitee.com/antlinker/flow/service/db"
	"github.com/pkg/errors"
)

// Flow 流程管理
type Flow struct {
	db *db.DB
}

// Init 初始化
func (a *Flow) Init(db *db.DB) *Flow {
	db.AddTableWithName(schema.Flows{}, schema.FlowsTableName)
	db.AddTableWithName(schema.FlowNodes{}, schema.FlowNodesTableName)
	db.AddTableWithName(schema.NodeRouters{}, schema.NodeRoutersTableName)
	db.AddTableWithName(schema.NodeAssignments{}, schema.NodeAssignmentsTableName)
	db.AddTableWithName(schema.FlowInstances{}, schema.FlowInstancesTableName)
	db.AddTableWithName(schema.NodeInstances{}, schema.NodeInstancesTableName)
	db.AddTableWithName(schema.NodeCandidates{}, schema.NodeCandidatesTableName)
	a.db = db
	return a
}

// CheckFlowCode 检查流程编号是否存在
func (a *Flow) CheckFlowCode(code string) (bool, error) {
	query := fmt.Sprintf("SELECT count(*) FROM %s WHERE deleted=0 AND code=?", schema.FlowsTableName)

	exists, err := a.db.CheckExists(query, code)
	if err != nil {
		return false, errors.Wrapf(err, "检查流程编号是否存在发生错误")
	}
	return exists, nil
}

// CreateFlowBasic 创建流程基础数据
func (a *Flow) CreateFlowBasic(flow *schema.Flows, nodes []*schema.FlowNodes, routers []*schema.NodeRouters, assignments []*schema.NodeAssignments) error {
	tran, err := a.db.Begin()
	if err != nil {
		return errors.Wrapf(err, "创建流程基础数据开启事物发生错误")
	}

	err = tran.Insert(flow)
	if err != nil {
		err = tran.Rollback()
		if err != nil {
			return errors.Wrapf(err, "创建流程基础数据回滚事物发生错误")
		}
		return errors.Wrapf(err, "插入流程数据发生错误")
	}

	for _, node := range nodes {
		err = tran.Insert(node)
		if err != nil {
			err = tran.Rollback()
			if err != nil {
				return errors.Wrapf(err, "创建流程基础数据回滚事物发生错误")
			}
			return errors.Wrapf(err, "插入流程节点数据发生错误")
		}
	}

	for _, router := range routers {
		err = tran.Insert(router)
		if err != nil {
			err = tran.Rollback()
			if err != nil {
				return errors.Wrapf(err, "创建流程基础数据回滚事物发生错误")
			}
			return errors.Wrapf(err, "插入流程节点路由数据发生错误")
		}
	}

	for _, assign := range assignments {
		err = tran.Insert(assign)
		if err != nil {
			err = tran.Rollback()
			if err != nil {
				return errors.Wrapf(err, "创建流程基础数据回滚事物发生错误")
			}
			return errors.Wrapf(err, "插入流程节点指派数据发生错误")
		}
	}

	err = tran.Commit()
	if err != nil {
		return errors.Wrapf(err, "创建流程基础数据提交事物发生错误")
	}
	return nil
}

// GetFlowByCode 根据编号查询流程数据
func (a *Flow) GetFlowByCode(code string) (*schema.Flows, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted=0 AND code=?", schema.FlowsTableName)

	var flow schema.Flows
	err := a.db.SelectOne(&flow, query, code)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "根据编号查询流程数据发生错误")
	}

	return &flow, nil
}
