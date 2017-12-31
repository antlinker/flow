package builtin

import (
	"fmt"

	"qlang.io/cl/qlang"
)

// Reg 注册内置函数，不注册不能使用
func init() {
	qlang.Import("", map[string]interface{}{
		"SliceStr": SliceStr,
		"Slice":    Slice,
	})

}

func SliceStr(in []map[string]interface{}, key string) (out []string) {
	for _, d := range in {
		out = append(out, fmt.Sprintf("%v", d[key]))
	}
	return
}

func Slice(in []map[string]interface{}, key string) (out []interface{}) {
	for _, d := range in {
		out = append(out, d)
	}
	return
}
