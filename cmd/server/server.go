package main

import (
	"log"
	"net/http"

	"github.com/nkcraddock/webhooker/mgmt"
	"github.com/nkcraddock/webhooker/mgmt/client"
	"github.com/nkcraddock/webhooker/webhooks"
)

var (
	cfg   config
	hooks webhooks.Store
)

func init() {
	cfg = loadConfig()
}

func main() {
	var locator client.ResourceLocator

	if cfg.ClientRoot == "" {
		locator = &client.BinDataLocator{}
	} else {
		locator = &client.FsLocator{cfg.ClientRoot}
	}

	clienthandler := client.NewHandler(locator)

	handlers := []mgmt.Handler{clienthandler}

	server, err := mgmt.NewMgmtServer(handlers)

	if err != nil {
		log.Println(err)
	}

	log.Println("Listening on", cfg.HostUrl)
	http.ListenAndServe(cfg.HostUrl, server)
}
