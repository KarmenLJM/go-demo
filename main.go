package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof" // 导入 pprof 包，注册 /debug/pprof 端点
	"os"
)

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/sum", sumHandler)

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

//$env:GOCOVERDIR="coverdata"
//go run -cover -coverpkg=./... .
