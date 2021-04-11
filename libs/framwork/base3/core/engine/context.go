package engine

import "net/http"

type Context struct {
	Writer http.ResponseWriter
}
