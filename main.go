package main

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof" // 导入 pprof 包，注册 /debug/pprof 端点
	"os"
)

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/sum", sumHandler)

	mux := http.NewServeMux()
	mux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))

	mux.HandleFunc("/whoami", func(w http.ResponseWriter, _ *http.Request) {
		h, _ := os.Hostname()
		fmt.Println(w, h)
	})

	//registerPprof(mux)

	//fmt.Println("Server started at :8080")
	fmt.Println("Server started at :8176")
	fmt.Println("GOCOVERDIR:", os.Getenv("GOCOVERDIR")) // 打印 GOCOVERDIR 环境变量的值
	//http.ListenAndServe(":8080", nil)
	http.ListenAndServe(":8176", nil)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

// http.HandleFunc("/sum", sumHandler)
func sumHandler(w http.ResponseWriter, r *http.Request) {
	a, b := 0, 0
	fmt.Sscanf(r.URL.Query().Get("a"), "%d", &a)
	fmt.Sscanf(r.URL.Query().Get("b"), "%d", &b)
	res := a + b
	w.Write([]byte(fmt.Sprintf("sum=%d", res)))
}

// 把所有 pprof 相关 handler 显式挂到自定义 mux 上
func registerPprof(mux *http.ServeMux) {
	// index 和四个固定函数
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// 其余的 profile（heap、goroutine、mutex、coverage…）用 Handler(name)
	for _, p := range []string{
		"allocs", "block", "goroutine", "heap",
		"mutex", "threadcreate", "coverage", // coverage 需 Go 1.22+
	} {
		mux.Handle("/debug/pprof/"+p, pprof.Handler(p))
	}
}

//$env:GOCOVERDIR="coverdata"
//go run -cover -coverpkg=./... .
