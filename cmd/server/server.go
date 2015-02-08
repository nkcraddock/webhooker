package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stderr)
	log.SetLevel(log.DebugLevel)
}

func main() {
	addr := flag.String("addr", ":3001", "the address to listen on")
	flag.Parse()

	http.HandleFunc("/", handler)
	log.Infof("Listening on %s", *addr)
	http.ListenAndServe(*addr, nil)
}

func handler(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(rw, "test")
}
