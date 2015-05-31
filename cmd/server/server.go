package main

import (
	"log"
	"net/http"

	"gopkg.in/redis.v3"

	"github.com/nkcraddock/webhooker/db"
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

func initRedisStore() webhooks.Store {
	return db.RedisHookerStore(func() *redis.Client {
		return redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			PoolSize: 10,
			DB:       1,
		})
	})
}

func initClientHandler() mgmt.Handler {
	var locator client.ResourceLocator

	if cfg.ClientRoot == "" {
		locator = &client.BinDataLocator{}
	} else {
		locator = &client.FsLocator{cfg.ClientRoot}
	}

	return client.NewHandler(locator)
}

func main() {
	store := initRedisStore()
	clientHandler := initClientHandler()
	hookHandler := mgmt.NewHooksHandler(store)

	handlers := []mgmt.Handler{hookHandler}

	server, err := mgmt.NewMgmtServer(clientHandler, handlers)

	if err != nil {
		log.Println(err)
	}

	log.Println("Listening on", cfg.HostUrl)
	http.ListenAndServe(cfg.HostUrl, server)
}
