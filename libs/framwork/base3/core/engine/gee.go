package engine

import (
	"fmt"
	"net/http"
)

type Engine struct {
	router map[string]HandlerFunc
}

func New() *Engine {
	return &Engine{}
}

type HandlerFunc func(w http.ResponseWriter, req *http.Request)

func (e *Engine) addRoute(method string, pattern string, handler HandlerFunc)  {
	key := fmt.Sprintf("%s-%s", method, pattern)
	e.router[key] = handler
}

func (e *Engine) Get(pattern string, handler HandlerFunc)  {
	e.addRoute("GET", pattern, handler)
}

func (e *Engine) Post(pattern string, handler HandlerFunc)  {
	e.addRoute("POST", pattern, handler)
}

func (e *Engine) Run(addr string)  {
	http.ListenAndServe(addr, e)
}

// 实现 http 下的 Handler interface
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	reqKey := fmt.Sprintf("%s-%s", req.Method, req.URL.Path)
	if handler, ok := e.router[reqKey]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 not found: %s", req.URL.Path)
	}
}
