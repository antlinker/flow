package schema

// 定义表名
const (
	FlowTableName            = "f_flow"
	NodeTableName            = "f_node"
	NodeRouterTableName      = "f_node_router"
	NodeAssignmentTableName  = "f_node_assignment"
	NodePropertyTableName    = "f_node_property"
	FlowInstanceTableName    = "f_flow_instance"
	NodeInstanceTableName    = "f_node_instance"
	NodeTimingTableName      = "f_node_timing"
	NodeCandidateTableName   = "f_node_candidate"
	FormTableName            = "f_form"
	FormFieldTableName       = "f_form_field"
	FieldOptionTableName     = "f_field_option"
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
	Status   int    `db:"status" structs:"status" json:"status"`                  // 流程状态(1:正常 2:禁用)
	Created  int64  `db:"created" structs:"created" json:"created"`               // 创建时间戳
	Updated  int64  `db:"updated" structs:"updated" json:"updated"`               // 更新时间戳
	Deleted  int64  `db:"deleted" structs:"deleted" json:"deleted"`               // 删除时间戳
}

// Node 流程节点
type Node struct {
	ID       int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`     // 唯一标识(自增ID)
	RecordID string `db:"record_id,size:36" structs:"record_id" json:"record_id"` // 记录内码(uuid)
	FlowID   string `db:"flow_id,size:36" structs:"flow_id" json:"flow_id"`       // 流程内码
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
	NodeID     string `db:"node_id,size:36" structs:"node_id" json:"node_id"`            // 节点内码
	Expression string `db:"expression,size:1024" structs:"expression" json:"expression"` // 执行表达式(基于qlang可提供多种内置函数支持，支持SQL查询)
	Created    int64  `db:"created" structs:"created" json:"created"`                    // 创建时间戳
	Updated    int64  `db:"updated" structs:"updated" json:"updated"`                    // 更新时间戳
	Deleted    int64  `db:"deleted" structs:"deleted" json:"deleted"`                    // 删除时间戳
}

// NodeProperty 节点属性
type NodeProperty struct {
	ID       int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`     // 唯一标识(自增ID)
	RecordID string `db:"record_id,size:36" structs:"record_id" json:"record_id"` // 记录内码(uuid)
	NodeID   string `db:"node_id,size:36" structs:"node_id" json:"node_id"`       // 节点内码
	Name     string `db:"name,size:50" structs:"name" json:"name"`                // 属性名称
	Value    string `db:"value,size:255" structs:"value" json:"value"`            // 属性值
	Created  int64  `db:"created" structs:"created" json:"created"`               // 创建时间戳
	Updated  int64  `db:"updated" structs:"updated" json:"updated"`               // 更新时间戳
	Deleted  int64  `db:"deleted" structs:"deleted" json:"deleted"`               // 删除时间戳
}

// FlowInstance 流程实例
type FlowInstance struct {
	ID         int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`     // 唯一标识(自增ID)
	RecordID   string `db:"record_id,size:36" structs:"record_id" json:"record_id"` // 记录内码(uuid)
	FlowID     string `db:"flow_id,size:36" structs:"flow_id" json:"flow_id"`       // 流程内码
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
	FlowInstanceID string `db:"flow_instance_id,size:36" structs:"flow_instance_id" json:"flow_instance_id"` // 流程实例内码
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

// NodeTiming 节点定时
type NodeTiming struct {
	ID             int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`                  // 唯一标识(自增ID)
	NodeInstanceID string `db:"node_instance_id" structs:"node_instance_id" json:"node_instance_id"` // 节点实例ID
	Flag           string `db:"flag" structs:"flag" json:"flag"`                                     // 标志
	Processor      string `db:"processor,size:36" structs:"processor" json:"processor"`              // 处理人
	Input          string `db:"input,size:1024" structs:"input" json:"input"`                        // 输入数据
	ExpiredAt      int64  `db:"expired_at" structs:"expired_at" json:"expired_at"`                   // 过期时间戳
	Created        int64  `db:"created" structs:"created" json:"created"`                            // 创建时间戳
	Deleted        int64  `db:"deleted" structs:"deleted" json:"deleted"`                            // 删除时间戳
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
	FlowID   string `db:"flow_id,size:36" structs:"flow_id" json:"flow_id"`       // 流程内码
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
	FormID       string `db:"form_id,size:36" structs:"form_id" json:"form_id"`                    // 表单内码
	Code         string `db:"code,size:50" structs:"code" json:"code"`                             // 字段编号
	Label        string `db:"label,size:50" structs:"label" json:"label"`                          // 字段标签
	TypeCode     string `db:"type_code,size:50" structs:"type_code" json:"type_code"`              // 字段类型
	DefaultValue string `db:"default_value,size:100" structs:"default_value" json:"default_value"` // 字段默认值
	Created      int64  `db:"created" structs:"created" json:"created"`                            // 创建时间戳
	Updated      int64  `db:"updated" structs:"updated" json:"updated"`                            // 更新时间戳
	Deleted      int64  `db:"deleted" structs:"deleted" json:"deleted"`                            // 删除时间戳
}

// FieldOption 字段选项
type FieldOption struct {
	ID        int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`         // 唯一标识(自增ID)
	RecordID  string `db:"record_id,size:36" structs:"record_id" json:"record_id"`     // 记录内码(uuid)
	FieldID   string `db:"field_id,size:36" structs:"field_id" json:"field_id"`        // 字段内码
	ValueID   string `db:"value_id,size:50" structs:"value_id" json:"value_id"`        // 选项值ID
	ValueName string `db:"value_name,size:100" structs:"value_name" json:"value_name"` // 选项值名称
	Created   int64  `db:"created" structs:"created" json:"created"`                   // 创建时间戳
	Updated   int64  `db:"updated" structs:"updated" json:"updated"`                   // 更新时间戳
	Deleted   int64  `db:"deleted" structs:"deleted" json:"deleted"`                   // 删除时间戳
}

// FieldProperty 字段属性
type FieldProperty struct {
	ID       int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`     // 唯一标识(自增ID)
	RecordID string `db:"record_id,size:36" structs:"record_id" json:"record_id"` // 记录内码(uuid)
	FieldID  string `db:"field_id,size:36" structs:"field_id" json:"field_id"`    // 字段内码
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
	FieldID          string `db:"field_id,size:36" structs:"field_id" json:"field_id"`                             // 字段内码
	ConstraintName   string `db:"constraint_name,size:50" structs:"constraint_name" json:"constraint_name"`        // 约束名称
	ConstraintConfig string `db:"constraint_config,size:100" structs:"constraint_config" json:"constraint_config"` // 约束配置
	Created          int64  `db:"created" structs:"created" json:"created"`                                        // 创建时间戳
	Updated          int64  `db:"updated" structs:"updated" json:"updated"`                                        // 更新时间戳
	Deleted          int64  `db:"deleted" structs:"deleted" json:"deleted"`                                        // 删除时间戳
}

// FlowQueryParam 流程查询参数
type FlowQueryParam struct {
	Code     string // 流程编号
	Name     string // 流程名称
	TypeCode string // 流程类型编号
	Status   int    // 流程状态(1:正常 2:禁用)
}

// FlowQueryResult 流程查询结果
type FlowQueryResult struct {
	ID       int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`     // 唯一标识(自增ID)
	RecordID string `db:"record_id,size:36" structs:"record_id" json:"record_id"` // 记录内码(uuid)
	Code     string `db:"code,size:50" structs:"code" json:"code"`                // 流程编号(唯一)
	Name     string `db:"name,size:50" structs:"name" json:"name"`                // 流程名称
	Version  int64  `db:"version" structs:"version" json:"version"`               // 版本号
	TypeCode string `db:"type_code,size:50" structs:"type_code" json:"type_code"` // 流程类型编号
	Status   int    `db:"status" structs:"status" json:"status"`                  // 流程状态(1:正常 2:禁用)
	Created  int64  `db:"created" structs:"created" json:"created"`               // 创建时间戳
	Memo     string `db:"memo,size:255" structs:"memo" json:"memo"`               // 流程备注
}

// FlowTodoResult 流程待办结果
type FlowTodoResult struct {
	RecordID       string  `db:"record_id" structs:"record_id" json:"record_id"`                      // 节点实例内码
	FlowInstanceID string  `db:"flow_instance_id" structs:"flow_instance_id" json:"flow_instance_id"` // 流程实例内码
	FlowName       string  `db:"flow_name" structs:"flow_name" json:"flow_name"`                      // 流程名称
	NodeID         string  `db:"node_id" structs:"node_id" json:"node_id"`                            // 节点内码
	NodeCode       string  `db:"node_code" structs:"node_code" json:"node_code"`                      // 节点编号
	NodeName       string  `db:"node_name" structs:"node_name" json:"node_name"`                      // 节点名称
	InputData      string  `db:"input_data" structs:"input_data" json:"input_data"`                   // 输入数据
	Launcher       string  `db:"launcher" structs:"launcher" json:"launcher"`                         // 发起人
	LaunchTime     int64   `db:"launch_time" structs:"launch_time" json:"launch_time"`                // 发起时间
	FormType       *string `db:"form_type" structs:"form_type" json:"form_type"`                      // 表单类型
	FormData       *string `db:"form_data" structs:"form_data" json:"form_data"`                      // 表单数据
}

// FlowHistoryResult 流程历史结果
type FlowHistoryResult struct {
	RecordID    string  `db:"record_id,size:36" structs:"record_id" json:"record_id"`      // 记录内码(uuid)
	NodeID      string  `db:"node_id,size:36" structs:"node_id" json:"node_id"`            // 节点ID
	NodeCode    string  `db:"node_code,size:36" structs:"node_code" json:"node_code"`      // 节点编号
	NodeName    string  `db:"node_name,size:36" structs:"node_name" json:"node_name"`      // 节点名称
	Processor   string  `db:"processor,size:36" structs:"processor" json:"processor"`      // 处理人
	ProcessTime int64   `db:"process_time" structs:"process_time" json:"process_time"`     // 处理时间(秒时间戳)
	InputData   string  `db:"input_data,size:1024" structs:"input_data" json:"input_data"` // 输入数据
	OutData     string  `db:"out_data,size:1024" structs:"out_data" json:"out_data"`       // 输出数据
	Status      int64   `db:"status" structs:"status" json:"status"`                       // 处理状态(1:待处理 2:已完成)
	FormType    *string `db:"form_type" structs:"form_type" json:"form_type"`              // 表单类型
	FormData    *string `db:"form_data" structs:"form_data" json:"form_data"`              // 表单数据
}

// FlowDoneResult 流程已办结果
type FlowDoneResult struct {
	RecordID       string  `db:"record_id" structs:"record_id" json:"record_id"`                      // 节点实例内码
	FlowInstanceID string  `db:"flow_instance_id" structs:"flow_instance_id" json:"flow_instance_id"` // 流程实例内码
	FlowName       string  `db:"flow_name" structs:"flow_name" json:"flow_name"`                      // 流程名称
	FlowStatus     int64   `db:"flow_status" structs:"flow_status" json:"flow_status"`                // 流程状态
	ProcessTime    int64   `db:"process_time" structs:"process_time" json:"process_time"`             // 处理时间
	NodeID         string  `db:"node_id" structs:"node_id" json:"node_id"`                            // 节点ID
	NodeName       string  `db:"node_name" structs:"node_name" json:"node_name"`                      // 节点名称
	OutData        string  `db:"out_data" structs:"out_data" json:"out_data"`                         // 输出数据
	Launcher       string  `db:"launcher" structs:"launcher" json:"launcher"`                         // 发起人
	LaunchTime     int64   `db:"launch_time" structs:"launch_time" json:"launch_time"`                // 发起时间
	FormType       *string `db:"form_type" structs:"form_type" json:"form_type"`                      // 表单类型
	FormData       *string `db:"form_data" structs:"form_data" json:"form_data"`                      // 表单数据
}

// FlowInstanceResult 流程实例结果
type FlowInstanceResult struct {
	ID         int64  `db:"id" structs:"id" json:"id"`                              // 流程实例自增ID
	RecordID   string `db:"record_id,size:36" structs:"record_id" json:"record_id"` // 流程实例记录内码
	FlowID     string `db:"flow_id,size:36" structs:"flow_id" json:"flow_id"`       // 流程内码
	Status     int64  `db:"status" structs:"status" json:"status"`                  // 流程状态(0:未开始 1:进行中 2:暂停 3:已停止 9:已完成)
	Launcher   string `db:"launcher,size:36" structs:"launcher" json:"launcher"`    // 发起人
	LaunchTime int64  `db:"launch_time" structs:"launch_time" json:"launch_time"`   // 发起时间
	FlowCode   string `db:"flow_code,size:36" structs:"flow_code" json:"flow_code"` // 流程编号
	FlowName   string `db:"flow_name,size:36" structs:"flow_name" json:"flow_name"` // 流程名称
}

// NodeOperating 节点操作
type NodeOperating struct {
	NodeGroup       []*Node
	RouterGroup     []*NodeRouter
	AssignmentGroup []*NodeAssignment
	PropertyGroup   []*NodeProperty
}

// All 获取所有节点操作的组
func (a *NodeOperating) All() []interface{} {
	var group []interface{}

	for _, item := range a.NodeGroup {
		group = append(group, item)
	}
	for _, item := range a.RouterGroup {
		group = append(group, item)
	}
	for _, item := range a.AssignmentGroup {
		group = append(group, item)
	}
	for _, item := range a.PropertyGroup {
		group = append(group, item)
	}

	return group
}

// FormOperating 表单操作
type FormOperating struct {
	FormGroup            []*Form
	FormFieldGroup       []*FormField
	FieldOptionGroup     []*FieldOption
	FieldPropertyGroup   []*FieldProperty
	FieldValidationGroup []*FieldValidation
}

// All 获取所有表单操作的组
func (a *FormOperating) All() []interface{} {
	var group []interface{}

	for _, item := range a.FormGroup {
		group = append(group, item)
	}
	for _, item := range a.FormFieldGroup {
		group = append(group, item)
	}
	for _, item := range a.FieldOptionGroup {
		group = append(group, item)
	}
	for _, item := range a.FieldPropertyGroup {
		group = append(group, item)
	}
	for _, item := range a.FieldValidationGroup {
		group = append(group, item)
	}

	return group
}
