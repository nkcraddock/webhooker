package main

import (
	"flag"
	"fmt"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":3001", "the address to listen on")
	flag.Parse()

	http.HandleFunc("/", handler)
	http.ListenAndServe(*addr, nil)
}

func handler(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(rw, "test")
}
