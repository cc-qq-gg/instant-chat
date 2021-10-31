package main

// 两个包都属于main包，没必要import
func main() {
	server := NewServer("192.168.1.13", 8888)
	server.Start()
}
