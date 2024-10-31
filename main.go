package main

// go build -o server.exe .\main.go .\server.go .\user.go
func main() {
	server := NewServer("127.0.0.1", 8888)
	server.Start()
}
