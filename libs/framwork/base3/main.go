package main

import (
	"fmt"

	"github.com/suhanyujie/go-utils/libs/framwork/base3/core/engine"
)

var (
	Port = 3001
)

func main() {
	r := engine.New()
	r.Get("/test", func(c *engine.Context) {
		c.Json(200, "who are you....")
	})
	r.Get("/query/:keyword", func(c *engine.Context) {
		c.Json(200, fmt.Sprintf("the query param is: %s\n", c.GetParam("keyword")))
	})
	fmt.Printf("server is start in port: %d\n", Port)
	addr := fmt.Sprintf(":%d", Port)
	r.Run(addr)
}
