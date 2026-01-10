package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	msgs, err := os.Open("messages.txt")
	if err != nil {
		panic(err)
	}
	defer msgs.Close()
	str := ""
	for {
		data := make([]byte, 8)
		n, err := msgs.Read(data)
		if err != nil && err != io.EOF {
			break
		}
		if strings.Contains(string(data[:n]), "\n") {
			firstNewline := strings.Index(string(data[:n]), "\n")
			str += string(data[:firstNewline])
			fmt.Printf("read: %s\n", str)
			remaining := data[firstNewline+1 : n]
			str = string(remaining)
		} else if err == io.EOF {
			fmt.Printf("read: %s\n", str+string(data[:n]))
			break
		} else {
			str += string(data[:n])
		}
	}
}
