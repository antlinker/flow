package schema

// 定义表名
const (
	FlowsTableName           = "f_flows"
	FlowNodesTableName       = "f_flow_nodes"
	NodeRoutersTableName     = "f_node_routers"
	NodeAssignmentsTableName = "f_node_assignments"
	FlowInstancesTableName   = "f_flow_instances"
	NodeInstancesTableName   = "f_node_instances"
	NodeCandidatesTableName  = "f_node_candidates"
)

// Flows 流程表
type Flows struct {
	ID       int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`     // 唯一标识(自增ID)
	RecordID string `db:"record_id,size:36" structs:"record_id" json:"record_id"` // 记录内码(uuid)
	Code     string `db:"code,size:50" structs:"code" json:"code"`                // 流程编号(唯一)
	Name     string `db:"name,size:50" structs:"name" json:"name"`                // 流程名称
	Version  int64  `db:"version" structs:"version" json:"version"`               // 版本号
	TypeCode string `db:"type_code,size:50" structs:"type_code" json:"type_code"` // 流程类型编号
	XML      string `db:"xml,size:1024" structs:"xml" json:"xml"`                 // XML数据
	Memo     string `db:"memo,size:255" structs:"memo" json:"memo"`               // 流程备注
	Creator  string `db:"creator,size:36" structs:"creator" json:"creator"`       // 创建人
	Created  int64  `db:"created" structs:"created" json:"created"`               // 创建时间戳
	Updator  string `db:"updator,size:36" structs:"updator" json:"updator"`       // 更新人
	Updated  int64  `db:"updated" structs:"updated" json:"updated"`               // 更新时间戳
	Deletor  string `db:"deletor,size:36" structs:"deletor" json:"deletor"`       // 删除人
	Deleted  int64  `db:"deleted" structs:"deleted" json:"deleted"`               // 删除时间戳
}

// FlowNodes 流程节点表
type FlowNodes struct {
	ID       int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`     // 唯一标识(自增ID)
	RecordID string `db:"record_id,size:36" structs:"record_id" json:"record_id"` // 记录内码(uuid)
	FlowID   string `db:"flow_id,size:36" structs:"flow_id" json:"flow_id"`       // 流程内码(flows.record_id)
	Code     string `db:"code,size:50" structs:"code" json:"code"`                // 节点编号
	Name     string `db:"name,size:50" structs:"name" json:"name"`                // 节点名称
	TypeCode string `db:"type_code,size:50" structs:"type_code" json:"type_code"` // 节点类型编号
	OrderNum string `db:"order_num,size:10" structs:"order_num" json:"order_num"` // 排序值
	FormID   string `db:"form_id,size:36" structs:"form_id" json:"form_id"`       // 表单内码
	Creator  string `db:"creator,size:36" structs:"creator" json:"creator"`       // 创建人
	Created  int64  `db:"created" structs:"created" json:"created"`               // 创建时间戳
	Updator  string `db:"updator,size:36" structs:"updator" json:"updator"`       // 更新人
	Updated  int64  `db:"updated" structs:"updated" json:"updated"`               // 更新时间戳
	Deletor  string `db:"deletor,size:36" structs:"deletor" json:"deletor"`       // 删除人
	Deleted  int64  `db:"deleted" structs:"deleted" json:"deleted"`               // 删除时间戳
}

// NodeRouters 节点路由表
type NodeRouters struct {
	ID              int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`                     // 唯一标识(自增ID)
	RecordID        string `db:"record_id,size:36" structs:"record_id" json:"record_id"`                 // 记录内码(uuid)
	SourceNodeID    string `db:"source_node_id,size:36" structs:"source_node_id" json:"source_node_id"`  // 源节点内码
	TargetNodeID    string `db:"target_node_id,size:36" structs:"target_node_id" json:"target_node_id"`  // 目标节点内码
	Expression      string `db:"expression,size:1024" structs:"expression" json:"expression"`            // 条件表达式(使用qlang作为表达式脚本语言(返回值bool))
	Explain         string `db:"explain,size:255" structs:"explain" json:"explain"`                      // 说明
	IsDefaultTarget int64  `db:"is_default_target" structs:"is_default_target" json:"is_default_target"` // 是否是默认节点(1:是 2:否)
	Creator         string `db:"creator,size:36" structs:"creator" json:"creator"`                       // 创建人
	Created         int64  `db:"created" structs:"created" json:"created"`                               // 创建时间戳
	Updator         string `db:"updator,size:36" structs:"updator" json:"updator"`                       // 更新人
	Updated         int64  `db:"updated" structs:"updated" json:"updated"`                               // 更新时间戳
	Deletor         string `db:"deletor,size:36" structs:"deletor" json:"deletor"`                       // 删除人
	Deleted         int64  `db:"deleted" structs:"deleted" json:"deleted"`                               // 删除时间戳
}

// NodeAssignments 节点指派表
type NodeAssignments struct {
	ID         int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`          // 唯一标识(自增ID)
	RecordID   string `db:"record_id,size:36" structs:"record_id" json:"record_id"`      // 记录内码(uuid)
	NodeID     string `db:"node_id,size:36" structs:"node_id" json:"node_id"`            // 节点内码(flow_nodes.record_id)
	Expression string `db:"expression,size:1024" structs:"expression" json:"expression"` // 执行表达式(基于qlang可提供多种内置函数支持，支持SQL查询)
	Creator    string `db:"creator,size:36" structs:"creator" json:"creator"`            // 创建人
	Created    int64  `db:"created" structs:"created" json:"created"`                    // 创建时间戳
	Updator    string `db:"updator,size:36" structs:"updator" json:"updator"`            // 更新人
	Updated    int64  `db:"updated" structs:"updated" json:"updated"`                    // 更新时间戳
	Deletor    string `db:"deletor,size:36" structs:"deletor" json:"deletor"`            // 删除人
	Deleted    int64  `db:"deleted" structs:"deleted" json:"deleted"`                    // 删除时间戳
}

// FlowInstances 流程实例表
type FlowInstances struct {
	ID         int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`     // 唯一标识(自增ID)
	RecordID   string `db:"record_id,size:36" structs:"record_id" json:"record_id"` // 记录内码(uuid)
	FlowID     string `db:"flow_id,size:36" structs:"flow_id" json:"flow_id"`       // 流程内码(flows.record_id)
	Status     int64  `db:"status" structs:"status" json:"status"`                  // 流程状态(0:未开始 1:进行中 2:暂停 3:已停止 9:已结束)
	Launcher   string `db:"launcher,size:36" structs:"launcher" json:"launcher"`    // 发起人
	LaunchTime int64  `db:"launch_time" structs:"launch_time" json:"launch_time"`   // 发起时间
	Creator    string `db:"creator,size:36" structs:"creator" json:"creator"`       // 创建人
	Created    int64  `db:"created" structs:"created" json:"created"`               // 创建时间戳
	Updator    string `db:"updator,size:36" structs:"updator" json:"updator"`       // 更新人
	Updated    int64  `db:"updated" structs:"updated" json:"updated"`               // 更新时间戳
	Deletor    string `db:"deletor,size:36" structs:"deletor" json:"deletor"`       // 删除人
	Deleted    int64  `db:"deleted" structs:"deleted" json:"deleted"`               // 删除时间戳
}

// NodeInstances 节点实例表
type NodeInstances struct {
	ID             int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`                          // 唯一标识(自增ID)
	RecordID       string `db:"record_id,size:36" structs:"record_id" json:"record_id"`                      // 记录内码(uuid)
	FlowInstanceID string `db:"flow_instance_id,size:36" structs:"flow_instance_id" json:"flow_instance_id"` // 流程实例内码(flows.record_id)
	NodeID         string `db:"node_id,size:36" structs:"node_id" json:"node_id"`                            // 节点内码
	Processor      string `db:"processor,size:36" structs:"processor" json:"processor"`                      // 处理人
	ProcessTime    int64  `db:"process_time" structs:"process_time" json:"process_time"`                     // 处理时间(秒时间戳)
	InputData      string `db:"input_data,size:1024" structs:"input_data" json:"input_data"`                 // 输入数据
	OutData        string `db:"out_data,size:1024" structs:"out_data" json:"out_data"`                       // 输出数据
	Status         int64  `db:"status" structs:"status" json:"status"`                                       // 处理状态(1:待处理 2:已完成)
	Creator        string `db:"creator,size:36" structs:"creator" json:"creator"`                            // 创建人
	Created        int64  `db:"created" structs:"created" json:"created"`                                    // 创建时间戳
	Updator        string `db:"updator,size:36" structs:"updator" json:"updator"`                            // 更新人
	Updated        int64  `db:"updated" structs:"updated" json:"updated"`                                    // 更新时间戳
	Deletor        string `db:"deletor,size:36" structs:"deletor" json:"deletor"`                            // 删除人
	Deleted        int64  `db:"deleted" structs:"deleted" json:"deleted"`                                    // 删除时间戳
}

// NodeCandidates 节点候选人表
type NodeCandidates struct {
	ID             int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`                          // 唯一标识(自增ID)
	RecordID       string `db:"record_id,size:36" structs:"record_id" json:"record_id"`                      // 记录内码(uuid)
	NodeInstanceID string `db:"node_instance_id,size:36" structs:"node_instance_id" json:"node_instance_id"` // 节点实例内码
	CandidateID    string `db:"candidate_id,size:36" structs:"candidate_id" json:"candidate_id"`             // 候选人ID(根据节点指派表达式生成)
	Creator        string `db:"creator,size:36" structs:"creator" json:"creator"`                            // 创建人
	Created        int64  `db:"created" structs:"created" json:"created"`                                    // 创建时间戳
	Updator        string `db:"updator,size:36" structs:"updator" json:"updator"`                            // 更新人
	Updated        int64  `db:"updated" structs:"updated" json:"updated"`                                    // 更新时间戳
	Deletor        string `db:"deletor,size:36" structs:"deletor" json:"deletor"`                            // 删除人
	Deleted        int64  `db:"deleted" structs:"deleted" json:"deleted"`                                    // 删除时间戳
}
