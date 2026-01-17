package main

import (
	"fmt"
	"io"
	"mgenc2077/httpfromtcp/internal/request"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	fmt.Println("Server is listening on port 42069")

	conn, err := listener.Accept()
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fmt.Println("Client connected")

	req, err := request.RequestFromReader(conn)
	if err != nil && err != io.EOF {
		fmt.Println("Error reading request:", err)
		return
	}

	fmt.Println("Request line:")
	fmt.Printf("- Method: %s\n", req.RequestLine.Method)
	fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
	fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)
	fmt.Println("Headers:")
	for key, value := range req.Headers {
		fmt.Printf("- %s: %s\n", key, value)
	}
	fmt.Println("Body:")
	fmt.Println(string(req.Body))
}
