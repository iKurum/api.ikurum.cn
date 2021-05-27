package route

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
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

	err := http.ListenAndServe(port, r)
	if err != nil {
		log.Fatal("ListenAndServe err:", err)
	}
}

func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/favicon.ico" {
		fmt.Println("URL Path:", req.URL.Path)
	}

	res.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	res.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	res.Header().Set("content-type", "application/json")             //返回数据格式是json

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
	fmt.Println("mux pattern:", pattern)

	mux = append(mux, URLHandlerContorller{f, "GET", pattern})
}

// POST 初始化路由
func (u URLHandlerContorller) POST(pattern string, f http.HandlerFunc) {
	fmt.Println("mux pattern:", pattern)

	mux = append(mux, URLHandlerContorller{f, "POST", pattern})
}
