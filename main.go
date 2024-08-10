package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
)

var mutex sync.Mutex

func logRequest(w http.ResponseWriter, r *http.Request) {
	// 从请求体中读取数据
	bodyBytes, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// 获取客户端 IP 和端口
	clientIP, clientPort, _ := net.SplitHostPort(r.RemoteAddr)

	mutex.Lock()
	// 打印请求信息
	fmt.Printf("---\n")
	fmt.Printf("Client: %s:%s\n", clientIP, clientPort)
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL)
	fmt.Printf("Headers:\n")
	for k, v := range r.Header {
		fmt.Printf("    %s: %v\n", k, v)
	}
	fmt.Printf("Body:\n%s\n", bodyBytes)
	mutex.Unlock()

	// 返回 200 OK
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/", logRequest)
	fmt.Println("Starting server on: 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
