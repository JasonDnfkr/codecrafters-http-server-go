package main

import (
	"fmt"
	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func ExtractUrlPath(conn net.Conn) {
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		os.Exit(1)
	}

	bufferStr := string(buffer)

	for i := 0; i < len(bufferStr); i++ {
		if bufferStr[i] == '/' {
			idx := i + 1
			for idx < len(bufferStr) && bufferStr[idx] == ' ' {
				idx++
			}
			if bufferStr[idx:idx+4] == "HTTP" {
				_, err := conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
				fmt.Printf("Response 200 Completed")
				if err != nil {
					fmt.Println("Error writing to connection: ", err.Error())
					os.Exit(1)
				}
				break
			} else {
				_, err := conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
				if err != nil {
					fmt.Println("Error writing to connection: ", err.Error())
					os.Exit(1)
				}
			}
		}
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	ExtractUrlPath(conn)
}
