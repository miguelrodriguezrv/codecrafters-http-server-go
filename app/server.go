package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type Server struct {
	router   *Router
	basePath string
}

func NewServer(router *Router, basePath string) *Server {
	if basePath == "" {
		basePath = "./"
	}
	return &Server{
		router:   router,
		basePath: basePath,
	}
}

func (s *Server) Start() error {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", PORT))
	if err != nil {
		return err
	}
	defer l.Close()

	log.Printf("Loaded %d handlers", s.router.NumRoutes())

	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	// Read HTTP request
	buf := make([]byte, 1024)
	keepConnOpen := true
	for keepConnOpen {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading request: ", err.Error())
			return
		}
		req := ParseHTTPRequest(string(buf[:n]))

		resp := NewResponseWriter()
		handler := s.router.Route(req)
		handler(resp, req)

		resp.Response.Body = handleCompression(req.Headers["Accept-Encoding"], resp)
		if connHeader, ok := req.Headers["Connection"]; ok && connHeader == "close" {
			resp.WriteHeader("Connection", "close")
			keepConnOpen = false
		}
		if _, err = conn.Write(resp.Response.ToBytes()); err != nil {
			fmt.Println("Error writing response: ", err.Error())
		}
	}
}

func handleCompression(encoding string, resp *ResponseWriter) []byte {
	for encoding := range strings.SplitSeq(encoding, ", ") {
		switch encoding {
		case "gzip":
			resp.WriteHeader("Content-Encoding", "gzip")
			return gzipCompress(resp.Response.Body)
		case "deflate":
			resp.WriteHeader("Content-Encoding", "deflate")
			return deflate(resp.Response.Body)
		}
	}
	return resp.Response.Body
}
