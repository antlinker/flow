package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/antlinker/flow"
	"github.com/antlinker/flow/service/db"

	"github.com/teambition/gear"

	_ "github.com/go-sql-driver/mysql"
)

var (
	staticRoot string
	dsn        string
	addr       string
)

func init() {
	flag.StringVar(&staticRoot, "static", "", "静态目录")
	flag.StringVar(&dsn, "dsn", "root:123456@tcp(127.0.0.1:3306)/flows?charset=utf8", "MySQL链接路径")
	flag.StringVar(&addr, "addr", ":6062", "监听地址")
}

func main() {
	flag.Parse()

	flow.Init(
		db.SetDSN(dsn),
		db.SetTrace(true),
	)

	serverOptions := []flow.ServerOption{
		flow.ServerStaticRootOption(staticRoot),
		flow.ServerPrefixOption("/flow/"),
		flow.ServerMiddlewareOption(filter),
	}

	http.Handle("/flow/", flow.StartServer(serverOptions...))

	log.Printf("服务运行在[%s]端口...", addr)
	log.Fatalf("%v", http.ListenAndServe(addr, nil))
}

func filter(ctx *gear.Context) error {
	fmt.Printf("请求参数：%s - %s \n", ctx.Path, ctx.Method)
	return nil
}
