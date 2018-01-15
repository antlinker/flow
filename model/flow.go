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
	db.AddTableWithName(schema.FlowForms{}, schema.FlowFormsTableName)
	db.AddTableWithName(schema.FormFields{}, schema.FormFieldsTableName)
	db.AddTableWithName(schema.FieldProperties{}, schema.FieldPropertiesTableName)
	db.AddTableWithName(schema.FieldValidation{}, schema.FieldValidationTableName)

	a.db = db
	return a
}

// CheckFlowCode 检查流程编号是否存在
func (a *Flow) CheckFlowCode(code string) (bool, error) {
	query := fmt.Sprintf("SELECT count(*) FROM %s WHERE deleted=0 AND flag=1 AND code=?", schema.FlowsTableName)

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
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted=0 AND flag=1 AND code=?", schema.FlowsTableName)

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

// GetNode 获取流程节点
func (a *Flow) GetNode(recordID string) (*schema.FlowNodes, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted=0 AND record_id=?", schema.FlowNodesTableName)

	var item schema.FlowNodes
	err := a.db.SelectOne(&item, query, recordID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "获取流程节点发生错误")
	}

	return &item, nil
}

// GetNodeByCode 根据节点编号获取流程节点
func (a *Flow) GetNodeByCode(flowID, nodeCode string) (*schema.FlowNodes, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted=0 AND flow_id=? AND code=? ORDER BY order_num", schema.FlowNodesTableName)

	var items []*schema.FlowNodes
	_, err := a.db.Select(&items, query, flowID, nodeCode)
	if err != nil {
		return nil, errors.Wrapf(err, "根据节点编号获取流程节点发生错误")
	} else if len(items) == 0 {
		return nil, nil
	}

	return items[0], nil
}

// GetFlowInstance 获取流程实例
func (a *Flow) GetFlowInstance(recordID string) (*schema.FlowInstances, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted=0 AND record_id=?", schema.FlowInstancesTableName)

	var item schema.FlowInstances
	err := a.db.SelectOne(&item, query, recordID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "获取流程实例发生错误")
	}

	return &item, nil
}

// GetNodeInstance 获取流程节点实例
func (a *Flow) GetNodeInstance(recordID string) (*schema.NodeInstances, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted=0 AND record_id=?", schema.NodeInstancesTableName)

	var item schema.NodeInstances
	err := a.db.SelectOne(&item, query, recordID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "获取流程节点实例发生错误")
	}

	return &item, nil
}

// QueryNodeRouters 查询节点路由
func (a *Flow) QueryNodeRouters(sourceNodeID string) ([]*schema.NodeRouters, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted=0 AND source_node_id=?", schema.NodeRoutersTableName)

	var items []*schema.NodeRouters
	_, err := a.db.Select(&items, query, sourceNodeID)
	if err != nil {
		return nil, errors.Wrapf(err, "查询节点路由发生错误")
	}

	return items, nil
}

// QueryNodeAssignments 查询节点指派
func (a *Flow) QueryNodeAssignments(nodeID string) ([]*schema.NodeAssignments, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted=0 AND node_id=?", schema.NodeAssignmentsTableName)

	var items []*schema.NodeAssignments
	_, err := a.db.Select(&items, query, nodeID)
	if err != nil {
		return nil, errors.Wrapf(err, "查询节点指派发生错误")
	}

	return items, nil
}

// CreateNodeInstance 创建流程节点实例
func (a *Flow) CreateNodeInstance(nodeInstance *schema.NodeInstances, nodeCandidates []*schema.NodeCandidates) error {
	tran, err := a.db.Begin()
	if err != nil {
		return errors.Wrapf(err, "创建流程节点实例开启事物发生错误")
	}

	err = tran.Insert(nodeInstance)
	if err != nil {
		err = tran.Rollback()
		if err != nil {
			return errors.Wrapf(err, "创建流程节点实例回滚事物发生错误")
		}
		return errors.Wrapf(err, "插入流程节点实例数据发生错误")
	}

	for _, c := range nodeCandidates {
		err = tran.Insert(c)
		if err != nil {
			err = tran.Rollback()
			if err != nil {
				return errors.Wrapf(err, "创建流程节点实例回滚事物发生错误")
			}
			return errors.Wrapf(err, "插入流程节点候选人数据发生错误")
		}
	}

	err = tran.Commit()
	if err != nil {
		return errors.Wrapf(err, "创建流程节点实例提交事物发生错误")
	}
	return nil
}

// UpdateNodeInstance 更新节点实例信息
func (a *Flow) UpdateNodeInstance(recordID string, info map[string]interface{}) error {
	_, err := a.db.UpdateByPK(schema.NodeInstancesTableName, db.M{"record_id": recordID}, db.M(info))
	if err != nil {
		return errors.Wrapf(err, "更新节点实例信息发生错误")
	}
	return nil
}

// CheckFlowInstanceTodo 检查流程实例待办事项
func (a *Flow) CheckFlowInstanceTodo(flowInstanceID string) (bool, error) {
	query := fmt.Sprintf("SELECT count(*) FROM %s WHERE deleted=0 AND status=1 AND flow_instance_id=?", schema.NodeInstancesTableName)
	n, err := a.db.SelectInt(query, flowInstanceID)
	if err != nil {
		return false, errors.Wrapf(err, "检查流程待办事项发生错误")
	}
	return n > 0, nil
}

// UpdateFlowInstance 更新流程实例信息
func (a *Flow) UpdateFlowInstance(recordID string, info map[string]interface{}) error {
	_, err := a.db.UpdateByPK(schema.FlowInstancesTableName, db.M{"record_id": recordID}, db.M(info))
	if err != nil {
		return errors.Wrapf(err, "更新流程实例信息发生错误")
	}
	return nil
}

// CreateFlowInstance 创建流程实例
func (a *Flow) CreateFlowInstance(flowInstance *schema.FlowInstances, nodeInstances ...*schema.NodeInstances) error {
	tran, err := a.db.Begin()
	if err != nil {
		return errors.Wrapf(err, "创建流程实例开启事物发生错误")
	}

	err = tran.Insert(flowInstance)
	if err != nil {
		err = tran.Rollback()
		if err != nil {
			return errors.Wrapf(err, "创建流程实例回滚事物发生错误")
		}
		return errors.Wrapf(err, "插入流程实例数据发生错误")
	}

	for _, n := range nodeInstances {
		err = tran.Insert(n)
		if err != nil {
			err = tran.Rollback()
			if err != nil {
				return errors.Wrapf(err, "创建流程实例回滚事物发生错误")
			}
			return errors.Wrapf(err, "插入流程节点实例数据发生错误")
		}
	}

	err = tran.Commit()
	if err != nil {
		return errors.Wrapf(err, "创建流程实例提交事物发生错误")
	}
	return nil
}

// QueryNodeCandidates 查询节点候选人
func (a *Flow) QueryNodeCandidates(nodeInstanceID string) ([]*schema.NodeCandidates, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted=0 AND node_instance_id=?", schema.NodeCandidatesTableName)

	var items []*schema.NodeCandidates
	_, err := a.db.Select(&items, query, nodeInstanceID)
	if err != nil {
		return nil, errors.Wrapf(err, "查询节点候选人发生错误")
	}

	return items, nil
}

// QueryTodoNodeInstances 查询用户的待办节点实例数据
func (a *Flow) QueryTodoNodeInstances(flowID, userID string) ([]*schema.NodeInstances, error) {
	query := fmt.Sprintf(`
SELECT *
FROM %s
WHERE deleted = 0 AND status = 1 AND record_id IN (SELECT node_instance_id
                                                   FROM %s
                                                   WHERE deleted = 0 AND candidate_id = ?) AND
      flow_instance_id IN (SELECT record_id
                           FROM %s
                           WHERE deleted = 0 AND status = 1 AND flow_id = ?)
		`, schema.NodeInstancesTableName, schema.NodeCandidatesTableName, schema.FlowInstancesTableName)

	var items []*schema.NodeInstances
	_, err := a.db.Select(&items, query, userID, flowID)
	if err != nil {
		return nil, errors.Wrapf(err, "查询用户的待办节点实例数据发生错误")
	}
	return items, nil
}
