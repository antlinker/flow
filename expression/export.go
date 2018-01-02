package expression

import "context"

var (
	defaultExp = CreateExecer("")
)

// ScriptImportAlias 设置导入脚本模块　并定义别名
func ScriptImportAlias(model, alias string) {
	defaultExp.ScriptImportAlias(model, alias)
}

// ScriptImport 设置导入脚本模块
func ScriptImport(model string) {
	defaultExp.ScriptImport(model)
}

// SetLibs 修改默认脚本库目录
func SetLibs(libs string) {
	defaultExp.SetLibs(libs)
}

// Exec 执行表达式
func Exec(ctx ExpContext, exp string) (*OutData, error) {
	return defaultExp.Exec(ctx, exp)
}

// ExecParam 执行表达式
func ExecParam(ctx context.Context, exp string, vars map[string]interface{}) (*OutData, error) {
	ectx := CreateExpContext(ctx)
	for key, v := range vars {
		ectx.AddVar(key, v)
	}
	return defaultExp.Exec(ectx, exp)
}

// ExecParamBool 执行表达式　返回布尔型
func ExecParamBool(ctx context.Context, exp string, vars map[string]interface{}) (bool, error) {
	return Bool(ExecParam(ctx, exp, vars))
}

// ExecParamSliceStr 执行表达式，返回字符串切片
func ExecParamSliceStr(ctx context.Context, exp string, vars map[string]interface{}) ([]string, error) {
	return SliceStr(ExecParam(ctx, exp, vars))
}

// ExecPredefineVar 执行表达式,传入预编译参数
func ExecPredefineVar(ctx context.Context, exp string, key string, predefinestr string) (*OutData, error) {
	ectx := CreateExpContext(ctx)
	ectx.PredefinedVar(key, predefinestr)
	return defaultExp.Exec(ectx, exp)
}

// ExecPredefineVarBool 执行表达式,传入预编译参数 返回布尔值
func ExecPredefineVarBool(ctx context.Context, exp string, key string, predefinestr string) (bool, error) {
	return Bool(ExecPredefineVar(ctx, exp, key, predefinestr))
}

// ExecPredefineVarSliceStr 执行表达式,传入预编译参数 返回字符串切片
func ExecPredefineVarSliceStr(ctx context.Context, exp string, key string, predefinestr string) ([]string, error) {
	return SliceStr(ExecPredefineVar(ctx, exp, key, predefinestr))
}

// ExecBool 执行表达式返回布尔值
func ExecBool(ctx ExpContext, exp string) (bool, error) {
	return Bool(defaultExp.Exec(ctx, exp))
}

// Bool 返回布尔值
func Bool(d *OutData, err ...error) (bool, error) {

	if len(err) > 1 && err[0] != nil {
		return false, err[0]
	}
	return d.Bool()
}

// SliceStr 返回字符串切片
func SliceStr(d *OutData, err ...error) ([]string, error) {

	if len(err) > 1 && err[0] != nil {
		return nil, err[0]
	}
	return d.SliceStr()
}
