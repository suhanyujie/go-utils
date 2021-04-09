package main

import (
	"github.com/suhanyujie/go-utils/libs/framwork/base3/core/engine"
)

func main() {
	r := engine.New()
	r.Run(":3001")
}
