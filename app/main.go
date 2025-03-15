package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

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

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	// Read HTTP request
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading request: ", err.Error())
		return
	}
	req := ParseHTTPRequest(string(buf[:n]))
	fmt.Printf("Got request: %v\n", req)

	var resp *HTTPResponse
	if req.URI == "/" {
		resp = &HTTPResponse{
			Status:       "200",
			StatusReason: "OK",
		}
	} else if len(req.URI) > 6 && strings.HasPrefix(req.URI, "/echo/") {
		resp = &HTTPResponse{
			Status:       "200",
			StatusReason: "OK",
			Body:         []byte(req.URI[6:]),
		}
		resp.AddHeader("Content-Type", "text/plain")
	} else {
		resp = &HTTPResponse{
			Status:       "404",
			StatusReason: "Not Found",
		}
	}
	resp.Write(conn)
}
