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
	DB *db.DB `inject:""`
}

// CreateFlow 创建流程数据
func (a *Flow) CreateFlow(flow *schema.Flow, nodes *schema.NodeOperating, forms *schema.FormOperating) error {
	tran, err := a.DB.Begin()
	if err != nil {
		return errors.Wrapf(err, "创建流程基础数据开启事物发生错误")
	}

	err = tran.Insert(flow)
	if err != nil {
		_ = tran.Rollback()
		return errors.Wrapf(err, "插入流程数据发生错误")
	}

	if list := nodes.All(); len(list) > 0 {
		err = tran.Insert(list...)
		if err != nil {
			_ = tran.Rollback()
			return errors.Wrapf(err, "插入节点数据发生错误")
		}
	}

	if list := forms.All(); len(list) > 0 {
		err = tran.Insert(list...)
		if err != nil {
			_ = tran.Rollback()
			return errors.Wrapf(err, "插入表单数据发生错误")
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
	err := a.DB.SelectOne(&flow, query, recordID)
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
	err := a.DB.SelectOne(&flow, query, code)
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
	err := a.DB.SelectOne(&item, query, recordID)
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
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted=0 AND flow_id=? AND code=? ORDER BY order_num LIMIT 1", schema.NodeTableName)

	var item schema.Node
	err := a.DB.SelectOne(&item, query, flowID, nodeCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "根据节点编号获取流程节点发生错误")
	}

	return &item, nil
}

// GetFlowInstance 获取流程实例
func (a *Flow) GetFlowInstance(recordID string) (*schema.FlowInstance, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted=0 AND record_id=? LIMIT 1", schema.FlowInstanceTableName)

	var item schema.FlowInstance
	err := a.DB.SelectOne(&item, query, recordID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "获取流程实例发生错误")
	}

	return &item, nil
}

// GetFlowInstanceByNode 根据节点实例获取流程实例
func (a *Flow) GetFlowInstanceByNode(nodeInstanceID string) (*schema.FlowInstance, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted=0 AND record_id IN (SELECT flow_instance_id FROM %s WHERE deleted=0 AND record_id=?) LIMIT 1", schema.FlowInstanceTableName, schema.NodeInstanceTableName)

	var item schema.FlowInstance
	err := a.DB.SelectOne(&item, query, nodeInstanceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "根据节点实例获取流程实例发生错误")
	}

	return &item, nil
}

// GetNodeInstance 获取流程节点实例
func (a *Flow) GetNodeInstance(recordID string) (*schema.NodeInstance, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted=0 AND record_id=? LIMIT 1", schema.NodeInstanceTableName)

	var item schema.NodeInstance
	err := a.DB.SelectOne(&item, query, recordID)
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
	_, err := a.DB.Select(&items, query, sourceNodeID)
	if err != nil {
		return nil, errors.Wrapf(err, "查询节点路由发生错误")
	}

	return items, nil
}

// QueryNodeAssignments 查询节点指派
func (a *Flow) QueryNodeAssignments(nodeID string) ([]*schema.NodeAssignment, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted=0 AND node_id=?", schema.NodeAssignmentTableName)

	var items []*schema.NodeAssignment
	_, err := a.DB.Select(&items, query, nodeID)
	if err != nil {
		return nil, errors.Wrapf(err, "查询节点指派发生错误")
	}

	return items, nil
}

// CreateNodeInstance 创建流程节点实例
func (a *Flow) CreateNodeInstance(nodeInstance *schema.NodeInstance, nodeCandidates []*schema.NodeCandidate) error {
	tran, err := a.DB.Begin()
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
	_, err := a.DB.UpdateByPK(schema.NodeInstanceTableName, db.M{"record_id": recordID}, db.M(info))
	if err != nil {
		return errors.Wrapf(err, "更新节点实例信息发生错误")
	}
	return nil
}

// CheckFlowInstanceTodo 检查流程实例待办事项
func (a *Flow) CheckFlowInstanceTodo(flowInstanceID string) (bool, error) {
	query := fmt.Sprintf("SELECT count(*) FROM %s WHERE deleted=0 AND status=1 AND flow_instance_id=?", schema.NodeInstanceTableName)
	n, err := a.DB.SelectInt(query, flowInstanceID)
	if err != nil {
		return false, errors.Wrapf(err, "检查流程待办事项发生错误")
	}
	return n > 0, nil
}

// UpdateFlowInstance 更新流程实例信息
func (a *Flow) UpdateFlowInstance(recordID string, info map[string]interface{}) error {
	_, err := a.DB.UpdateByPK(schema.FlowInstanceTableName, db.M{"record_id": recordID}, db.M(info))
	if err != nil {
		return errors.Wrapf(err, "更新流程实例信息发生错误")
	}
	return nil
}

// CreateFlowInstance 创建流程实例
func (a *Flow) CreateFlowInstance(flowInstance *schema.FlowInstance, nodeInstances ...*schema.NodeInstance) error {
	tran, err := a.DB.Begin()
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
	_, err := a.DB.Select(&items, query, nodeInstanceID)
	if err != nil {
		return nil, errors.Wrapf(err, "查询节点候选人发生错误")
	}

	return items, nil
}

// QueryTodo 查询用户的待办数据
func (a *Flow) QueryTodo(flowCode, userID string) ([]*schema.FlowTodoResult, error) {
	query := fmt.Sprintf(`
	SELECT
	  ni.record_id,
	  ni.flow_instance_id,
	  ni.input_data,
	  ni.node_id,
	  f.data 'form_data',
	  fi.launcher,
	  fi.launch_time,
		n.code 'node_code',
		n.name 'node_name'
	FROM %s ni
	  JOIN %s fi ON ni.flow_instance_id = fi.record_id AND fi.deleted = ni.deleted
	  LEFT JOIN %s n ON ni.node_id = n.record_id AND n.deleted = ni.deleted
	  LEFT JOIN %s f ON n.form_id = f.record_id AND f.deleted = n.deleted
	WHERE ni.deleted = 0 AND ni.status = 1 AND fi.status = 1 AND fi.flow_id IN (SELECT record_id FROM %s WHERE deleted=0 AND flag=1 AND code=?)
	      AND ni.record_id IN (SELECT node_instance_id
	                           FROM %s
	                           WHERE deleted = 0 AND candidate_id = ?)
	`, schema.NodeInstanceTableName, schema.FlowInstanceTableName, schema.NodeTableName, schema.FormTableName, schema.FlowTableName, schema.NodeCandidateTableName)

	var items []*schema.FlowTodoResult
	_, err := a.DB.Select(&items, query, flowCode, userID)
	if err != nil {
		return nil, errors.Wrapf(err, "查询用户的待办数据发生错误")
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

	n, err := a.DB.SelectInt(fmt.Sprintf("SELECT count(*) FROM %s %s", schema.FlowTableName, where), args...)
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
	_, err = a.DB.Select(&items, query, args...)
	if err != nil {
		return 0, nil, errors.Wrapf(err, "查询分页数据发生错误")
	}

	return n, items, err
}

// DeleteFlow 删除流程
func (a *Flow) DeleteFlow(flowID string) error {
	tran, err := a.DB.Begin()
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

// QueryHistory 查询流程实例历史数据
func (a *Flow) QueryHistory(flowInstanceID string) ([]*schema.FlowHistoryResult, error) {
	query := fmt.Sprintf("SELECT ni.record_id,ni.processor,ni.process_time,ni.out_data,ni.status,n.code 'node_code',n.name 'node_name' FROM %s ni JOIN %s n ON ni.node_id=n.record_id AND n.deleted=ni.deleted WHERE ni.deleted=0 AND ni.flow_instance_id=? AND n.type_code='userTask' ORDER BY ni.status DESC,ni.process_time", schema.NodeInstanceTableName, schema.NodeTableName)

	var items []*schema.FlowHistoryResult
	_, err := a.DB.Select(&items, query, flowInstanceID)
	if err != nil {
		return nil, errors.Wrapf(err, "查询流程实例历史数据发生错误")
	}
	return items, nil
}

// QueryDoneIDs 查询已办理的流程实例ID列表
func (a *Flow) QueryDoneIDs(flowCode, userID string) ([]string, error) {
	query := fmt.Sprintf("SELECT record_id FROM %s WHERE deleted=0 AND flow_id IN (SELECT record_id FROM %s WHERE deleted=0 AND flag=1 AND code=?) AND record_id IN(SELECT flow_instance_id FROM %s WHERE deleted=0 AND status=2 AND processor=?)", schema.FlowInstanceTableName, schema.FlowTableName, schema.NodeInstanceTableName)

	var items []*schema.FlowInstance
	_, err := a.DB.Select(&items, query, flowCode, userID)
	if err != nil {
		return nil, errors.Wrapf(err, "查询已办理的流程数据发生错误")
	}

	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.RecordID
	}

	return ids, nil
}
