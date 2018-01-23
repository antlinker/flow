package main

import (
	"flag"
	"net/http"

	"github.com/LyricTian/logger"

	"gitee.com/antlinker/flow"
	"gitee.com/antlinker/flow/service/db"
	_ "github.com/go-sql-driver/mysql"
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
	}

	http.Handle("/flow/", flow.StartServer(serverOptions...))

	logger.Infof("服务运行在[6062]端口...")
	logger.Fatalf("%v", http.ListenAndServe(":6062", nil))
}
