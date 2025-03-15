package main

import (
	"strings"
)

type EndpointHandler func(w *ResponseWriter, r *HTTPRequest)

type Router struct {
	routes        map[string]EndpointHandler
	patternRoutes []struct {
		pattern string
		handler EndpointHandler
	}
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]EndpointHandler),
	}
}

func (r *Router) NumRoutes() int {
	return len(r.routes) + len(r.patternRoutes)
}

func (r *Router) HandlePath(path string, handler EndpointHandler) {
	r.routes[path] = handler
}

func (r *Router) HandlePattern(pattern string, handler EndpointHandler) {
	r.patternRoutes = append(r.patternRoutes, struct {
		pattern string
		handler EndpointHandler
	}{pattern, handler})
}

func (r *Router) Route(req *HTTPRequest) EndpointHandler {
	// First check exact routes
	if handler, ok := r.routes[req.URI]; ok {
		return handler
	}

	// Then check pattern routes
	for _, route := range r.patternRoutes {
		if strings.HasPrefix(req.URI, route.pattern) {
			return route.handler
		}
	}

	// Return 404 handler if no match
	return func(w *ResponseWriter, r *HTTPRequest) {
		w.WriteStatus(404)
		w.WriteHeader("Content-Type", "text/plain")
	}
}
