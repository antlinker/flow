package main

import (
	"flag"
	"fmt"
	"net/http"

	"gitee.com/antlinker/flow"
	"gitee.com/antlinker/flow/service/db"
	"github.com/LyricTian/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/teambition/gear"
)

var (
	staticRoot string
)

func init() {
	flag.StringVar(&staticRoot, "static", "", "静态目录")
}

func main() {
	flag.Parse()

	flow.Init(&db.Config{
		DSN:   "root:123456@tcp(127.0.0.1:3306)/flows?charset=utf8",
		Trace: true,
	})

	serverOptions := []flow.ServerOption{
		flow.ServerStaticRootOption(staticRoot),
		flow.ServerPrefixOption("/flow/"),
		flow.ServerMiddlewareOption(filter),
	}

	http.Handle("/flow/", flow.StartServer(serverOptions...))

	logger.Infof("服务运行在[6062]端口...")
	logger.Fatalf("%v", http.ListenAndServe(":6062", nil))
}

func filter(ctx *gear.Context) error {
	fmt.Printf("请求参数：%s - %s \n", ctx.Path, ctx.Method)
	return nil
}
