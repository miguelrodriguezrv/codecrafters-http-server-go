package main

import (
	"fmt"
)

type ResponseWriter struct {
	Response *HTTPResponse
}

func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{
		Response: &HTTPResponse{
			Status:       "200",
			StatusReason: "OK",
			Headers:      make(map[string]string),
			Body:         []byte{},
		},
	}
}

func (w *ResponseWriter) Write(p []byte) (n int, err error) {
	w.Response.Body = append(w.Response.Body, p...)
	return len(p), nil
}

func (w *ResponseWriter) WriteStatus(statusCode int) {
	w.Response.Status = fmt.Sprintf("%d", statusCode)
	w.Response.StatusReason = StatusText(statusCode)
}

func (w *ResponseWriter) WriteHeader(key, value string) {
	if w.Response.Headers == nil {
		w.Response.Headers = make(map[string]string)
	}
	w.Response.Headers[key] = value
}

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

func (r *HTTPResponse) ToBytes() []byte {
	// Calculate Content-Length
	contentLength := len(r.Body)
	r.AddHeader("Content-Length", fmt.Sprintf("%d", contentLength))

	// Convert headers to byte array
	headersBytes := make([]byte, 0)
	for key, value := range r.Headers {
		headersBytes = append(headersBytes, fmt.Appendf(nil, "%s: %s\r\n", key, value)...)
	}
	headersBytes = append(headersBytes, []byte("\r\n")...)

	// Write the status line
	statusLineBytes := fmt.Appendf(nil, "HTTP/1.1 %s %s\r\n", r.Status, r.StatusReason)

	// Write the headers
	headersBytes = append(statusLineBytes, headersBytes...)

	// Write the body
	bodyBytes := append(headersBytes, r.Body...)

	return bodyBytes
}
