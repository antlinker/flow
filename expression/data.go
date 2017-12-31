package expression

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"qlang.io/spec"
)

// OutData 输出数据
type OutData struct {
	Result interface{} // 数据结果
}

// IsUndefined 有未定义变量
func (d *OutData) IsUndefined() bool {
	return d != nil && d.Result == spec.Undefined
}

// IsNil 检测返回结果是不是空值
func (d *OutData) IsNil() bool {
	return d == nil || d.Result == nil
}

// String 获取字符串类型数据
func (d OutData) String() (string, error) {
	if d.IsUndefined() {
		return "", errors.Errorf("未定义变量：spec.Undefined")
	}
	if d.IsNil() {
		return "", nil
	}

	return result2string(d.Result)
}

// Bool 获取布尔类型数据
func (d OutData) Bool() (bool, error) {
	if d.IsUndefined() {
		return false, errors.Errorf("未定义变量：spec.Undefined")
	}
	if d.IsNil() {
		return false, nil
	}
	return result2bool(d.Result)
}

// Int 获取整数类型数据
func (d OutData) Int() (int, error) {
	if d.IsUndefined() {
		return 0, errors.Errorf("未定义变量：spec.Undefined")
	}
	if d.IsNil() {
		return 0, nil
	}
	return result2int(d.Result)
}

// SliceStr 获取字符串切片类型数据
func (d OutData) SliceStr() ([]string, error) {
	if d.IsUndefined() {
		return nil, errors.Errorf("未定义变量：spec.Undefined")
	}
	if d.IsNil() {
		return nil, nil
	}
	r, ok := d.Result.([]string)
	if ok {
		return r, nil
	}
	return nil, errors.Errorf("返回值的类型错误:%s", r)
}

// Float 获取浮点数类型数据
func (d OutData) Float() (float64, error) {
	if d.IsUndefined() {
		return 0, errors.Errorf("未定义变量：spec.Undefined")
	}
	if d.IsNil() {
		return 0, nil
	}
	return result2float(d.Result)
}
func result2string(result interface{}) (string, error) {

	return fmt.Sprintf("%v", result), nil
}
func result2float(result interface{}) (float64, error) {

	switch result.(type) {
	case bool:
		if result.(bool) {
			return 1, nil
		}
		return 0, nil
	case int, int8, int16, int32, int64, uint, uint16, uint32, uint64:
		return float64(result.(int)), nil
	case float32, float64:
		return result.(float64), nil
	case string:
		return strconv.ParseFloat(result.(string), 64)

	case uint8:
		return float64(result.(uint8)), nil
	default:
		return 0, fmt.Errorf("不能处理的类型:%s", result)
	}
	//return 0, fmt.Errorf("返回值的类型错误:%s", r)
}
func result2int(result interface{}) (int, error) {

	switch result.(type) {
	case bool:
		if result.(bool) {
			return 1, nil
		}
		return 0, nil
	case int, int8, int16, int32, int64, uint, uint16, uint32, uint64:
		return result.(int), nil
	case float32, float64:
		return int(result.(float64)), nil
	case string:
		return strconv.Atoi(result.(string))

	case uint8:
		return int(result.(uint8)), nil
	default:
		v := reflect.ValueOf(result)

		switch v.Kind() {
		case reflect.Slice, reflect.Map:
			return v.Len(), nil
		default:
			if !v.IsValid() {
				return 0, nil
			}
			if v.IsNil() {
				return 0, nil
			}
			return 0, fmt.Errorf("不能处理的类型:%s", v)
		}
	}
}

func result2bool(result interface{}) (bool, error) {

	//	fmt.Printf("type:%s\n", result)

	switch result.(type) {
	case bool:
		return result.(bool), nil
	case int, int8, int16, int32, int64, uint, uint16, uint32, uint64:
		return result.(int) != 0, nil
	case float32, float64:
		return result.(float64) != 0, nil
	case string:
		r := strings.ToLower(result.(string))
		return result.(string) != "" && r != "false" && r != "off", nil

	case uint8:
		return result.(uint8) != 0, nil
	default:
		v := reflect.ValueOf(result)

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
