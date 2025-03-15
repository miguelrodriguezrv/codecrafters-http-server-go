package main

import (
	"fmt"
	"net"
	"os"
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

type HTTPResponse struct {
	Status       string
	StatusReason string
	Headers      map[string]string
	Body         []byte
}

func (r HTTPResponse) Write(conn net.Conn) error {
	// Write the status line
	_, err := conn.Write(fmt.Appendf(nil, "HTTP/1.1 %s %s\r\n", r.Status, r.StatusReason))
	if err != nil {
		return err
	}

	// Write the headers
	for key, value := range r.Headers {
		_, err := conn.Write(fmt.Appendf(nil, "%s: %s\r\n", key, value))
		if err != nil {
			return err
		}
	}
	// Mark end of headers
	conn.Write([]byte("\r\n"))

	// Write the body
	_, err = conn.Write(r.Body)
	if err != nil {
		return err
	}

	return nil
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	resp := HTTPResponse{
		Status:       "200",
		StatusReason: "OK",
	}
	resp.Write(conn)
}
