package model

import (
	"gitee.com/antlinker/flow/service/db"
)

// All 模型集合
type All struct {
	Flow *Flow
}

// Init 初始化模型集合
func (a *All) Init(db *db.DB) *All {
	a.Flow = new(Flow).Init(db)

	// 创建不存在的业务表
	err := db.CreateTablesIfNotExists()
	if err != nil {
		panic(err)
	}
	return a
}
