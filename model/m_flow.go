package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/antlinker/flow/schema"
	"github.com/antlinker/flow/service/db"
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
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted=0 AND flag=1 AND status=1 AND code=? ORDER BY version DESC LIMIT 1", schema.FlowTableName)

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

// QueryFlowByCode 根据流程编号查询流程数据
func (a *Flow) QueryFlowByCode(flowCode string) ([]*schema.Flow, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted=0 AND flag=1 AND status=1 AND code=? ORDER BY version DESC", schema.FlowTableName)

	var items []*schema.Flow
	_, err := a.DB.Select(&items, query, flowCode)
	if err != nil {
		return nil, errors.Wrapf(err, "根据流程编号查询流程数据发生错误")
	}

	return items, nil
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
func (a *Flow) QueryTodo(typeCode, flowCode, userID string) ([]*schema.FlowTodoResult, error) {
	var args []interface{}
	query := fmt.Sprintf(`
		SELECT
		  ni.record_id,
		  ni.flow_instance_id,
		  ni.input_data,
		  ni.node_id,
		  f.data 'form_data',
		  f.type_code 'form_type',
		  fi.launcher,
		  fi.launch_time,
			n.code 'node_code',
			n.name 'node_name',
			fw.name 'flow_name'
		FROM %s ni
		  JOIN %s fi ON ni.flow_instance_id = fi.record_id AND fi.deleted = ni.deleted
		  LEFT JOIN %s n ON ni.node_id = n.record_id AND n.deleted = ni.deleted
		  LEFT JOIN %s f ON n.form_id = f.record_id AND f.deleted = n.deleted
			LEFT JOIN %s fw ON n.flow_id = fw.record_id AND fw.deleted=n.deleted
		WHERE ni.deleted = 0 AND ni.status = 1 AND fi.status = 1 AND ni.record_id IN (SELECT node_instance_id FROM %s WHERE deleted = 0 AND candidate_id = ?)
		`, schema.NodeInstanceTableName, schema.FlowInstanceTableName, schema.NodeTableName, schema.FormTableName, schema.FlowTableName, schema.NodeCandidateTableName)

	args = append(args, userID)
	if typeCode != "" {
		query = fmt.Sprintf("%s AND fi.flow_id IN (SELECT record_id FROM %s WHERE deleted=0 AND flag=1 AND type_code=?)", query, schema.FlowTableName)
		args = append(args, typeCode)
	} else if flowCode != "" {
		query = fmt.Sprintf("%s AND fi.flow_id IN (SELECT record_id FROM %s WHERE deleted=0 AND flag=1 AND code=?)", query, schema.FlowTableName)
		args = append(args, flowCode)
	}
	query = fmt.Sprintf("%s ORDER BY ni.id", query)

	var items []*schema.FlowTodoResult
	_, err := a.DB.Select(&items, query, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "查询用户的待办数据发生错误")
	}
	return items, nil
}

// GetTodoByID 根据ID获取待办
func (a *Flow) GetTodoByID(nodeInstanceID string) (*schema.FlowTodoResult, error) {
	query := fmt.Sprintf(`
		SELECT
		  ni.record_id,
		  ni.flow_instance_id,
		  ni.input_data,
		  ni.node_id,
		  f.data 'form_data',
		  f.type_code 'form_type',
		  fi.launcher,
		  fi.launch_time,
			n.code 'node_code',
			n.name 'node_name',
			fw.name 'flow_name'
		FROM %s ni
		  JOIN %s fi ON ni.flow_instance_id = fi.record_id AND fi.deleted = ni.deleted
		  LEFT JOIN %s n ON ni.node_id = n.record_id AND n.deleted = ni.deleted
		  LEFT JOIN %s f ON n.form_id = f.record_id AND f.deleted = n.deleted
			LEFT JOIN %s fw ON n.flow_id = fw.record_id AND fw.deleted=n.deleted
		WHERE ni.deleted = 0 AND ni.status = 1 AND fi.status = 1 AND ni.record_id=?
		`, schema.NodeInstanceTableName, schema.FlowInstanceTableName, schema.NodeTableName, schema.FormTableName, schema.FlowTableName)

	var item schema.FlowTodoResult
	err := a.DB.SelectOne(&item, query, nodeInstanceID)
	if err != nil {
		return nil, errors.Wrapf(err, "根据ID获取待办发生错误")
	}
	return &item, nil
}

// GetDoneByID 根据ID获取已办
func (a *Flow) GetDoneByID(nodeInstanceID string) (*schema.FlowDoneResult, error) {
	table := fmt.Sprintf(`%s ni
		JOIN %s fi ON ni.flow_instance_id = fi.record_id AND fi.deleted = ni.deleted
		LEFT JOIN %s n ON ni.node_id = n.record_id AND n.deleted = ni.deleted
		LEFT JOIN %s f ON n.form_id = f.record_id AND f.deleted = n.deleted
		LEFT JOIN %s fw ON n.flow_id = fw.record_id AND fw.deleted=n.deleted`, schema.NodeInstanceTableName, schema.FlowInstanceTableName, schema.NodeTableName, schema.FormTableName, schema.FlowTableName)

	where := "WHERE ni.deleted = 0 AND ni.status = 2 AND n.type_code='userTask' AND ni.record_id=?"

	fieldsSelect := `
	ni.record_id,
	ni.flow_instance_id,
	ni.out_data,
	ni.process_time,
	f.data 'form_data',
	f.type_code 'form_type',
	fi.status 'flow_status',
	fi.launcher,
	fi.launch_time,
	n.record_id 'node_id',
	n.name 'node_name',
	fw.name 'flow_name'`

	query := fmt.Sprintf("SELECT %s FROM %s %s", fieldsSelect, table, where)

	var item schema.FlowDoneResult
	err := a.DB.SelectOne(&item, query, nodeInstanceID)
	if err != nil {
		return nil, errors.Wrapf(err, "根据ID获取已办发生错误")
	}

	return &item, nil
}

// QueryDone 查询用户的已办数据
func (a *Flow) QueryDone(typeCode, flowCode, userID string, lastTime int64, count int) ([]*schema.FlowDoneResult, error) {
	table := fmt.Sprintf(`%s ni
		JOIN %s fi ON ni.flow_instance_id = fi.record_id AND fi.deleted = ni.deleted
		LEFT JOIN %s n ON ni.node_id = n.record_id AND n.deleted = ni.deleted
		LEFT JOIN %s f ON n.form_id = f.record_id AND f.deleted = n.deleted
		LEFT JOIN %s fw ON n.flow_id = fw.record_id AND fw.deleted=n.deleted`, schema.NodeInstanceTableName, schema.FlowInstanceTableName, schema.NodeTableName, schema.FormTableName, schema.FlowTableName)

	where := "WHERE ni.deleted = 0 AND ni.status = 2 AND n.type_code='userTask' AND ni.processor=?"
	args := []interface{}{userID}

	if typeCode != "" {
		where = fmt.Sprintf("%s AND fi.flow_id IN (SELECT record_id FROM %s WHERE deleted=0 AND flag=1 AND type_code=?)", where, schema.FlowTableName)
		args = append(args, typeCode)
	} else if flowCode != "" {
		where = fmt.Sprintf("%s AND fi.flow_id IN (SELECT record_id FROM %s WHERE deleted=0 AND flag=1 AND code=?)", where, schema.FlowTableName)
		args = append(args, flowCode)
	}

	if lastTime > 0 {
		where = fmt.Sprintf("%s AND ni.process_time<?", where)
		args = append(args, lastTime)
	}

	fieldsSelect := `
	ni.record_id,
	ni.flow_instance_id,
	ni.out_data,
	ni.process_time,
	f.data 'form_data',
	f.type_code 'form_type',
	fi.status 'flow_status',
	fi.launcher,
	fi.launch_time,
	n.record_id 'node_id',
	n.name 'node_name',
	fw.name 'flow_name'`

	query := fmt.Sprintf("SELECT %s FROM %s %s ORDER BY ni.process_time DESC LIMIT %d", fieldsSelect, table, where, count)

	var items []*schema.FlowDoneResult
	_, err := a.DB.Select(&items, query, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "查询用户的已办数据发生错误")
	}
	return items, nil
}

// GetDoneCount 获取已办数量
func (a *Flow) GetDoneCount(userID string) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE deleted=0 AND status=2 AND processor=?", schema.NodeInstanceTableName)

	n, err := a.DB.SelectInt(query, userID)
	if err != nil {
		return 0, errors.Wrapf(err, "获取已办数量发生错误")
	}
	return n, nil
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

	_, err = tran.Exec(fmt.Sprintf("UPDATE %s SET deleted=? WHERE deleted=0 AND node_id IN(SELECT record_id FROM %s WHERE deleted=0 AND flow_id=?)", schema.NodePropertyTableName, schema.NodeTableName), ctimeUnix, flowID)
	if err != nil {
		_ = tran.Rollback()
		return errors.Wrapf(err, "删除流程节点属性发生错误")
	}

	_, err = tran.Exec(fmt.Sprintf("UPDATE %s SET deleted=? WHERE deleted=0 AND flow_id=?", schema.NodeTableName), ctimeUnix, flowID)
	if err != nil {
		_ = tran.Rollback()
		return errors.Wrapf(err, "删除流程节点发生错误")
	}

	_, err = tran.Exec(fmt.Sprintf("UPDATE %s SET deleted=? WHERE deleted=0 AND flow_id=?", schema.FormTableName), ctimeUnix, flowID)
	if err != nil {
		_ = tran.Rollback()
		return errors.Wrapf(err, "删除流程表单发生错误")
	}

	err = tran.Commit()
	if err != nil {
		return errors.Wrapf(err, "删除流程提交事物发生错误")
	}
	return nil
}

// QueryHistory 查询流程实例历史数据
func (a *Flow) QueryHistory(flowInstanceID string) ([]*schema.FlowHistoryResult, error) {
	query := fmt.Sprintf(`
		SELECT
		ni.record_id,
		ni.processor,
		ni.process_time,
		ni.input_data,
		ni.out_data,
		ni.status,
		n.record_id 'node_id',
		n.code 'node_code',
		n.name 'node_name',
		f.data 'form_data',
		f.type_code 'form_type'
		FROM %s ni JOIN %s n ON ni.node_id=n.record_id AND n.deleted=ni.deleted
		LEFT JOIN %s f ON n.form_id = f.record_id AND f.deleted = n.deleted
		WHERE ni.deleted=0 AND ni.flow_instance_id=? AND n.type_code='userTask'
		ORDER BY ni.status DESC,ni.process_time
		`, schema.NodeInstanceTableName, schema.NodeTableName, schema.FormTableName)

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

// QueryFlowIDsByType 根据类型查询流程ID列表
func (a *Flow) QueryFlowIDsByType(typeCode string) ([]string, error) {
	query := fmt.Sprintf("SELECT record_id FROM %s WHERE deleted=0 AND flag=1 AND status=1 AND type_code=?", schema.FlowTableName)

	var items []*schema.Flow
	_, err := a.DB.Select(&items, query, typeCode)
	if err != nil {
		return nil, errors.Wrapf(err, "根据类型查询流程ID列表发生错误")
	}

	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.RecordID
	}
	return ids, nil
}

// QueryFlowByIDs 根据流程ID查询流程数据
func (a *Flow) QueryFlowByIDs(flowIDs []string) ([]*schema.FlowQueryResult, error) {
	query := fmt.Sprintf("SELECT code,MAX(version)'version' FROM %s WHERE deleted=0 AND flag=1 AND status=1 AND record_id IN(?)  GROUP BY code ORDER BY code", schema.FlowTableName)

	query, args, err := a.DB.In(query, flowIDs)
	if err != nil {
		return nil, errors.Wrapf(err, "根据流程ID查询流程数据发生错误")
	}

	var items []*schema.Flow
	_, err = a.DB.Select(&items, query, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "根据流程ID查询流程数据发生错误")
	} else if len(items) == 0 {
		return nil, nil
	}

	result := make([]*schema.FlowQueryResult, len(items))
	for i, item := range items {
		flowResult, verr := a.GetFlowQueryResultByCodeAndVersion(item.Code, item.Version)
		if verr != nil {
			return nil, verr
		}
		result[i] = flowResult
	}

	return result, nil
}

// GetFlowFormByNodeID 获取流程节点表单
func (a *Flow) GetFlowFormByNodeID(nodeID string) (*schema.Form, error) {
	node, err := a.GetNode(nodeID)
	if err != nil {
		return nil, err
	} else if node == nil || node.FormID == "" {
		return nil, nil
	}

	return a.GetForm(node.FormID)
}

// GetNodeByFlowAndTypeCode 根据流程ID和节点类型获取节点数据
func (a *Flow) GetNodeByFlowAndTypeCode(flowID, typeCode string) (*schema.Node, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted=0 AND flow_id=? AND type_code=?", schema.NodeTableName)

	var item schema.Node
	err := a.DB.SelectOne(&item, query, flowID, typeCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "根据流程ID和节点类型获取节点数据发生错误")
	}

	return &item, nil
}

// GetForm 获取流程表单
func (a *Flow) GetForm(formID string) (*schema.Form, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted=0 AND record_id=?", schema.FormTableName)

	var item schema.Form
	err := a.DB.SelectOne(&item, query, formID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "获取流程表单发生错误")
	}
	return &item, nil
}

// Update 更新流程信息
func (a *Flow) Update(recordID string, info map[string]interface{}) error {
	_, err := a.DB.UpdateByPK(schema.FlowTableName, db.M{"record_id": recordID}, db.M(info))
	if err != nil {
		return errors.Wrapf(err, "更新流程信息发生错误")
	}
	return nil
}

// QueryNodeProperty 查询节点属性
func (a *Flow) QueryNodeProperty(nodeID string) ([]*schema.NodeProperty, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted=0 AND node_id=?", schema.NodePropertyTableName)

	var items []*schema.NodeProperty
	_, err := a.DB.Select(&items, query, nodeID)
	if err != nil {
		return nil, errors.Wrapf(err, "查询节点属性发生错误")
	}

	return items, nil
}

// CreateNodeTiming 创建定时节点
func (a *Flow) CreateNodeTiming(item *schema.NodeTiming) error {
	err := a.DB.Insert(item)
	if err != nil {
		return errors.Wrapf(err, "创建节点定时发生错误")
	}
	return nil
}

// UpdateNodeTiming 更新定时节点
func (a *Flow) UpdateNodeTiming(nodeInstanceID string, info map[string]interface{}) error {
	_, err := a.DB.UpdateByPK(schema.NodeTimingTableName, db.M{"node_instance_id": nodeInstanceID}, db.M(info))
	if err != nil {
		return errors.Wrapf(err, "更新节点定时发生错误")
	}
	return nil
}

// QueryExpiredNodeTiming 查询到期的定时节点
func (a *Flow) QueryExpiredNodeTiming() ([]*schema.NodeTiming, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted=0 AND expired_at < ? ORDER BY expired_at", schema.NodeTimingTableName)

	var items []*schema.NodeTiming
	_, err := a.DB.Select(&items, query, time.Now().Unix())
	if err != nil {
		return nil, errors.Wrapf(err, "查询到期的节点定时发生错误")
	}
	return items, nil
}

// QueryLaunchFlowInstanceResult 查询发起的流程实例数据
func (a *Flow) QueryLaunchFlowInstanceResult(launcher, typeCode, flowCode string, lastID int64, count int) ([]*schema.FlowInstanceResult, error) {
	var args []interface{}
	query := fmt.Sprintf("SELECT fi.id,fi.record_id,fi.flow_id,fi.status,fi.launcher,fi.launch_time,f.code 'flow_code',f.name 'flow_name' FROM %s fi LEFT JOIN %s f ON fi.flow_id=f.record_id AND f.deleted=0 WHERE fi.deleted=0", schema.FlowInstanceTableName, schema.FlowTableName)
	query = fmt.Sprintf("%s AND fi.launcher=?", query)
	args = append(args, launcher)

	if typeCode != "" {
		query = fmt.Sprintf("%s AND f.type_code=?", query)
		args = append(args, typeCode)
	} else if flowCode != "" {
		query = fmt.Sprintf("%s AND f.code=?", query)
		args = append(args, flowCode)
	}

	if lastID > 0 {
		query = fmt.Sprintf("%s AND fi.id<?", query)
		args = append(args, lastID)
	}

	query = fmt.Sprintf("%s ORDER BY fi.id DESC LIMIT %d", query, count)

	var items []*schema.FlowInstanceResult
	_, err := a.DB.Select(&items, query, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "查询发起的流程实例数据发生错误")
	}

	return items, nil
}

// QueryTodoFlowInstanceResult 查询待办的流程实例数据
func (a *Flow) QueryTodoFlowInstanceResult(userID, typeCode, flowCode string, lastID int64, count int) ([]*schema.FlowInstanceResult, error) {
	var args []interface{}
	query := fmt.Sprintf("SELECT fi.id,fi.record_id,fi.flow_id,fi.status,fi.launcher,fi.launch_time,f.code 'flow_code',f.name 'flow_name' FROM %s fi LEFT JOIN %s f ON fi.flow_id=f.record_id AND f.deleted=0 WHERE fi.deleted=0 AND fi.status = 1", schema.FlowInstanceTableName, schema.FlowTableName)
	query = fmt.Sprintf("%s AND fi.record_id IN(SELECT flow_instance_id FROM %s WHERE deleted=0 AND status=1 AND record_id IN(SELECT node_instance_id FROM %s WHERE deleted=0 AND candidate_id=?))", query, schema.NodeInstanceTableName, schema.NodeCandidateTableName)
	args = append(args, userID)

	if typeCode != "" {
		query = fmt.Sprintf("%s AND f.type_code=?", query)
		args = append(args, typeCode)
	} else if flowCode != "" {
		query = fmt.Sprintf("%s AND f.code=?", query)
		args = append(args, flowCode)
	}

	if lastID > 0 {
		query = fmt.Sprintf("%s AND fi.id<?", query)
		args = append(args, lastID)
	}

	query = fmt.Sprintf("%s ORDER BY fi.id DESC LIMIT %d", query, count)

	var items []*schema.FlowInstanceResult
	_, err := a.DB.Select(&items, query, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "查询发起的流程实例数据发生错误")
	}

	return items, nil
}

// QueryHandleFlowInstanceResult 查询处理的流程实例结果
func (a *Flow) QueryHandleFlowInstanceResult(processor, typeCode, flowCode string, lastID int64, count int) ([]*schema.FlowInstanceResult, error) {
	var args []interface{}
	query := fmt.Sprintf("SELECT fi.id,fi.record_id,fi.flow_id,fi.status,fi.launcher,fi.launch_time,f.code 'flow_code',f.name 'flow_name' FROM %s fi LEFT JOIN %s f ON fi.flow_id=f.record_id AND f.deleted=0 WHERE fi.deleted=0", schema.FlowInstanceTableName, schema.FlowTableName)
	query = fmt.Sprintf("%s AND fi.launcher!=?", query)
	args = append(args, processor)
	query = fmt.Sprintf("%s AND fi.record_id IN(SELECT flow_instance_id FROM %s WHERE deleted=0 AND status=2 AND processor=?)", query, schema.NodeInstanceTableName)
	args = append(args, processor)

	if typeCode != "" {
		query = fmt.Sprintf("%s AND f.type_code=?", query)
		args = append(args, typeCode)
	} else if flowCode != "" {
		query = fmt.Sprintf("%s AND f.code=?", query)
		args = append(args, flowCode)
	}

	if lastID > 0 {
		query = fmt.Sprintf("%s AND fi.id<?", query)
		args = append(args, lastID)
	}

	query = fmt.Sprintf("%s ORDER BY fi.id DESC LIMIT %d", query, count)

	var items []*schema.FlowInstanceResult
	_, err := a.DB.Select(&items, query, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "查询发起的流程实例数据发生错误")
	}

	return items, nil
}

// QueryLastNodeInstance 查询节点实例
func (a *Flow) QueryLastNodeInstance(flowInstanceID string) (*schema.NodeInstance, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted=0 AND flow_instance_id=? ORDER BY id DESC LIMIT 1", schema.NodeInstanceTableName)

	var item schema.NodeInstance
	err := a.DB.SelectOne(&item, query, flowInstanceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "查询节点实例发生错误")
	}

	return &item, nil
}

// -----------------------------web查询操作(start)-------------------------------

// QueryAllFlowPage 查询流程分页数据
func (a *Flow) QueryAllFlowPage(params schema.FlowQueryParam, pageIndex, pageSize uint) (int64, []*schema.FlowQueryResult, error) {
	var (
		where = "WHERE deleted=0 AND flag=1"
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

	if v := params.TypeCode; v != "" {
		where = fmt.Sprintf("%s AND type_code=?", where)
		args = append(args, v)
	}

	if v := params.Status; v > 0 {
		where = fmt.Sprintf("%s AND status=?", where)
		args = append(args, v)
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

// QueryGroupFlowPage 查询流程分组分页数据
func (a *Flow) QueryGroupFlowPage(params schema.FlowQueryParam, pageIndex, pageSize uint) (int64, []*schema.FlowQueryResult, error) {
	var (
		where = "WHERE deleted=0 AND flag=1"
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

	if v := params.TypeCode; v != "" {
		where = fmt.Sprintf("%s AND type_code=?", where)
		args = append(args, v)
	}

	if v := params.Status; v > 0 {
		where = fmt.Sprintf("%s AND status=?", where)
		args = append(args, v)
	}

	query := fmt.Sprintf("SELECT code,MAX(version)'version' FROM %s %s GROUP BY code ORDER BY code", schema.FlowTableName, where)

	var items []*schema.Flow
	_, err := a.DB.Select(&items, query, args...)
	if err != nil {
		return 0, nil, errors.Wrapf(err, "查询分页数据发生错误")
	} else if len(items) == 0 {
		return 0, nil, nil
	}

	start := int((pageIndex - 1) * pageSize)
	end := int(start + int(pageSize))

	var data []*schema.Flow
	if l := len(items); l > start {
		if l > end {
			data = items[start:end]
		} else {
			data = items[start:]
		}
	}

	result := make([]*schema.FlowQueryResult, len(data))
	for i, item := range data {
		flowResult, verr := a.GetFlowQueryResultByCodeAndVersion(item.Code, item.Version)
		if verr != nil {
			return 0, nil, verr
		}
		result[i] = flowResult
	}

	return int64(len(items)), result, err
}

// GetFlowQueryResultByCodeAndVersion 根据编号和版本获取流程结果
func (a *Flow) GetFlowQueryResultByCodeAndVersion(code string, version int64) (*schema.FlowQueryResult, error) {
	query := fmt.Sprintf("SELECT id,record_id,created,code,name,version,type_code,status,memo FROM %s", schema.FlowTableName)
	query = fmt.Sprintf("%s WHERE deleted=0 AND flag=1 AND code=? AND version=?", query)

	var item schema.FlowQueryResult
	err := a.DB.SelectOne(&item, query, code, version)
	if err != nil {
		return nil, errors.Wrapf(err, "查询流程结果发生错误")
	}

	return &item, nil
}

// QueryFlowVersion 查询流程版本数据
func (a *Flow) QueryFlowVersion(code string) ([]*schema.FlowQueryResult, error) {
	query := fmt.Sprintf("SELECT id,record_id,created,code,name,version,type_code,status,memo FROM %s", schema.FlowTableName)
	query = fmt.Sprintf("%s WHERE deleted=0 AND flag=1 AND code=? ORDER BY version", query)

	var items []*schema.FlowQueryResult
	_, err := a.DB.Select(&items, query, code)
	if err != nil {
		return nil, errors.Wrapf(err, "查询流程版本数据发生错误")
	}

	return items, nil
}

// -----------------------------web查询操作(end)---------------------------------
