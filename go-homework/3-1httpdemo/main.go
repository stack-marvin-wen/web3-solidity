package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "欢迎访问路径: %s\n", r.URL.Path)
		fmt.Fprintf(w, "欢迎访问方法: %s\n", r.Method)
	})
	http.HandleFunc("/hello", helloHandleFunc)
	fmt.Println("服务器启动, 监听端口8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("服务器启动失败:", err)
	}
}
func helloHandleFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "欢迎访问路径: %s\n", r.URL.Path)
	fmt.Fprintf(w, "欢迎访问方法: %s\n", r.Method)
}
