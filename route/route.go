package route

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
)

// Mux 路由表
var Mux URLHandlerContorller

// Router 路由接口，实现ServeHTTP
type Router struct{}

// URLHandlerContorller 路由控制器
type URLHandlerContorller struct {
	Func    func(http.ResponseWriter, *http.Request)
	Method  string
	Pattern string
}

// 路由对象
var mux []URLHandlerContorller

// Listen 监听端口
func (r *Router) Listen(port string) {
	if b := strings.HasPrefix(port, ":"); !b {
		port = ":" + port
	}
	fmt.Printf("启动服务: 127.0.0.1%s\n", port)

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
	fmt.Println("notify sigs")
	httpSrv.Shutdown(ctx)
	fmt.Println("http shutdown")
}

func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/favicon.ico" {
		fmt.Println("URL Path:", req.URL.Path)
	}

	res.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	res.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	res.Header().Add("x-content-type-options", "nosniff")
	res.Header().Del("Content-Type")
	res.Header().Add("Content-Type", "application/json;utf-8")

	for _, URLHandlerContorller := range mux {
		// fmt.Printf("url: %s, req: %s\n", URLHandlerContorller.Pattern, req.URL.Path)
		if m, _ := regexp.MatchString(URLHandlerContorller.Pattern, req.URL.Path); m {
			if req.Method == URLHandlerContorller.Method {
				URLHandlerContorller.Func(res, req)
				return
			}
		}
	}

	res.Write([]byte("404"))
}

// GET 初始化路由
func (u URLHandlerContorller) GET(pattern string, f http.HandlerFunc) {
	// fmt.Println("get pattern:", pattern)

	mux = append(mux, URLHandlerContorller{f, "GET", pattern})
}

// POST 初始化路由
func (u URLHandlerContorller) POST(pattern string, f http.HandlerFunc) {
	// fmt.Println("post pattern:", pattern)

	mux = append(mux, URLHandlerContorller{f, "POST", pattern})
}
