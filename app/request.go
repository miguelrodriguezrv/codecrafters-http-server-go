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
func ParseHTTPRequest(request string) *HTTPRequest {
	// Split the request into headers and body
	parts := strings.SplitN(request, "\r\n\r\n", 2)
	headerPart := parts[0]

	// Initialize the request struct
	parsedRequest := &HTTPRequest{
		Headers: make(map[string]string),
		Body:    nil,
	}

	// Parse the body if it exists
	if len(parts) > 1 {
		parsedRequest.Body = []byte(parts[1])
	}

	// Split the header part into lines
	lines := strings.Split(headerPart, "\r\n")

	// Parse the request line (first line)
	if len(lines) > 0 {
		requestLineParts := strings.Split(lines[0], " ")
		if len(requestLineParts) >= 3 {
			parsedRequest.Method = requestLineParts[0]
			parsedRequest.URI = requestLineParts[1]
			parsedRequest.Protocol = requestLineParts[2]
		}
	}

	// Parse the headers (remaining lines)
	for i := 1; i < len(lines); i++ {
		if lines[i] == "" {
			continue
		}
		headerParts := strings.SplitN(lines[i], ": ", 2)
		if len(headerParts) == 2 {
			parsedRequest.Headers[headerParts[0]] = headerParts[1]
		}
	}

	return parsedRequest
}
