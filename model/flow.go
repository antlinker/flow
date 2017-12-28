package model

import (
	"gitee.com/antlinker/flow/schema"
	"gitee.com/antlinker/flow/service/db"
)

// Flow 流程管理
type Flow struct {
	db *db.DB
}

// Init 初始化
func (a *Flow) Init(db *db.DB) *Flow {
	db.AddTableWithName(schema.Flows{}, schema.FlowsTableName)
	a.db = db
	return a
}
