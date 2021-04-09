package main

import (
	"fmt"
	"log"
	"net/http"
)

type Engine struct {

}

// 实现 http 下的 Handler interface
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.Path=%q", req.URL.Path)
	case "/hello":
		for k, v := range req.Header {
			fmt.Fprintf(w, "header %q=%q\n", k, v)
		}
	}
}

func main() {
	engine := new(Engine)
	log.Fatal(http.ListenAndServe(":3001", engine))
}
