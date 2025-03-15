package main

import (
	"io"
	"log"
	"os"
	"path"
)

const (
	PORT = "4221"
)

func main() {
	router := NewRouter()
	router.HandlePath("/", indexEndpoint)
	router.HandlePattern("/echo/", echoEndpoint)
	router.HandlePath("/user-agent", userAgentEndpoint)
	router.HandlePattern("/files/", fileEndpoint)

	server := NewServer(router, "")
	if err := server.Start(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func indexEndpoint(w *ResponseWriter, req *HTTPRequest) {
	w.WriteStatus(200)
	w.WriteHeader("Content-Type", "text/plain")
}

func echoEndpoint(w *ResponseWriter, req *HTTPRequest) {
	if len(req.URI) < 7 {
		w.WriteStatus(400)
		return
	}
	w.WriteStatus(200)
	w.WriteHeader("Content-Type", "text/plain")
	w.Write([]byte(req.URI[6:]))
}

func userAgentEndpoint(w *ResponseWriter, req *HTTPRequest) {
	userAgent := req.Headers["User-Agent"]
	w.WriteStatus(200)
	w.WriteHeader("Content-Type", "text/plain")
	w.Write([]byte(userAgent))
}

func fileEndpoint(w *ResponseWriter, req *HTTPRequest) {
	if len(req.URI) < 7 {
		w.WriteStatus(400)
		return
	}
	var directory string
	for i, arg := range os.Args {
		if arg == "--directory" && i+1 < len(os.Args) {
			directory = os.Args[i+1]
			break
		}
	}
	if directory == "" {
		directory = "."
	}
	parsedFilePath := path.Join(directory, req.URI[6:])

	switch req.Method {
	case "GET":
		fileEndpointGET(w, parsedFilePath)
		return
	case "POST":
		fileEndpointPOST(w, req, parsedFilePath)
		return
	default:
		w.WriteStatus(405)
		w.WriteHeader("Allow", "GET, POST")
	}
}

func fileEndpointGET(w *ResponseWriter, parsedFilePath string) {
	if _, err := os.Stat(parsedFilePath); os.IsNotExist(err) {
		w.WriteStatus(404)
		return
	}
	file, err := os.Open(parsedFilePath)
	if err != nil {
		w.WriteStatus(500)
		log.Println("Failed to open file:", err)
		return
	}
	defer file.Close()

	body := make([]byte, 1024)
	for {
		n, err := file.Read(body)
		if err != nil && err != io.EOF {
			log.Println("Failed to read file:", err)
			w.WriteStatus(500)
			return
		}
		if n == 0 {
			break
		}
		w.Write(body[:n])
	}
	w.WriteStatus(200)
	w.WriteHeader("Content-Type", "application/octet-stream")
}

func fileEndpointPOST(w *ResponseWriter, req *HTTPRequest, parsedFilePath string) {
	file, err := os.Create(parsedFilePath)
	if err != nil {
		w.WriteStatus(500)
		log.Println("Failed to create file:", err)
		return
	}
	defer file.Close()

	if _, err := file.Write(req.Body); err != nil {
		log.Println("Failed to write file:", err)
		w.WriteStatus(500)
		return
	}
	w.WriteStatus(201)
}
