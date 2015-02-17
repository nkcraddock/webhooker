package main

import (
	"net/http"

	"github.com/emicklei/go-restful/swagger"

	"github.com/emicklei/go-restful"

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

func init_swagger(container *restful.Container) {
	swag := swagger.Config{
		WebServices: container.RegisteredWebServices(),
		// WebServicesUrl:  "http://localhost:3001",
		ApiPath:         "/apidocs.json",
		SwaggerPath:     "/apidocs/",
		SwaggerFilePath: "/home/nathan/dev/swagger-ui/dist",
	}

	swagger.RegisterSwaggerService(swag, container)
}

func main() {
	container := restful.NewContainer()
	webhooks.Register(container, hooks, rabbit)

	init_swagger(container)
	http.ListenAndServe(cfg.HostUrl, container)
}
