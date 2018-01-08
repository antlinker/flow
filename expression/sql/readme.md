# sql 扩展模块使用

``` go

import (
    "gitee.com/antlinker/flow/expression"
    "gitee.com/antlinker/flow/expression/sql"
)
// db 为已经初始化后的数据库连接
sql.Reg(db)

exp := expression.CreateExecer("")
// sql 表示成当前目录下开始导入 sql/sql.ql脚本，如果
exp.ScriptImportAlias("sql/sql.ql", "sql")

// 这样就可以在脚本中使用 sql.Query(query,args...) sql.Count(query,args...) sql.One(query,args...) sql.querySliceStr 四个函数
// 也可以使用sqlctx.QueryDB(ctx,db,query,args) sqlctx.CountDB(ctx,db,query,args) sqlctx.OneDB(ctx,db,query,args)
// 也可以使用sqlctx.Query(ctx,query,args) sqlctx.Count(ctx,query,args) sqlctx.One(ctx,query,args)
```