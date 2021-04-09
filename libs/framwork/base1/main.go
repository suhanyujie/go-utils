package main

import (
	"fmt"
	"log"
	"net/http"
)

func main()  {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/hello", HelloHandler)
	log.Fatal(http.ListenAndServe(":3001", nil))
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
}

func HelloHandler(w http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		fmt.Fprintf(w, "Header [%q] = %q\n", k, v)
	}
}
