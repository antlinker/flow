package flow

import (
	"context"
)

// Renderer 表单渲染器
type Renderer interface {
	Render(context.Context, *NodeFormResult) ([]byte, error)
}
