package expression

import (
	"io"
	"text/template"
)

const (
	expEvnTplstr = `
{{range $k, $v := .Import}}import "{{$k}}"{{if ne $v ""}} as {{$v}}{{end}}
{{end}}
// 自动生成Execer 预编译环境
{{range .ExecerVar}}{{.Key}} = {{.Value}}
{{end}}
// 自动生成ExpContext 预编译环境
{{range .CtxVar}}{{.Key}} = {{.Value}}
{{end}}
{{.ResultKey}} = {{.Exp}}
`
)

var (
	expEvnTpl = template.Must(template.New("exetpl").Parse(expEvnTplstr))
)

type tplOption struct {
	Import    map[string]string
	ExecerVar []pairs
	CtxVar    []pairs
	ResultKey string
	Exp       string
}

func parseExeTpl(wr io.Writer, tpl *tplOption) {
	expEvnTpl.Execute(wr, tpl)
}
