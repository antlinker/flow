package schema

// 定义表名
const (
	FlowTableName            = "f_flow"
	NodeTableName            = "f_node"
	NodeRouterTableName      = "f_node_router"
	NodeAssignmentTableName  = "f_node_assignment"
	FlowInstanceTableName    = "f_flow_instance"
	NodeInstanceTableName    = "f_node_instance"
	NodeCandidateTableName   = "f_node_candidate"
	FormTableName            = "f_form"
	FormFieldTableName       = "f_form_field"
	FieldPropertyTableName   = "f_field_property"
	FieldValidationTableName = "f_field_validation"
)

// Flow 流程
type Flow struct {
	ID       int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`     // 唯一标识(自增ID)
	RecordID string `db:"record_id,size:36" structs:"record_id" json:"record_id"` // 记录内码(uuid)
	Code     string `db:"code,size:50" structs:"code" json:"code"`                // 流程编号
	Name     string `db:"name,size:50" structs:"name" json:"name"`                // 流程名称
	Version  int64  `db:"version" structs:"version" json:"version"`               // 版本号
	TypeCode string `db:"type_code,size:50" structs:"type_code" json:"type_code"` // 流程类型编号
	XML      string `db:"xml,size:1024" structs:"xml" json:"xml"`                 // XML数据
	Memo     string `db:"memo,size:255" structs:"memo" json:"memo"`               // 流程备注
	Flag     int64  `db:"flag" structs:"flag" json:"flag"`                        // 流程标志(1:主流程 2:子流程)
	ParentID string `db:"parent_id,size:36" structs:"parent_id" json:"parent_id"` // 父级流程内码
	Created  int64  `db:"created" structs:"created" json:"created"`               // 创建时间戳
	Updated  int64  `db:"updated" structs:"updated" json:"updated"`               // 更新时间戳
	Deleted  int64  `db:"deleted" structs:"deleted" json:"deleted"`               // 删除时间戳
}

// Node 流程节点
type Node struct {
	ID       int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`     // 唯一标识(自增ID)
	RecordID string `db:"record_id,size:36" structs:"record_id" json:"record_id"` // 记录内码(uuid)
	FlowID   string `db:"flow_id,size:36" structs:"flow_id" json:"flow_id"`       // 流程内码(flows.record_id)
	Code     string `db:"code,size:50" structs:"code" json:"code"`                // 节点编号
	Name     string `db:"name,size:50" structs:"name" json:"name"`                // 节点名称
	TypeCode string `db:"type_code,size:50" structs:"type_code" json:"type_code"` // 节点类型编号
	OrderNum string `db:"order_num,size:10" structs:"order_num" json:"order_num"` // 排序值
	FormID   string `db:"form_id,size:36" structs:"form_id" json:"form_id"`       // 表单内码
	Created  int64  `db:"created" structs:"created" json:"created"`               // 创建时间戳
	Updated  int64  `db:"updated" structs:"updated" json:"updated"`               // 更新时间戳
	Deleted  int64  `db:"deleted" structs:"deleted" json:"deleted"`               // 删除时间戳
}

// NodeRouter 节点路由
type NodeRouter struct {
	ID              int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`                     // 唯一标识(自增ID)
	RecordID        string `db:"record_id,size:36" structs:"record_id" json:"record_id"`                 // 记录内码(uuid)
	SourceNodeID    string `db:"source_node_id,size:36" structs:"source_node_id" json:"source_node_id"`  // 源节点内码
	TargetNodeID    string `db:"target_node_id,size:36" structs:"target_node_id" json:"target_node_id"`  // 目标节点内码
	Expression      string `db:"expression,size:1024" structs:"expression" json:"expression"`            // 条件表达式(使用qlang作为表达式脚本语言(返回值bool))
	Explain         string `db:"explain,size:255" structs:"explain" json:"explain"`                      // 说明
	IsDefaultTarget int64  `db:"is_default_target" structs:"is_default_target" json:"is_default_target"` // 是否是默认节点(1:是 2:否)
	Created         int64  `db:"created" structs:"created" json:"created"`                               // 创建时间戳
	Updated         int64  `db:"updated" structs:"updated" json:"updated"`                               // 更新时间戳
	Deleted         int64  `db:"deleted" structs:"deleted" json:"deleted"`                               // 删除时间戳
}

// NodeAssignment 节点指派
type NodeAssignment struct {
	ID         int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`          // 唯一标识(自增ID)
	RecordID   string `db:"record_id,size:36" structs:"record_id" json:"record_id"`      // 记录内码(uuid)
	NodeID     string `db:"node_id,size:36" structs:"node_id" json:"node_id"`            // 节点内码(flow_nodes.record_id)
	Expression string `db:"expression,size:1024" structs:"expression" json:"expression"` // 执行表达式(基于qlang可提供多种内置函数支持，支持SQL查询)
	Created    int64  `db:"created" structs:"created" json:"created"`                    // 创建时间戳
	Updated    int64  `db:"updated" structs:"updated" json:"updated"`                    // 更新时间戳
	Deleted    int64  `db:"deleted" structs:"deleted" json:"deleted"`                    // 删除时间戳
}

// FlowInstance 流程实例
type FlowInstance struct {
	ID         int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`     // 唯一标识(自增ID)
	RecordID   string `db:"record_id,size:36" structs:"record_id" json:"record_id"` // 记录内码(uuid)
	FlowID     string `db:"flow_id,size:36" structs:"flow_id" json:"flow_id"`       // 流程内码(flows.record_id)
	Status     int64  `db:"status" structs:"status" json:"status"`                  // 流程状态(0:未开始 1:进行中 2:暂停 3:已停止 9:已完成)
	Launcher   string `db:"launcher,size:36" structs:"launcher" json:"launcher"`    // 发起人
	LaunchTime int64  `db:"launch_time" structs:"launch_time" json:"launch_time"`   // 发起时间
	Created    int64  `db:"created" structs:"created" json:"created"`               // 创建时间戳
	Updated    int64  `db:"updated" structs:"updated" json:"updated"`               // 更新时间戳
	Deleted    int64  `db:"deleted" structs:"deleted" json:"deleted"`               // 删除时间戳
}

// NodeInstance 节点实例表
type NodeInstance struct {
	ID             int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`                          // 唯一标识(自增ID)
	RecordID       string `db:"record_id,size:36" structs:"record_id" json:"record_id"`                      // 记录内码(uuid)
	FlowInstanceID string `db:"flow_instance_id,size:36" structs:"flow_instance_id" json:"flow_instance_id"` // 流程实例内码(flows.record_id)
	NodeID         string `db:"node_id,size:36" structs:"node_id" json:"node_id"`                            // 节点内码
	Processor      string `db:"processor,size:36" structs:"processor" json:"processor"`                      // 处理人
	ProcessTime    int64  `db:"process_time" structs:"process_time" json:"process_time"`                     // 处理时间(秒时间戳)
	InputData      string `db:"input_data,size:1024" structs:"input_data" json:"input_data"`                 // 输入数据
	OutData        string `db:"out_data,size:1024" structs:"out_data" json:"out_data"`                       // 输出数据
	Status         int64  `db:"status" structs:"status" json:"status"`                                       // 处理状态(1:待处理 2:已完成)
	Created        int64  `db:"created" structs:"created" json:"created"`                                    // 创建时间戳
	Updated        int64  `db:"updated" structs:"updated" json:"updated"`                                    // 更新时间戳
	Deleted        int64  `db:"deleted" structs:"deleted" json:"deleted"`                                    // 删除时间戳
}

// NodeCandidate 节点候选人
type NodeCandidate struct {
	ID             int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`                          // 唯一标识(自增ID)
	RecordID       string `db:"record_id,size:36" structs:"record_id" json:"record_id"`                      // 记录内码(uuid)
	NodeInstanceID string `db:"node_instance_id,size:36" structs:"node_instance_id" json:"node_instance_id"` // 节点实例内码
	CandidateID    string `db:"candidate_id,size:36" structs:"candidate_id" json:"candidate_id"`             // 候选人ID(根据节点指派表达式生成)
	Created        int64  `db:"created" structs:"created" json:"created"`                                    // 创建时间戳
	Updated        int64  `db:"updated" structs:"updated" json:"updated"`                                    // 更新时间戳
	Deleted        int64  `db:"deleted" structs:"deleted" json:"deleted"`                                    // 删除时间戳
}

// Form 流程表单
type Form struct {
	ID       int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`     // 唯一标识(自增ID)
	RecordID string `db:"record_id,size:36" structs:"record_id" json:"record_id"` // 记录内码(uuid)
	FlowID   string `db:"flow_id,size:36" structs:"flow_id" json:"flow_id"`       // 流程内码(flows.record_id)
	Code     string `db:"code,size:50" structs:"code" json:"code"`                // 表单编号(唯一)
	Name     string `db:"name,size:50" structs:"name" json:"name"`                // 表单名称
	TypeCode string `db:"type_code,size:50" structs:"type_code" json:"type_code"` // 表单类型(URL:表单链接路径 META:表单元数据)
	Data     string `db:"data,size:1024" structs:"data" json:"data"`              // 表单数据
	Created  int64  `db:"created" structs:"created" json:"created"`               // 创建时间戳
	Updated  int64  `db:"updated" structs:"updated" json:"updated"`               // 更新时间戳
	Deleted  int64  `db:"deleted" structs:"deleted" json:"deleted"`               // 删除时间戳
}

// FormField 表单字段
type FormField struct {
	ID           int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`                  // 唯一标识(自增ID)
	RecordID     string `db:"record_id,size:36" structs:"record_id" json:"record_id"`              // 记录内码(uuid)
	FormID       string `db:"form_id,size:36" structs:"form_id" json:"form_id"`                    // 表单内码(flow_forms.record_id)
	Code         string `db:"code,size:50" structs:"code" json:"code"`                             // 字段编号
	Label        string `db:"label,size:50" structs:"label" json:"label"`                          // 字段标签
	TypeCode     string `db:"type_code,size:50" structs:"type_code" json:"type_code"`              // 字段类型
	DefaultValue string `db:"default_value,size:100" structs:"default_value" json:"default_value"` // 字段默认值
	Created      int64  `db:"created" structs:"created" json:"created"`                            // 创建时间戳
	Updated      int64  `db:"updated" structs:"updated" json:"updated"`                            // 更新时间戳
	Deleted      int64  `db:"deleted" structs:"deleted" json:"deleted"`                            // 删除时间戳
}

// FieldProperty 字段属性
type FieldProperty struct {
	ID       int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`     // 唯一标识(自增ID)
	RecordID string `db:"record_id,size:36" structs:"record_id" json:"record_id"` // 记录内码(uuid)
	FieldID  string `db:"field_id,size:36" structs:"field_id" json:"field_id"`    // 字段内码(form_fields.record_id)
	Code     string `db:"code,size:50" structs:"code" json:"code"`                // 属性编号
	Value    string `db:"value,size:100" structs:"value" json:"value"`            // 属性值
	Created  int64  `db:"created" structs:"created" json:"created"`               // 创建时间戳
	Updated  int64  `db:"updated" structs:"updated" json:"updated"`               // 更新时间戳
	Deleted  int64  `db:"deleted" structs:"deleted" json:"deleted"`               // 删除时间戳
}

// FieldValidation 字段校验
type FieldValidation struct {
	ID               int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`                              // 唯一标识(自增ID)
	RecordID         string `db:"record_id,size:36" structs:"record_id" json:"record_id"`                          // 记录内码(uuid)
	FieldID          string `db:"field_id,size:36" structs:"field_id" json:"field_id"`                             // 字段内码(form_fields.record_id)
	ConstraintName   string `db:"constraint_name,size:50" structs:"constraint_name" json:"constraint_name"`        // 约束名称
	ConstraintConfig string `db:"constraint_config,size:100" structs:"constraint_config" json:"constraint_config"` // 约束配置
	Created          int64  `db:"created" structs:"created" json:"created"`                                        // 创建时间戳
	Updated          int64  `db:"updated" structs:"updated" json:"updated"`                                        // 更新时间戳
	Deleted          int64  `db:"deleted" structs:"deleted" json:"deleted"`                                        // 删除时间戳
}

// FlowQueryParam 流程查询参数
type FlowQueryParam struct {
	Code string `db:"code" structs:"code" json:"code"` // 流程编号
	Name string `db:"name" structs:"name" json:"name"` // 流程名称
}

// FlowQueryResult 流程查询结果
type FlowQueryResult struct {
	ID       int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`     // 唯一标识(自增ID)
	RecordID string `db:"record_id,size:36" structs:"record_id" json:"record_id"` // 记录内码(uuid)
	Code     string `db:"code,size:50" structs:"code" json:"code"`                // 流程编号(唯一)
	Name     string `db:"name,size:50" structs:"name" json:"name"`                // 流程名称
	Version  int64  `db:"version" structs:"version" json:"version"`               // 版本号
	Created  int64  `db:"created" structs:"created" json:"created"`               // 创建时间戳
}
