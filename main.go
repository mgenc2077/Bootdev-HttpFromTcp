package main

import (
	"fmt"
	"io"
	"os"
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
	msgs, err := os.Open("messages.txt")
	if err != nil {
		panic(err)
	}
	defer msgs.Close()
	linesChan := getLinesChannel(msgs)
	for line := range linesChan {
		fmt.Printf("read: %s\n", line)
	}
}
