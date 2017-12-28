package flow

// Execer 表达式执行器
type Execer interface {
	// 执行表达式返回布尔类型的值
	ExecToBool(exp, params []byte) (bool, error)

	// 执行表达式返回字符串切片类型的值
	ExecToStringSlice(exp, params []byte) ([]string, error)
}
