package main

import (
	"fmt"
	"os"
)

func main() {
	msgs, err := os.Open("messages.txt")
	if err != nil {
		panic(err)
	}
	defer msgs.Close()
	for {
		data := make([]byte, 8)
		n, err := msgs.Read(data)
		if err != nil {
			break
		}
		fmt.Printf("read: %s\n", string(data[:n]))
	}
}
