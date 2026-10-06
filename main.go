package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
const port = ":42069"

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("error listening for tcp traffic: %s\n", err.Error())
	}
	defer listener.Close()

	fmt.Println("Listening for traffic on ", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}

		fmt.Println("Accepted connection from", conn.RemoteAddr())

		linesChan := getLinesChannel(conn)

		for line := range linesChan {
			fmt.Println(line)
		}
		fmt.Println("Closing connection from", conn.RemoteAddr())
	}
}

func getLinesChannel(file io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func() {
		defer file.Close()
		defer close(lines)
		currentLinesContents := ""

		for {
			buffer := make([]byte, 8, 8)
			n, err := file.Read(buffer)
			if err != nil {
				if currentLinesContents != "" {
					lines <- currentLinesContents
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Println("Error reading from file:", err.Error())
				return
			}
			str := string(buffer[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				lines <- fmt.Sprintf("%s%s", currentLinesContents, parts[i])
				currentLinesContents = ""
			}
			currentLinesContents += parts[len(parts)-1]
		}
	}()
	return lines
}
