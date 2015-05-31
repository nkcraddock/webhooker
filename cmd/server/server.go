package main

import (
	"net/http"

	"github.com/emicklei/go-restful/swagger"

	"github.com/emicklei/go-restful"

	"github.com/nkcraddock/webhooker/domain"
)

var (
	cfg   config
	hooks domain.Store
)

func init() {
	cfg = loadConfig()
}

func init_swagger(container *restful.Container) {
	swag := swagger.Config{
		WebServices: container.RegisteredWebServices(),
		ApiPath:     "/api/docs.json",
	}

	swagger.RegisterSwaggerService(swag, container)
}

func main() {
	container := restful.NewContainer()

	init_swagger(container)
	http.ListenAndServe(cfg.HostUrl, container)
}
