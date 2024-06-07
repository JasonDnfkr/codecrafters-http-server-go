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
const RESPONSE_OK = "HTTP/1.1 200 OK"
const RESPONSE_FAILED = "HTTP/1.1 404 Not Found"

func ResponseHandler(conn net.Conn) {
	defer conn.Close()

	createStatusLine := func(isSuccess bool) string {
		if isSuccess {
			return RESPONSE_OK
		} else {
			return RESPONSE_FAILED
		}
	}

	addHeaders := func(headerType string, headerValue string, header *map[string]string) {
		if *header == nil {
			*header = make(map[string]string)
		}
		(*header)[headerType] = headerValue
	}

	buildHeader := func(header map[string]string) string {
		var resp strings.Builder
		if len(header) == 0 {
			return ""
		}
		for key, value := range header {
			resp.WriteString(fmt.Sprintf("%s: %s%s", key, value, CRLF))
		}

		return resp.String()
	}

	createHttpResponse := func(statusLine, header, body string) string {
		return fmt.Sprintf("%s%s%s%s%s", statusLine, CRLF, header, CRLF, body)
	}

	for {
		buffer := make([]byte, 1024)
		_, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from connection: ", err.Error())
			os.Exit(1)
		}

		request := string(buffer)

		var statusLine string
		var headers map[string]string
		var body string

		lines := strings.Split(request, CRLF)
		fmt.Println("======= HTTP REQUEST =======")
		for _, line := range lines {
			fmt.Println(line)
		}
		fmt.Println("======= HTTP REQUEST END =======")

		statusLine = createStatusLine(true)
		addHeaders("Content-Type", "text/plain", &headers)
		addHeaders("Content-Length", "0", &headers)

		path := strings.Split(lines[0], " ")[1]
		if path == "/" {
			statusLine = createStatusLine(true)
		} else if strings.Split(path, "/")[1] == "echo" {
			body = strings.Split(path, "/")[2]
			addHeaders("Content-Length", strconv.Itoa(len(body)), &headers)
		} else if strings.Split(path, "/")[1] == "user-agent" {
			fmt.Println("get user agent header")
			for _, line := range lines {
				if strings.HasPrefix(line, "User-Agent:") {
					// get foobar/1.2.3 ...
					content := strings.Split(line, " ")[1]
					addHeaders("Content-Length", strconv.Itoa(len(content)), &headers)
					body = content
				}
			}
		} else {
			statusLine = createStatusLine(false)
			headers = make(map[string]string)
			body = ""
		}

		resp := createHttpResponse(statusLine, buildHeader(headers), body)

		fmt.Println("======== RESPONSE ========")
		fmt.Println(resp)
		fmt.Println("======== RESPONSE END =======")

		conn.Write([]byte(resp))
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

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go ResponseHandler(conn)
	}
}
