package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

func NewClient(serverIp string, serverPort int) *Client {
	// 创建客户端
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}
	// 连接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client.conn = conn
	return client
}

// 每个go文件中都有init函数
var serverIp string
var serverPort int

// 通过命令行设置ip和端口号
// ./client -h，可以设置的默认信息
func init() {
	flag.StringVar(&serverIp, "ip", "192.168.1.13", "设置服务器IP地址（默认192.168.1.13）")
	flag.IntVar(&serverPort, "port", 8888, "设置端口好（默认888）")
}
func main() {
	// 解析命令行
	flag.Parse()

	client := NewClient("192.168.1.13", 8888)
	if client == nil {
		fmt.Println(">>>>连接失败")
		return
	}
	fmt.Println(">>>>连接成功")
	// 启动客户端业务
	select {}
}
