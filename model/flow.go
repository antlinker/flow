package model

import (
	"database/sql"
	"fmt"
	"time"

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
	db.AddTableWithName(schema.Flow{}, schema.FlowTableName)
	db.AddTableWithName(schema.Node{}, schema.NodeTableName)
	db.AddTableWithName(schema.NodeRouter{}, schema.NodeRouterTableName)
	db.AddTableWithName(schema.NodeAssignment{}, schema.NodeAssignmentTableName)
	db.AddTableWithName(schema.FlowInstance{}, schema.FlowInstanceTableName)
	db.AddTableWithName(schema.NodeInstance{}, schema.NodeInstanceTableName)
	db.AddTableWithName(schema.NodeCandidate{}, schema.NodeCandidateTableName)
	db.AddTableWithName(schema.Form{}, schema.FormTableName)
	db.AddTableWithName(schema.FormField{}, schema.FormFieldTableName)
	db.AddTableWithName(schema.FieldProperty{}, schema.FieldPropertyTableName)
	db.AddTableWithName(schema.FieldValidation{}, schema.FieldValidationTableName)

	a.db = db
	return a
}

// CreateFlowBasic 创建流程基础数据
func (a *Flow) CreateFlowBasic(flow *schema.Flow, nodes []*schema.Node, routers []*schema.NodeRouter, assignments []*schema.NodeAssignment) error {
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

// GetFlow 获取流程数据
func (a *Flow) GetFlow(recordID string) (*schema.Flow, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted=0 AND record_id=? LIMIT 1", schema.FlowTableName)

	var flow schema.Flow
	err := a.db.SelectOne(&flow, query, recordID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "获取流程数据发生错误")
	}

	return &flow, nil
}

// GetFlowByCode 根据编号查询流程数据
func (a *Flow) GetFlowByCode(code string) (*schema.Flow, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted=0 AND flag=1 AND code=? ORDER BY version DESC LIMIT 1", schema.FlowTableName)

	var flow schema.Flow
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
func (a *Flow) GetNode(recordID string) (*schema.Node, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted=0 AND record_id=?", schema.NodeTableName)

	var item schema.Node
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
func (a *Flow) GetNodeByCode(flowID, nodeCode string) (*schema.Node, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted=0 AND flow_id=? AND code=? ORDER BY order_num", schema.NodeTableName)

	var items []*schema.Node
	_, err := a.db.Select(&items, query, flowID, nodeCode)
	if err != nil {
		return nil, errors.Wrapf(err, "根据节点编号获取流程节点发生错误")
	} else if len(items) == 0 {
		return nil, nil
	}

	return items[0], nil
}

// GetFlowInstance 获取流程实例
func (a *Flow) GetFlowInstance(recordID string) (*schema.FlowInstance, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted=0 AND record_id=?", schema.FlowInstanceTableName)

	var item schema.FlowInstance
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
func (a *Flow) GetNodeInstance(recordID string) (*schema.NodeInstance, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted=0 AND record_id=?", schema.NodeInstanceTableName)

	var item schema.NodeInstance
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
func (a *Flow) QueryNodeRouters(sourceNodeID string) ([]*schema.NodeRouter, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted=0 AND source_node_id=?", schema.NodeRouterTableName)

	var items []*schema.NodeRouter
	_, err := a.db.Select(&items, query, sourceNodeID)
	if err != nil {
		return nil, errors.Wrapf(err, "查询节点路由发生错误")
	}

	return items, nil
}

// QueryNodeAssignments 查询节点指派
func (a *Flow) QueryNodeAssignments(nodeID string) ([]*schema.NodeAssignment, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted=0 AND node_id=?", schema.NodeAssignmentTableName)

	var items []*schema.NodeAssignment
	_, err := a.db.Select(&items, query, nodeID)
	if err != nil {
		return nil, errors.Wrapf(err, "查询节点指派发生错误")
	}

	return items, nil
}

// CreateNodeInstance 创建流程节点实例
func (a *Flow) CreateNodeInstance(nodeInstance *schema.NodeInstance, nodeCandidates []*schema.NodeCandidate) error {
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
	_, err := a.db.UpdateByPK(schema.NodeInstanceTableName, db.M{"record_id": recordID}, db.M(info))
	if err != nil {
		return errors.Wrapf(err, "更新节点实例信息发生错误")
	}
	return nil
}

// CheckFlowInstanceTodo 检查流程实例待办事项
func (a *Flow) CheckFlowInstanceTodo(flowInstanceID string) (bool, error) {
	query := fmt.Sprintf("SELECT count(*) FROM %s WHERE deleted=0 AND status=1 AND flow_instance_id=?", schema.NodeInstanceTableName)
	n, err := a.db.SelectInt(query, flowInstanceID)
	if err != nil {
		return false, errors.Wrapf(err, "检查流程待办事项发生错误")
	}
	return n > 0, nil
}

// UpdateFlowInstance 更新流程实例信息
func (a *Flow) UpdateFlowInstance(recordID string, info map[string]interface{}) error {
	_, err := a.db.UpdateByPK(schema.FlowInstanceTableName, db.M{"record_id": recordID}, db.M(info))
	if err != nil {
		return errors.Wrapf(err, "更新流程实例信息发生错误")
	}
	return nil
}

// CreateFlowInstance 创建流程实例
func (a *Flow) CreateFlowInstance(flowInstance *schema.FlowInstance, nodeInstances ...*schema.NodeInstance) error {
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
func (a *Flow) QueryNodeCandidates(nodeInstanceID string) ([]*schema.NodeCandidate, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted=0 AND node_instance_id=?", schema.NodeCandidateTableName)

	var items []*schema.NodeCandidate
	_, err := a.db.Select(&items, query, nodeInstanceID)
	if err != nil {
		return nil, errors.Wrapf(err, "查询节点候选人发生错误")
	}

	return items, nil
}

// QueryTodoNodeInstances 查询用户的待办节点实例数据
func (a *Flow) QueryTodoNodeInstances(flowID, userID string) ([]*schema.NodeInstance, error) {
	query := fmt.Sprintf(`
SELECT *
FROM %s
WHERE deleted = 0 AND status = 1 AND record_id IN (SELECT node_instance_id
                                                   FROM %s
                                                   WHERE deleted = 0 AND candidate_id = ?) AND
      flow_instance_id IN (SELECT record_id
                           FROM %s
                           WHERE deleted = 0 AND status = 1 AND flow_id = ?)
		`, schema.NodeInstanceTableName, schema.NodeCandidateTableName, schema.FlowInstanceTableName)

	var items []*schema.NodeInstance
	_, err := a.db.Select(&items, query, userID, flowID)
	if err != nil {
		return nil, errors.Wrapf(err, "查询用户的待办节点实例数据发生错误")
	}
	return items, nil
}

// QueryFlowPage 查询流程分页数据
func (a *Flow) QueryFlowPage(params schema.FlowQueryParam, pageIndex, pageSize uint) (int64, []*schema.FlowQueryResult, error) {
	var (
		where = "WHERE deleted=0"
		args  []interface{}
	)

	if code := params.Code; code != "" {
		where = fmt.Sprintf("%s AND code LIKE ?", where)
		args = append(args, "%"+code+"%")
	}

	if name := params.Name; name != "" {
		where = fmt.Sprintf("%s AND name LIKE ?", where)
		args = append(args, "%"+name+"%")
	}

	n, err := a.db.SelectInt(fmt.Sprintf("SELECT count(*) FROM %s %s", schema.FlowTableName, where), args...)
	if err != nil {
		return 0, nil, errors.Wrapf(err, "查询分页数据发生错误")
	} else if n == 0 {
		return 0, nil, nil
	}

	query := fmt.Sprintf("SELECT id,record_id,created,code,name,version FROM %s %s ORDER BY id DESC", schema.FlowTableName, where)
	if pageIndex > 0 && pageSize > 0 {
		query = fmt.Sprintf("%s limit %d,%d", query, (pageIndex-1)*pageSize, pageSize)
	}

	var items []*schema.FlowQueryResult
	_, err = a.db.Select(&items, query, args...)
	if err != nil {
		return 0, nil, errors.Wrapf(err, "查询分页数据发生错误")
	}

	return n, items, err
}

// DeleteFlow 删除流程
func (a *Flow) DeleteFlow(flowID string) error {
	tran, err := a.db.Begin()
	if err != nil {
		return errors.Wrapf(err, "删除流程开启事物发生错误")
	}

	ctimeUnix := time.Now().Unix()
	_, err = tran.Exec(fmt.Sprintf("UPDATE %s SET deleted=? WHERE deleted=0 AND record_id=?", schema.FlowTableName), ctimeUnix, flowID)
	if err != nil {
		_ = tran.Rollback()
		return errors.Wrapf(err, "删除流程发生错误")
	}

	_, err = tran.Exec(fmt.Sprintf("UPDATE %s SET deleted=? WHERE deleted=0 AND source_node_id IN(SELECT record_id FROM %s WHERE deleted=0 AND flow_id=?)", schema.NodeRouterTableName, schema.NodeTableName), ctimeUnix, flowID)
	if err != nil {
		_ = tran.Rollback()
		return errors.Wrapf(err, "删除流程节点路由发生错误")
	}

	_, err = tran.Exec(fmt.Sprintf("UPDATE %s SET deleted=? WHERE deleted=0 AND node_id IN(SELECT record_id FROM %s WHERE deleted=0 AND flow_id=?)", schema.NodeAssignmentTableName, schema.NodeTableName), ctimeUnix, flowID)
	if err != nil {
		_ = tran.Rollback()
		return errors.Wrapf(err, "删除流程节点指派发生错误")
	}

	_, err = tran.Exec(fmt.Sprintf("UPDATE %s SET deleted=? WHERE deleted=0 AND flow_id=?", schema.NodeTableName), ctimeUnix, flowID)
	if err != nil {
		_ = tran.Rollback()
		return errors.Wrapf(err, "删除流程节点发生错误")
	}

	err = tran.Commit()
	if err != nil {
		return errors.Wrapf(err, "删除流程提交事物发生错误")
	}
	return nil
}
