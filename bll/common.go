package bll

import (
	"gitee.com/antlinker/flow/model"
	"gitee.com/antlinker/flow/service/db"
)

// Bll 业务处理
type Bll struct {
	Models *model.All
}

// All 业务处理集合
type All struct {
	Models *model.All
	Flow   *Flow
}

// Init 初始化业务处理集合
func (a *All) Init(db *db.DB) *All {
	a.Models = new(model.All).Init(db)

	bll := &Bll{Models: a.Models}
	a.Flow = &Flow{bll}
	return a
}
