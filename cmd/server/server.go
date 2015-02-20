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
	conn  *db.Connection
	cfg   config
	hooks webhooks.Store
)

func init() {
	cfg = loadConfig()
	init_rabbit(cfg.RabbitUri, cfg.RabbitUsername, cfg.RabbitPassword)
}

func init_rabbit(uri string, username string, password string) {
	var err error
	rabbit, err := rabbithole.NewClient(uri, username, password)

	if err != nil {
		panic("Failed to connect to rabbit")
	}

	hooks = webhooks.NewRabbitStore(rabbit, "meathooks")
}

func init_swagger(container *restful.Container) {
	swag := swagger.Config{
		WebServices:     container.RegisteredWebServices(),
		ApiPath:         "/apidocs.json",
		SwaggerPath:     "/apidocs/",
		SwaggerFilePath: "/home/nathan/dev/swagger-ui/dist",
	}

	swagger.RegisterSwaggerService(swag, container)
}

func main() {
	container := restful.NewContainer()
	webhooks.Register(container, hooks)

	init_swagger(container)
	http.ListenAndServe(cfg.HostUrl, container)
}
