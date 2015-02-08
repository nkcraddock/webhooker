package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stderr)
	log.SetLevel(log.DebugLevel)
}

func main() {
	addr := flag.String("addr", ":3001", "the address to listen on")
	flag.Parse()

	r := mux.NewRouter()

	r.HandleFunc("/", home)
	r.HandleFunc("/webhooks", WebhooksPost).Methods("POST")
	r.HandleFunc("/webhooks/{id}", WebhooksGet).Methods("GET")
	r.HandleFunc("/webhooks", WebhooksList).Methods("GET")

	http.Handle("/", r)
	log.Infof("Listening on %s", *addr)
	http.ListenAndServe(*addr, nil)
}

func home(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "HUH")
}
