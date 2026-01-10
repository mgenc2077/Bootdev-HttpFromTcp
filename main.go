package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func() {
		defer close(lines)
		str := ""
		for {
			data := make([]byte, 8)
			n, err := f.Read(data)
			if err != nil && err != io.EOF {
				break
			}
			if strings.Contains(string(data[:n]), "\n") {
				firstNewline := strings.Index(string(data[:n]), "\n")
				str += string(data[:firstNewline])
				lines <- str
				remaining := data[firstNewline+1 : n]
				str = string(remaining)
			} else if err == io.EOF {
				if len(str)+n > 0 {
					lines <- str + string(data[:n])
				}
				break
			} else {
				str += string(data[:n])
			}
		}
	}()
	return lines
}

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

	// Read lines from the connection
	linesChan := getLinesChannel(conn)
	for line := range linesChan {
		fmt.Println(line)
	}

	// print on chanel closed
	if _, ok := <-linesChan; !ok {
		fmt.Println("Channel closed by client")
	}
}
