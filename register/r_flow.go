package register

import (
	"gitee.com/antlinker/flow/schema"
	"gitee.com/antlinker/flow/service/db"
)

// FlowDBMap 注册流程相关的数据库映射
func FlowDBMap(db *db.DB) {
	db.AddTableWithName(schema.Flow{}, schema.FlowTableName)
	db.AddTableWithName(schema.Node{}, schema.NodeTableName)
	db.AddTableWithName(schema.NodeRouter{}, schema.NodeRouterTableName)
	db.AddTableWithName(schema.NodeAssignment{}, schema.NodeAssignmentTableName)
	db.AddTableWithName(schema.FlowInstance{}, schema.FlowInstanceTableName)
	db.AddTableWithName(schema.NodeInstance{}, schema.NodeInstanceTableName)
	db.AddTableWithName(schema.NodeCandidate{}, schema.NodeCandidateTableName)
	db.AddTableWithName(schema.Form{}, schema.FormTableName)
	db.AddTableWithName(schema.FormField{}, schema.FormFieldTableName)
	db.AddTableWithName(schema.FieldOption{}, schema.FieldOptionTableName)
	db.AddTableWithName(schema.FieldProperty{}, schema.FieldPropertyTableName)
	db.AddTableWithName(schema.FieldValidation{}, schema.FieldValidationTableName)
}
