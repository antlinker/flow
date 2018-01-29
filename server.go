package flow

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/teambition/gear"
	"github.com/teambition/gear/logging"
	"github.com/teambition/gear/middleware/static"
)

type serverOptions struct {
	prefix      string
	staticRoot  string
	middlewares []gear.Middleware
}

// ServerOption 流程服务配置
type ServerOption func(*serverOptions)

// ServerPrefixOption 访问前缀
func ServerPrefixOption(prefix string) ServerOption {
	return func(opts *serverOptions) {
		opts.prefix = prefix
	}
}

// ServerStaticRootOption 静态文件目录
func ServerStaticRootOption(staticRoot string) ServerOption {
	return func(opts *serverOptions) {
		opts.staticRoot = staticRoot
	}
}

// ServerMiddlewareOption 中间件
func ServerMiddlewareOption(middlewares ...gear.Middleware) ServerOption {
	return func(opts *serverOptions) {
		opts.middlewares = middlewares
	}
}

// Server 流程管理服务
type Server struct {
	opts   serverOptions
	engine *Engine
	app    *gear.App
}

// Init 初始化
func (a *Server) Init(engine *Engine, opts ...ServerOption) *Server {
	a.engine = engine

	var o serverOptions
	for _, opt := range opts {
		opt(&o)
	}

	if o.prefix == "" {
		o.prefix = "/"
	}
	a.opts = o

	app := gear.New()

	app.UseHandler(logging.Default())

	for _, m := range a.opts.middlewares {
		app.Use(m)
	}

	if a.opts.staticRoot != "" {
		app.Use(newStaticMiddleware(a))
	}

	app.UseHandler(newRouterMiddleware(a))
	a.app = app

	return a
}

func (a *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.app.ServeHTTP(w, r)
}

// 静态文件中间件
func newStaticMiddleware(srv *Server) gear.Middleware {
	staticRoot, prefix := srv.opts.staticRoot, srv.opts.prefix
	staticMiddleware := static.New(static.Options{
		Root:        staticRoot,
		Prefix:      prefix,
		StripPrefix: true,
	})

	routerPrefix := regexp.MustCompile(`^(api)/.*`)
	return func(ctx *gear.Context) error {
		path := strings.TrimPrefix(ctx.Path, prefix)
		if routerPrefix.MatchString(path) {
			return nil
		}

		_, verr := os.Stat(filepath.Join(staticRoot, path))
		if verr != nil && os.IsNotExist(verr) {
			http.ServeFile(ctx.Res, ctx.Req, filepath.Join(staticRoot, "index.html"))
			return nil
		}

		return staticMiddleware(ctx)
	}
}

func newRouterMiddleware(srv *Server) gear.Handler {
	router := gear.NewRouter(gear.RouterOptions{
		Root:       fmt.Sprintf("%sapi", srv.opts.prefix),
		IgnoreCase: true,
	})

	api := new(API).Init(srv.engine)
	router.Get("/flow/page", api.QueryFlowPage)
	router.Get("/flow/:id", api.GetFlow)
	router.Delete("/flow/:id", api.DeleteFlow)
	router.Post("/flow", api.SaveFlow)

	return router
}
