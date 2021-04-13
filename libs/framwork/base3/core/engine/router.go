package engine

import (
	"fmt"
	"net/http"
)

type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
	}
}

func (e *router) addRoute(method string, pattern string, handler HandlerFunc) {
	key := fmt.Sprintf("%s-%s", method, pattern)
	e.handlers[key] = handler
}

func (r *router) Get(pattern string, handler HandlerFunc) {
	r.addRoute("GET", pattern, handler)
}

func (r *router) Post(pattern string, handler HandlerFunc) {
	r.addRoute("POST", pattern, handler)
}

func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
