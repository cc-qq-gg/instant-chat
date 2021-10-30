package main

// 两个包都属于main包，没必要import
func main() {
	server := NewServer("127.0.0.1", 8888)
	server.Start()
}
