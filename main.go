package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	runtimepprof "runtime/pprof"
	//_ "net/http/pprof" // 导入 pprof 包，注册 /debug/pprof 端点
	"os"
)

func init() {
	fmt.Println("GOCOVERDIR =", os.Getenv("GOCOVERDIR"))
	if runtimepprof.Lookup("coverage") != nil {
		fmt.Println("coverage enabled = true")
	} else {
		fmt.Println("coverage enabled = false")
	}
}

func main() {
	mux := http.NewServeMux()

	// 业务
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/sum", sumHandler)

	// who-am-I
	mux.HandleFunc("/whoami", func(w http.ResponseWriter, _ *http.Request) {
		h, _ := os.Hostname()
		fmt.Fprintln(w, h) // 写到 HTTP 响应里
	})

	// pprof
	registerPprof(mux)

	fmt.Println("Server started at :8176")
	fmt.Println("GOCOVERDIR:", os.Getenv("GOCOVERDIR")) // 打印 GOCOVERDIR 环境变量的值

	log.Println("listen :8176")
	log.Fatal(http.ListenAndServe(":8176", mux)) // ← 传入 mux  // nil = DefaultServeMux
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

	// 其它 profile（一定要把 coverage 带上）
	for _, p := range []string{
		"allocs", "block", "goroutine", "heap",
		"mutex", "threadcreate", "coverage",
	} {
		mux.Handle("/debug/pprof/"+p, pprof.Handler(p))
	}
}

//$env:GOCOVERDIR="coverdata"
//go run -cover -coverpkg=./... .
