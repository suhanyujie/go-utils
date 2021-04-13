package engine

import (
	"net/http"
)

type Engine struct {
	router *router
}

func New() *Engine {
	return &Engine{
		router: newRouter(),
	}
}

type HandlerFunc func(c *Context)

func (e *Engine) Run(addr string) {
	http.ListenAndServe(addr, e)
}

// 实现 http 下的 Handler interface
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	e.router.handle(c)
}
