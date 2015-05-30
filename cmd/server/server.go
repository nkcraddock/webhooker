package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/emicklei/go-restful/swagger"

	"github.com/emicklei/go-restful"

	"github.com/michaelklishin/rabbit-hole"

	"github.com/nkcraddock/webhooker/webhooks"
)

var (
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

	hooks = webhooks.NewRabbitStore(rabbit, "webhooker")
}

func init_swagger(container *restful.Container) {
	cur, _ := os.Getwd()
	swag := swagger.Config{
		WebServices:     container.RegisteredWebServices(),
		ApiPath:         "/api/docs.json",
		SwaggerPath:     "/api/docs/",
		SwaggerFilePath: filepath.Join(cur, "swagger"),
	}

	swagger.RegisterSwaggerService(swag, container)
}

func main() {
	container := restful.NewContainer()
	webhooks.RegisterHooks(container, hooks)
	webhooks.RegisterHookers(container, hooks)

	init_swagger(container)
	http.ListenAndServe(cfg.HostUrl, container)
}
