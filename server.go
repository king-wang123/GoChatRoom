package main

import (
	"fmt"
	"net"
	"strconv"
)

type Server struct {
	Ip   string
	Port int
}

func NewServer(ip string, port int) *Server {
	return &Server{Ip: ip, Port: port}
}

func (this *Server) Handler(conn net.Conn) {
	fmt.Println("Successfully established connection")
}

func (this *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", this.Ip+":"+strconv.Itoa(this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}

	//close socket
	defer listener.Close()

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
