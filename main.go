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

	var b bytes.Buffer
	b.WriteString("---\n")
	b.WriteString(fmt.Sprintf("Client: %s:%s\n", clientIP, clientPort))
	b.WriteString(fmt.Sprintf("Method: %s\n", r.Method))
	b.WriteString(fmt.Sprintf("URL: %s\n", r.URL))
	b.WriteString("Headers:\n")
	for k, v := range r.Header {
		b.WriteString(fmt.Sprintf("    %s: %s\n", k, v))
	}
	b.WriteString(fmt.Sprintf("Body:\n%s\n", bodyBytes))

	mutex.Lock()
	fmt.Print(b.String())
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
