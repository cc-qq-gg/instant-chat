package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
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

func (client *Client) UpdateName() bool {
	fmt.Println("请输入用户名：")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("con.Write error: ", err)
		return false
	}
	return true
}

// 处理server返回消息，显示到标准输出
func (client *Client) DealResponse() {
	// 一旦client.conn有数据，就直接copy到stout标准输出
	// 永久阻塞，等同于下面for循环
	io.Copy(os.Stdout, client.conn)
	// for {
	// 	buf := make([]byte, 1024)
	// 	client.conn.Read(buf)
	// 	fmt.Print(buf)
	// }
}

func (client *Client) PublicChat() {
	var chatMsg string
	// 提示用户输入信息
	fmt.Printf(">>>>请输入聊天内容, 输入exit退出\n")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		// 消息不为空时发送
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Printf("conn.write error: %v", err)
				break
			}
		}
		chatMsg = ""
		fmt.Printf(">>>>请输入聊天内容, 输入exit退出\n")
		fmt.Scanln(&chatMsg)
	}
}

func (client *Client) ShowOnlineUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Printf("ShowOnlineUsers conn.write error: %v", err)
		return
	}
}

func (client *Client) PrivateChat() {
	// 选择私聊的用户名
	var remoteName string
	var chatMsg string
	// 展示在线用户
	client.ShowOnlineUsers()
	fmt.Printf(">>>>请输入用户名, 输入exit退出\n")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Printf(">>>>请输入消息内容, 输入exit退出\n")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			// 消息不为空发送
			if len(remoteName) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("PriviteChat.conn.write.error", err)
					break
				}
			}
			chatMsg = ""
			fmt.Printf(">>>>请输入消息内容, 输入exit退出\n")
			fmt.Scanln(&chatMsg)
		}
		client.ShowOnlineUsers()
		fmt.Printf(">>>>请输入用户名, 输入exit退出\n")
		fmt.Scanln(&remoteName)
	}

}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}
		// 根据不同模式处理不同的业务
		switch client.flag {
		case 1:
			// 群聊模式
			client.PublicChat()
			break
		case 2:
			// 私聊模式
			client.PrivateChat()
			break
		case 3:
			client.UpdateName()
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

	// 单独开启一个goroutine处理server返回的消息
	go client.DealResponse()

	// 启动客户端业务
	client.Run()
}
