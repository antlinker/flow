package sql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/antlinker/flow/expression"

	"qlang.io/cl/qlang"
)

// execDB 从ctx获取数据库
// 没有数据库发出panic
func execDB(ctx context.Context) *sql.DB {
	db := expression.FromExpContextForDB(ctx)
	if db == nil {
		panic(fmt.Errorf("没有指定数据库不能查询"))
	}
	return db
}

// getDB 获取指定数据库
// 如果ctx没有指定数据库返回defaultDB
// 如果defaultDB为nil则发出panic
func getDB(ctx context.Context, defaultDB *sql.DB) *sql.DB {
	db := expression.FromExpContextForDB(ctx)
	if db == nil {
		db = defaultDB
	}
	if db == nil {
		panic(fmt.Errorf("没有指定数据库不能查询"))
	}
	return db
}

// Reg 注册数据库DB
// 有默认数据库操作
// 也支持多数据库
func Reg(defaultDB *sql.DB) {
	qlang.Import("sqlctx", map[string]interface{}{
		"QueryDB": QueryDB,
		"Query": func(ctx context.Context, query string, args ...interface{}) []map[string]interface{} {

			return QueryDB(ctx, getDB(ctx, defaultDB), query, args...)
		},
		"CountDB": QueryDBCount,
		"Count": func(ctx context.Context, query string, args ...interface{}) int {
			return QueryDBCount(ctx, getDB(ctx, defaultDB), query, args...)
		},
		"OneDB": QueryOneDB,
		"One": func(query string, ctx context.Context, args ...interface{}) map[string]interface{} {
			return QueryOneDB(ctx, getDB(ctx, defaultDB), query, args...)
		},
	})

}

// RegMoreDB 注册多数据库支持
// 没有默认数据库
func RegMoreDB() {
	qlang.Import("sqlctx", map[string]interface{}{
		"QueryDB": QueryDB,
		"Query": func(ctx context.Context, query string, args ...interface{}) []map[string]interface{} {
			return QueryDB(ctx, execDB(ctx), query, args...)
		},
		"CountDB": QueryDBCount,
		"Count": func(ctx context.Context, query string, args ...interface{}) int {
			return QueryDBCount(ctx, execDB(ctx), query, args...)
		},
		"OneDB": QueryOneDB,
		"One": func(query string, ctx context.Context, args ...interface{}) map[string]interface{} {
			return QueryOneDB(ctx, execDB(ctx), query, args...)
		},
	})

}

// QueryDB 查询sql返回的所有行
func QueryDB(ctx context.Context, db *sql.DB, query string, args ...interface{}) (out []map[string]interface{}) {
	// var rows *sql.Rows
	// var err error
	rows, err := db.QueryContext(ctx, query, args...)

	if err != nil {
		panic(fmt.Sprintf("查询失败:%s  %v ==> %v", query, args, err))
	}
	cols, err := rows.Columns() // Remember to check err afterwards
	if err != nil {
		panic(fmt.Sprintf("查询失败:%s  %v ==> %v", query, args, err))
	}
	vals := make([]interface{}, len(cols))
	for i := range cols {
		vals[i] = new(sql.RawBytes)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(vals...)
		if err != nil {
			panic(fmt.Sprintf("提取数据失败:%s  %v ==> %v", query, args, err))
		}
		vmap := make(map[string]interface{})
		for i, col := range cols {
			var s string
			rb, ok := vals[i].(*sql.RawBytes)
			if ok {
				s = string(*rb)
			}
			vmap[col] = s
		}
		out = append(out, vmap)
	}

	return
}

// QueryDBCount 查询sql匹配的条数
func QueryDBCount(ctx context.Context, db *sql.DB, query string, args ...interface{}) (count int) {

	query = "select count(*) num from (" + query + ")"

	rows, err := db.QueryContext(ctx, query, args...)

	if err != nil {
		panic(fmt.Sprintf("查询失败:%s  %v ==> %v", query, args, err))
	}

	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			panic(fmt.Sprintf("提取数据失败:%s  %v ==> %v", query, args, err))
		}
	}
	return
}

// QueryOneDB 查询sql返回的第一条记录
func QueryOneDB(ctx context.Context, db *sql.DB, query string, args ...interface{}) (out map[string]interface{}) {
	// var rows *sql.Rows
	// var err error
	if strings.Index(strings.ToLower(query), "limit") < 0 {
		query = "select * from (" + query + ") limit 1"
	}
	rows, err := db.QueryContext(ctx, query, args...)

	if err != nil {
		panic(fmt.Sprintf("查询失败:%s  %v ==> %v", query, args, err))
	}
	cols, err := rows.Columns() // Remember to check err afterwards
	if err != nil {
		panic(fmt.Sprintf("查询失败:%s  %v ==> %v", query, args, err))
	}
	vals := make([]interface{}, len(cols))
	for i := range cols {
		vals[i] = new(sql.RawBytes)
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(vals...)
		if err != nil {
			panic(fmt.Sprintf("提取数据失败:%s  %v ==> %v", query, args, err))
		}
		out = make(map[string]interface{})
		for i, col := range cols {
			out[col] = vals[i]

		}
	}

	return
}
