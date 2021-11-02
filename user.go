package main

import "net"

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

// 用户广播消息业务
func (this *User) DoMessage(msg string) {
	// 查询在线用户，并返回给当前用户
	if msg == "who" {
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			this.C <- "[" + user.Addr + "] " + user.Name + "online...\n"
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		this.Rename(msg)
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
