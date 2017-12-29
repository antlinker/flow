package builtin

import (
	"fmt"
	"reflect"
	"strings"
	"sync/atomic"

	"qlang.io/cl/qlang"
	"qlang.io/spec"
)

var (
	counter = int64(0)
)

func creResultKey() string {
	c := atomic.AddInt64(&counter, 1)
	return fmt.Sprintf("__result%d__", c)
}

// ResultBool 表达式执行返回布尔型
func ResultBool(scriptCode string) bool {
	key := creResultKey()
	r, err := execbool([]byte(key+" = "+scriptCode), key)
	if err != nil {
		panic(fmt.Sprintf("执行失败:%v", err))
	}
	return r
}

// ResultString 表达式执行返回字符串
func ResultString(scriptCode string) string {
	key := creResultKey()
	out, err := exec([]byte(key+"="+scriptCode), key)
	if err != nil {
		panic(fmt.Sprintf("执行失败:%v", err))
	}
	return fmt.Sprintf("%v", out)
}

// ResultStringSlice 表达式执行返回字符串切片
func ResultStringSlice(scriptCode string) []string {
	key := creResultKey()
	out, err := exec([]byte(key+"="+scriptCode), key)
	if err != nil {
		panic(fmt.Sprintf("执行失败:%v", err))
	}

	r, ok := out.([]string)
	if ok {
		return r
	}
	panic(fmt.Errorf("输出格式错误：%s", out))
}

// ResultIntSlice 表达式执行返回整数切片
func ResultIntSlice(scriptCode string) []int {
	key := creResultKey()
	out, err := exec([]byte(key+" = "+scriptCode), key)
	if err != nil {
		panic(fmt.Sprintf("执行失败:%v", err))
	}

	r, ok := out.([]int)
	if ok {
		return r
	}
	panic(fmt.Errorf("输出格式错误：%s", out))
}

// ResultFloatSlice 表达式执行返回整数切片
func ResultFloatSlice(scriptCode string) ([]float64, error) {
	key := creResultKey()
	out, err := exec([]byte(key+"="+scriptCode), key)
	if err != nil {
		return nil, err
	}

	r, ok := out.([]float64)
	if ok {
		return r, nil
	}
	return nil, fmt.Errorf("输出格式错误：%s", out)
}
func execbool(scriptCode []byte, resultKey string) (bool, error) {

	out, err := exec(scriptCode, resultKey)
	if err != nil {
		// 错误处理
		fmt.Printf("error:%s\n", err)
		return false, err
	}
	if out == spec.Undefined {
		return false, nil
	}
	fmt.Printf("type:%s\n", out)

	switch out.(type) {
	case bool:
		return out.(bool), nil
	case int, int8, int16, int32, int64, uint, uint16, uint32, uint64:
		return out.(int) != 0, nil
	case float32, float64:
		return out.(float64) != 0, nil
	case string:
		r := strings.ToLower(out.(string))
		return out.(string) != "" && r != "false" && r != "off", nil

	case uint8:
		return out.(uint8) != 0, nil
	default:
		v := reflect.ValueOf(out)

		switch v.Kind() {
		case reflect.Slice, reflect.Map:
			return v.Len() > 0, nil
		default:
			if !v.IsValid() {
				return false, nil
			}
			if v.IsNil() {
				return false, nil
			}
			return false, fmt.Errorf("不能处理的类型:%s", v)
		}
	}
}

func exec(scriptCode []byte, resultKey string) (interface{}, error) {
	ql := qlang.New()

	err := ql.Exec(scriptCode, "")
	if err != nil {
		// 错误处理
		return false, err
	}
	out := ql.Var(resultKey)
	return out, nil
}
