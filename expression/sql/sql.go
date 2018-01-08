package sql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"qlang.io/cl/qlang"
)

func Reg(defaultDB *sql.DB) {
	qlang.Import("sqlctx", map[string]interface{}{
		"QueryDB": QueryDB,
		"Query": func(ctx context.Context, query string, args ...interface{}) []map[string]interface{} {
			return QueryDB(defaultDB, ctx, query, args...)
		},
		"CountDB": QueryDBCount,
		"Count": func(ctx context.Context, query string, args ...interface{}) int {
			return QueryDBCount(defaultDB, ctx, query, args...)
		},
		"OneDB": QueryOneDB,
		"One": func(query string, ctx context.Context, args ...interface{}) map[string]interface{} {
			return QueryOneDB(defaultDB, ctx, query, args...)
		},
	})
	// qlang.Import("sql", map[string]interface{}{
	// 	"QueryDB": func(db *sql.DB, query string, args ...interface{}) []map[string]interface{} {
	// 		return QueryDB(db, context.Background(), query, args...)
	// 	},
	// 	"Query": func(query string, args ...interface{}) []map[string]interface{} {
	// 		return QueryDB(defaultDB, context.Background(), query, args...)
	// 	},
	// 	"CountDB": func(db *sql.DB, query string, args ...interface{}) int {
	// 		return QueryDBCount(db, context.Background(), query, args...)
	// 	},
	// 	"Count": func(query string, args ...interface{}) int {
	// 		return QueryDBCount(defaultDB, context.Background(), query, args...)
	// 	},
	// 	"OneDB": func(db *sql.DB, query string, args ...interface{}) map[string]interface{} {
	// 		return QueryOneDB(db, context.Background(), query, args...)
	// 	},
	// 	"One": func(query string, args ...interface{}) map[string]interface{} {
	// 		return QueryOneDB(defaultDB, context.Background(), query, args...)
	// 	},
	// })

}

// QueryDB 查询sql返回的所有行
func QueryDB(db *sql.DB, ctx context.Context, query string, args ...interface{}) (out []map[string]interface{}) {
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
	for i, _ := range cols {
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
			vmap[col] = vals[i]

		}
		out = append(out, vmap)
	}

	return
}

// QueryDBCount 查询sql匹配的条数
func QueryDBCount(db *sql.DB, ctx context.Context, query string, args ...interface{}) (count int) {

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
func QueryOneDB(db *sql.DB, ctx context.Context, query string, args ...interface{}) (out map[string]interface{}) {
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
	for i, _ := range cols {
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
