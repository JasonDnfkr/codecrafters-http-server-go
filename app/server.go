package main

import (
	"bytes"
	"compress/gzip"
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

func createStatusLine(isSuccess bool) string {
	if isSuccess {
		return RESPONSE_OK
	} else {
		return RESPONSE_FAILED
	}
}

func addHeaders(headerType string, headerValue string, header *map[string]string) {
	if *header == nil {
		*header = make(map[string]string)
	}
	(*header)[headerType] = headerValue
}

func buildHeader(header map[string]string) string {
	var resp strings.Builder
	if len(header) == 0 {
		return ""
	}
	for key, value := range header {
		resp.WriteString(fmt.Sprintf("%s: %s%s", key, value, CRLF))
	}

	return resp.String()
}

func createHttpResponse(statusLine, header, body string) string {
	return fmt.Sprintf("%s%s%s%s%s", statusLine, CRLF, header, CRLF, body)
}

func getHeaders(request string) map[string]string {
	lines := strings.Split(request, "\r\n")

	headers := make(map[string]string)

	for i, line := range lines {
		if i == 0 {
			continue
		}
		pos := strings.Index(line, ":")
		if pos < 0 {
			break
		}
		header := line[0:pos]
		value := strings.TrimSpace(line[pos+1:])
		headers[header] = value
	}

	return headers
}

func getBody(request string) string {
	idx := strings.Index(request, CRLF+CRLF)
	idx += len(CRLF + CRLF)
	return request[idx:]
}

func ResponseHandler(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 512)
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		os.Exit(1)
	}

	request := string(bytes.Trim(buffer, "\x00"))

	requestHeaders := getHeaders(request)
	requestBody := getBody(request)

	var statusLine string
	var headers map[string]string
	var body string

	lines := strings.Split(request, CRLF)
	fmt.Println("\n======= HTTP REQUEST =======")
	for _, line := range lines {
		fmt.Println(line)
	}
	fmt.Println("\n======= HTTP REQUEST END =======")

	statusLine = createStatusLine(true)
	addHeaders("Content-Type", "text/plain", &headers)
	addHeaders("Content-Length", "0", &headers)

	lineSplit := strings.Split(lines[0], " ")
	method := lineSplit[0]
	path := lineSplit[1]

	switch method {
	case "GET":
		if path == "/" {
			statusLine = createStatusLine(true)
		} else if strings.Split(path, "/")[1] == "echo" {
			// get response body
			body = strings.Split(path, "/")[2]

			// found gzip
			gzipFound := false

			// get Accept-Encoding
			encodingLine := requestHeaders["Accept-Encoding"]
			encodings := strings.Split(encodingLine, ",")
			for _, encoding := range encodings {
				encoding = strings.TrimSpace(encoding)
				if encoding == "gzip" {
					gzipFound = true
				}
			}

			if gzipFound {
				statusLine = createStatusLine(true)
				addHeaders("Content-Encoding", "gzip", &headers)

				// gzip compress
				var buf bytes.Buffer
				writer := gzip.NewWriter(&buf)
				_, err := writer.Write([]byte(requestBody))
				if err != nil {
					fmt.Println("gzip error: " + err.Error())
					os.Exit(1)
				}
				err = writer.Close()
				if err != nil {
					fmt.Println("gzip error: " + err.Error())
					os.Exit(1)
				}

				fmt.Printf("----------------------===-%d\n", len(buf.String()))

				body = requestBody
				addHeaders("Content-Length", strconv.Itoa(len(buf.String())), &headers)
			} else {
				addHeaders("Content-Length", strconv.Itoa(len(body)), &headers)
			}
			addHeaders("Content-Type", "text/plain", &headers)
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
		} else if strings.Split(path, "/")[1] == "files" {
			fileName := strings.Split(path, "/")[2]
			dir := os.Args[2]
			data, err := os.ReadFile(dir + fileName)
			fmt.Printf("fileName: %s, dir: %s\n", fileName, dir)
			if err != nil {
				statusLine = createStatusLine(false)
				headers = make(map[string]string)
				body = ""
			} else {
				statusLine = createStatusLine(true)
				addHeaders("Content-Type", "application/octet-stream", &headers)
				addHeaders("Content-Length", strconv.Itoa(len(data)), &headers)
				body = string(data)
			}
		} else {
			statusLine = createStatusLine(false)
			headers = make(map[string]string)
			body = ""
		}

	case "POST":
		if strings.Split(path, "/")[1] == "files" {
			fileName := strings.Split(path, "/")[2]
			dir := os.Args[2]
			fmt.Printf("fileName: %s, dir: %s\n", fileName, dir)
			if err != nil {
				statusLine = createStatusLine(false)
				headers = make(map[string]string)
				body = ""
			} else {
				fmt.Println("\n======= BODY =======")
				fmt.Println(requestBody)
				fmt.Println("\n======= BODY END =======")

				err = os.WriteFile(dir+fileName, []byte(requestBody), 0644)
				if err != nil {
					fmt.Println("Error writing to file: ", err.Error())
					os.Exit(1)
				}

				statusLine = "HTTP/1.1 201 Created"
				headers = make(map[string]string)
				body = ""
			}
		}
	}

	resp := createHttpResponse(statusLine, buildHeader(headers), body)

	fmt.Println("======== RESPONSE ========")
	fmt.Println(resp)
	fmt.Println("======== RESPONSE END =======")

	conn.Write([]byte(resp))
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
