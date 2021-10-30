package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
}

// N大写表示，该方法对外开放
// 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

func (this *Server) Handler(conn net.Conn) {
	// ...当前连接的任务
	fmt.Println("连接建立成功")
}

// 启动server
func (this *Server) Start() {
	// socket listen
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen error: ", err)
	}
	// 防止遗忘关闭
	// close listen socket
	defer listen.Close()
	// 循环处理请求
	for {
		// 阻塞等待
		// accept
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("net.Accept error", err)
			continue
		}

		// 如果成功，表明有一个链接进来
		// 为了不耽误下次Accept，开启一个异步的go程
		// do handler
		go this.Handler(conn)
	}
}
