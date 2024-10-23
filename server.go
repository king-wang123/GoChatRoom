package main

import (
	"fmt"
	"net"
	"strconv"
	"sync"
)

type Server struct {
	Ip        string
	Port      int
	OnlineMap map[string]*User // online user map
	mapLock   sync.RWMutex
	Message   chan string // broadcast message
}

func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
}

func (this *Server) ListenBroadcast() {
	for {
		msg := <-this.Message

		// send message to all online users
		this.mapLock.RLock()
		for _, user := range this.OnlineMap {
			user.C <- msg
		}
		this.mapLock.RUnlock()
	}
}

func (this *Server) Broadcast(user *User, message string) {
	sendMsg := "[" + user.Name + "]: " + message
	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	//fmt.Println("Successfully established connection")
	user := NewUser(conn)

	// add user to online map
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()
	// broadcast message
	this.Broadcast(user, "has joined the chat")

	// block the handler
	select {}
}

// Start server
func (this *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", this.Ip+":"+strconv.Itoa(this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}

	//close socket
	defer listener.Close()

	go this.ListenBroadcast()

	// accept
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err:", err)
			continue
		}

		//do handler
		go this.Handler(conn)
	}

}
