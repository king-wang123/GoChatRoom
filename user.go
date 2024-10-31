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
	go user.ListenMessage()
	return user
}

func (this *User) Online() {
	// add user to online map
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()
	// broadcast message
	this.server.Broadcast(this, "has joined the chat")
}

func (this *User) Offline() {
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	this.server.Broadcast(this, "has left the chat")
	this.conn.Close()
}

func (this *User) Send(msg string) {
	this.conn.Write([]byte(msg))
}

func (this *User) DoMessage(msg string) {
	if msg == "$who" {
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Name + "] " + "is online\n"
			this.Send(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 8 && msg[:8] == "$rename|" {
		newName := msg[8:]
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.Send("This name is already taken, please choose another one.\n")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.Name = newName
			this.server.OnlineMap[this.Name] = this
			this.server.mapLock.Unlock()
			this.Send("Your name has been changed to " + newName + "\n")
		}
	} else {
		this.server.Broadcast(this, msg)
	}
}

func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
