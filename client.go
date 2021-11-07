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
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	// 创建客户端
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       9999,
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

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1.群聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.修改用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>>请输入合法范围内的数字<<<<")
		return false
	}
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {

		}
		// 根据不同模式处理不同的业务
		switch client.flag {
		case 1:
			fmt.Println("群聊模式")
			break
		case 2:
			fmt.Println("私聊模式")
			break
		case 3:
			fmt.Println("修改用户名")
			break
		case 0:
			fmt.Println("退出")
			break
		}
	}
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
	client.Run()
}
