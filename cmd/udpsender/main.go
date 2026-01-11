package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	udpaddr, _ := net.ResolveUDPAddr("udp", "localhost:42069")
	conn, _ := net.DialUDP("udp", nil, udpaddr)
	defer conn.Close()

	msg := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		line, err := msg.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
		}
		conn.Write([]byte(line))
	}
}
