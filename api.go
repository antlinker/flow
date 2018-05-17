package flow

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/antlinker/flow/schema"
	"github.com/teambition/gear"
)

// API 提供API管理
type API struct {
	engine *Engine
}

// Init 初始化
func (a *API) Init(engine *Engine) *API {
	a.engine = engine
	return a
}

// 获取分页的页索引
func (a *API) pageIndex(ctx *gear.Context) uint {
	if v := ctx.Query("current"); v != "" {
		i, _ := strconv.Atoi(v)
		return uint(i)
	}
	return 1
}

// 获取分页的页大小
func (a *API) pageSize(ctx *gear.Context) uint {
	if v := ctx.Query("pageSize"); v != "" {
		i, _ := strconv.Atoi(v)
		if i > 40 {
			i = 40
		}
		return uint(i)
	}
	return 10
}

// QueryFlowPage 查询流程分页数据
func (a *API) QueryFlowPage(ctx *gear.Context) error {
	pageIndex, pageSize := a.pageIndex(ctx), a.pageSize(ctx)
	params := schema.FlowQueryParam{
		Code: ctx.Query("code"),
		Name: ctx.Query("name"),
	}

	total, items, err := a.engine.flowBll.QueryAllFlowPage(params, pageIndex, pageSize)
	if err != nil {
		return gear.ErrInternalServerError.From(err)
	}

	response := map[string]interface{}{
		"list": items,
		"pagination": map[string]interface{}{
			"total":    total,
			"current":  pageIndex,
			"pageSize": pageSize,
		},
	}

	return ctx.JSON(http.StatusOK, response)
}

// GetFlow 获取流程数据
func (a *API) GetFlow(ctx *gear.Context) error {
	item, err := a.engine.flowBll.GetFlow(ctx.Param("id"))
	if err != nil {
		return gear.ErrInternalServerError.From(err)
	}
	return ctx.JSON(http.StatusOK, item)
}

type saveFlowRequest struct {
	XML string `json:"xml"`
}

func (a *saveFlowRequest) Validate() error {
	if len(a.XML) == 0 {
		return errors.New("请求含有空数据")
	}
	return nil
}

// SaveFlow 保存流程
func (a *API) SaveFlow(ctx *gear.Context) error {
	var req saveFlowRequest
	if err := ctx.ParseBody(&req); err != nil {
		return gear.ErrBadRequest.From(err)
	}

	_, err := a.engine.CreateFlow([]byte(req.XML))
	if err != nil {
		return gear.ErrInternalServerError.From(err)
	}
	return ctx.JSON(http.StatusOK, "ok")
}

// DeleteFlow 删除流程数据
func (a *API) DeleteFlow(ctx *gear.Context) error {
	err := a.engine.flowBll.DeleteFlow(ctx.Param("id"))
	if err != nil {
		return gear.ErrInternalServerError.From(err)
	}
	return ctx.JSON(http.StatusOK, "ok")
}
