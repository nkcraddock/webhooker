package main

import (
	"log"
	"net/http"

	"github.com/nkcraddock/webhooker/domain"
	"github.com/nkcraddock/webhooker/mgmt"
)

var (
	cfg   config
	hooks domain.Store
)

func init() {
	cfg = loadConfig()
}

func main() {
	handlers := make([]mgmt.Handler, 0)
	server, err := mgmt.NewMgmtServer(handlers)

	if err != nil {
		log.Println(err)
	}

	http.ListenAndServe(cfg.HostUrl, server)
}
