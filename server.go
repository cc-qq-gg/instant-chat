package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
	// 在线用户列表
	OnlineMap map[string]*User
	// OnlineMap是全局的，需要加锁
	mapLock sync.Mutex
	// 消息广播的channel
	MessageChannel chan string
}

// N大写表示，该方法对外开放
// 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:             ip,
		Port:           port,
		OnlineMap:      make(map[string]*User),
		MessageChannel: make(chan string),
	}
	return server
}

// 广播消息的方法
func (this *Server) BoradCast(user *User, msg string) {
	sendMsg := "[" + user.Name + "]" + msg
	this.MessageChannel <- sendMsg
}

// 监听MessageChannel，广播给全部在线user
func (this *Server) ListenMessage() {
	for {
		msg := <-this.MessageChannel
		// 将msg发送给全部在线用户
		this.mapLock.Lock()
		for _, client := range this.OnlineMap {
			client.C <- msg
		}
		this.mapLock.Unlock()
	}
}
func (this *Server) Handler(conn net.Conn) {
	// ...当前连接的任务
	fmt.Println("连接建立成功")
	// 用户上线，加入到OnlineMap中
	user := NewUser(conn)
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()

	// 广播当前用户上线消息
	this.BoradCast(user, "已上线")

	// 阻塞当前handler，否则当前goroutine会dead，里面的子goroutine也dead
	select {}

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

	// 启动监听MessageChannel的gotoutine
	go this.ListenMessage()
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
