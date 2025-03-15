package main

import "strings"

type HTTPRequest struct {
	Method   string
	URI      string
	Protocol string
	Headers  map[string]string
	Body     []byte
}

// Parse the request string into an HTTPRequest struct
func ParseHTTPRequest(request string) HTTPRequest {
	// Split the request string into parts
	requestParts := strings.Split(request, "\r\n")
	// Parse the request line
	parts := strings.Split(requestParts[0], " ")
	parsedRequest := HTTPRequest{
		Method:   parts[0],
		URI:      parts[1],
		Protocol: parts[2],
		Headers:  make(map[string]string),
		Body:     nil,
	}
	// Parse the headers
	for _, header := range requestParts[1:] {
		if header == "" {
			break
		}
		headerParts := strings.Split(header, ": ")
		parsedRequest.Headers[headerParts[0]] = headerParts[1]
	}
	// Parse the body
	if len(requestParts) > 2 {
		parsedRequest.Body = []byte(requestParts[len(requestParts)-1])
	}

	return parsedRequest
}
