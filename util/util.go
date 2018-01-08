package util

import (
	"strconv"

	"github.com/fatih/structs"
	"github.com/satori/go.uuid"
)

// UUID 获取UUID
func UUID() string {
	return uuid.NewV4().String()
}

// StructToMap 转换struct为map
func StructToMap(s interface{}) map[string]interface{} {
	return structs.Map(s)
}

// StringToInt 字符串转数值
func StringToInt(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// StringRemoveDuplicates 字符串切片去重
func StringRemoveDuplicates(strs []string) (result []string) {
	for i, l := 0, len(strs); i < l; i++ {
		var exist bool
		for j := 0; j < i; j++ {
			if strs[j] == strs[i] {
				exist = true
				break
			}
		}
		if !exist {
			result = append(result, strs[i])
		}
	}

	return
}
