package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	// 启动监听当前user channel消息的goroutine
	go user.ListenMessage()
	return user
}

// 用户上线
func (this *User) Online() {
	// 添加到onlineMap中
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()
	this.server.BoradCast(this, "上线")
}

// 用户下线
func (this *User) Offline() {
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()
	this.server.BoradCast(this, "下线")
}
func (this *User) Rename(newName string) {
	// 检查是否被暂用
	newName = newName[7:]
	_, ok := this.server.OnlineMap[newName]
	if ok {
		this.SendMsg("【" + newName + "】" + "已被占用" + "\n")
	} else {
		this.server.mapLock.Lock()
		delete(this.server.OnlineMap, this.Name)
		this.server.OnlineMap[newName] = this
		this.server.mapLock.Unlock()
		this.Name = newName
		this.SendMsg("修改成功，当前用户名：" + newName + "\n")
	}
}
func (this *User) PriviteChat(msg string) {
	// 获取to用户名
	toUserName := strings.Split(msg, "|")[1]
	if toUserName == "" {
		this.SendMsg("消息格式不正确，请使用\"to|张三|消息\"的格式")
		return
	}
	toUser, ok := this.server.OnlineMap[toUserName]
	if !ok {
		this.SendMsg("该用户不存" + "\n")
		return
	}
	toMsg := strings.Split(msg, "|")[2]
	if toMsg == "" {
		this.SendMsg("请输入消息内容" + "\n")
		return
	}
	toUser.SendMsg(toMsg)
}
func (this *User) GetOnlineUsers() {
	this.server.mapLock.Lock()
	for _, user := range this.server.OnlineMap {
		this.C <- "[" + user.Addr + "] " + user.Name + ":online...\n"
	}
	this.server.mapLock.Unlock()
}

// 用户广播消息业务
func (this *User) DoMessage(msg string) {
	// 查询在线用户，并返回给当前用户
	if msg == "who" {
		this.GetOnlineUsers()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		this.Rename(msg)
	} else if len(msg) > 3 && msg[:3] == "to|" {
		this.PriviteChat(msg)
	} else {
		this.server.BoradCast(this, msg)
	}
}

// 给当前用户的客户端发送消息
func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

// 监听当前User channel的方法，有消息就发给客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.SendMsg(msg + "\n")
	}
}
