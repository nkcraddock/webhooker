package main

import (
	"net/http"
	"os"

	"github.com/emicklei/go-restful"

	"github.com/gorilla/handlers"
	"github.com/michaelklishin/rabbit-hole"

	"github.com/nkcraddock/meathooks/db"
	"github.com/nkcraddock/meathooks/webhooks"
)

var (
	conn   *db.Connection
	cfg    config
	hooks  webhooks.Store
	rabbit *rabbithole.Client
)

func init() {
	cfg = loadConfig()
	init_mongo(cfg.MongoUrl, cfg.MongoDb)
	init_rabbit(cfg.RabbitUri, cfg.RabbitUsername, cfg.RabbitPassword)
}

func init_mongo(url string, database string) {
	var err error
	conn, err = db.Dial(url, database)

	if err != nil {
		panic("Failed to connect to mongo")
	}

	hooks = webhooks.NewMongoStore(conn)
}

func init_rabbit(uri string, username string, password string) {
	var err error
	rabbit, err = rabbithole.NewClient(uri, username, password)

	if err != nil {
		panic("Failed to connect to rabbit")
	}
}

func main() {
	container := restful.NewContainer()
	container.Router(restful.CurlyRouter{})

	webhooks.Register(container, hooks, rabbit)

	server := &http.Server{Addr: cfg.HostUrl, Handler: container}
	server.ListenAndServe()
	// 	chain := alice.New(loggerHandler).Then(router)
	//	http.ListenAndServe(cfg.HostUrl, chain)
}

func loggerHandler(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, next)
}
