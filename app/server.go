package main

import (
	"fmt"
	"strconv"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

const CRLF = "\r\n"

func Response(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		return
	}

	bufferString := string(buffer)
	fmt.Println(bufferString)

	path := strings.Split(bufferString, " ")[1]
	fmt.Println("[Path]" + path)
	if path == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else if strings.Split(path, "/")[1] == "echo" {
		str := strings.Split(path, "/")[2]
		var resp strings.Builder
		resp.WriteString("HTTP/1.1 200 OK")
		resp.WriteString(CRLF)
		resp.WriteString("Content-Type: text/plain")
		resp.WriteString(CRLF)
		resp.WriteString("Content-Length: " + strconv.Itoa(len(str)))
		resp.WriteString(CRLF + CRLF)
		resp.WriteString(str)

		conn.Write([]byte(resp.String()))
	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
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

	Response(conn)
}
