package main

import (
	"fmt"
	"net"
)

type HTTPResponse struct {
	Status       string
	StatusReason string
	Headers      map[string]string
	Body         []byte
}

func (r *HTTPResponse) AddHeader(key, value string) {
	if r.Headers == nil {
		r.Headers = make(map[string]string)
	}
	r.Headers[key] = value
}

func (r *HTTPResponse) Write(conn net.Conn) error {
	// Calculate Content-Length
	contentLength := len(r.Body)
	r.AddHeader("Content-Length", fmt.Sprintf("%d", contentLength))

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
