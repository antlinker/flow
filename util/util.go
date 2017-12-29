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

func StringToInt(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
