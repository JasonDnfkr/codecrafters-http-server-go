package main

import (
	"fmt"
	"strings"
)

func gHeaders(request string) map[string]string {
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

func main() {
	request := "GET / HTTP/1.1\r\n" + "GET /echo/foo HTTP/1.1\nHost: localhost:4221\nUser-Agent: curl/7.64.1\nAccept-Encoding: gzip\r\n\r\nbody"
	headers := gHeaders(request)
	for k, v := range headers {
		fmt.Printf("%s: %s\n", k, v)
	}
}
