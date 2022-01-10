package route

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"api.ikurum.cn/config"
	"api.ikurum.cn/util"
	"api.ikurum.cn/util/logs"
)

var (
	r_Mux URLHandlerContorller   // Mux 路由表
	mux   []URLHandlerContorller // 路由对象
)

// Router 路由接口，实现ServeHTTP
type Router struct{}

// URLHandlerContorller 路由控制器
type URLHandlerContorller struct {
	Func    func(http.ResponseWriter, *http.Request)
	Method  string
	Pattern string
}

// Listen 监听端口
func (r *Router) Listen(port string) {
	logs.Init()
	go util.StartToken()

	if b := strings.HasPrefix(port, ":"); !b {
		port = ":" + port
	}
	logs.Warning("init mysql ip: ", config.DB["ip"])
	logs.Info("启动服务: http://127.0.0.1", port)

	server := &http.Server{
		Addr:    port,
		Handler: r,
	}
	go server.ListenAndServe()

	listenSignal(context.Background(), server)
}

func listenSignal(ctx context.Context, httpSrv *http.Server) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-sigs
	logs.Warning("notify sigs")
	httpSrv.Shutdown(ctx)
	logs.Warning("http shutdown")
}

func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/favicon.ico" {
		logs.Info("URL Path:", req.URL.Path)
	}

	res.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	res.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	res.Header().Add("x-content-type-options", "nosniff")
	res.Header().Del("Content-Type")
	res.Header().Add("Content-Type", "application/json;utf-8")

	for _, URLHandlerContorller := range mux {
		if m, _ := regexp.MatchString(URLHandlerContorller.Pattern, req.URL.Path); m {
			if req.Method == URLHandlerContorller.Method {
				URLHandlerContorller.Func(res, req)
				return
			}
		}
	}

	res.Write([]byte("404"))
}

func GET(pattern string, f http.HandlerFunc) {
	r_Mux.GET(pattern, f)
}

func POST(pattern string, f http.HandlerFunc) {
	r_Mux.GET(pattern, f)
}

// GET 初始化路由
func (u URLHandlerContorller) GET(pattern string, f http.HandlerFunc) {
	println("GET 初始化路由: ", pattern)
	mux = append(mux, URLHandlerContorller{f, "GET", pattern})
}

// POST 初始化路由
func (u URLHandlerContorller) POST(pattern string, f http.HandlerFunc) {
	println("POST 初始化路由: ", pattern)
	mux = append(mux, URLHandlerContorller{f, "POST", pattern})
}
