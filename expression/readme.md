# expression 表达式解析模块

## 使用说明

``` go

    import (
        "ant-flow/expression"
    )
    // db 为已经初始化后的数据库连接
    sql.Reg(db)

    exp := expression.CreateExecer("")
    // sql 表示成当前目录下开始导入 sql/sql.ql脚本，如果
    exp.ScriptImportAlias("sql/sql.ql", "sql")

    // 这样就可以在脚本中使用 sql.Query(query,args...) sql.Count(query,args...) sql.One(query,args...) 三个函数
    // 也可以使用sqlctx.QueryDB(ctx,db,query,args) sqlctx.CountDB(ctx,db,query,args) sqlctx.OneDB(ctx,db,query,args)
    // 也可以使用sqlctx.Query(ctx,query,args) sqlctx.Count(ctx,query,args) sqlctx.One(ctx,query,args)


	exp.PredefinedJson("global", map[string]interface{}{
		"test_1": 1,
		"test_a": "a",
	})
	exp.PredefinedVar("fun1", `fn(a ) {
		return 1==a
	}`)
	exp.PredefinedVar("fun2", `fn(a ) {
		return 1==a
    }`)

    ectx := expression.CreateExpContext(context.Background())

	ectx.AddVar("ctx_10", 10)
	ectx.AddVar("ctx_a", "a")

	out, err :=exp.Exec(ectx, "1+1")
	if err != nil {
        //return 0, err
        // TODO
    }
    // 打印输出2
    fmt.Println(out.Int())

    out, err :=exp.Exec(ectx, "ctx_10*10")
    if err != nil {
        //return 0, err
        // TODO
    }

    // 打印输出10
    fmt.Println(out.Int())

    out, err =exp.Exec(ectx, `sql.Count("select * from table")`)
    if err != nil {
        //return 0, err
        // TODO
    }

    // 打印输出查询到的表记录数量
    fmt.Println(out.Int())

    out, err:=exp.Exec(ectx, `SliceStr(sql.Query("select name from table"),"name")`)
    if err != nil {
        //return 0, err
        // TODO
    }

    // 打印输出查询到的表记录数量
    fmt.Println(out.SliceStr())
    out, err:=exp.Exec(ectx, `sql.querySliceStr("select name from table","name")`)
    if err != nil {
        //return 0, err
        // TODO
    }

    // 打印输出查询到的表记录数量
    fmt.Println(out.SliceStr())
```