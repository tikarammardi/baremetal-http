package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

var addr = "localhost:42069"

func main() {
	listener, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialUDP("udp", nil, listener)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		conn.Write([]byte(message))
	}
}
