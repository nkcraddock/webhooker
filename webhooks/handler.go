package webhooks

import (
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful"

	"github.com/michaelklishin/rabbit-hole"
)

type WebhooksResource struct {
	hooks  Store
	rabbit *rabbitFarm
}

func Register(c *restful.Container, store Store, rabbit *rabbithole.Client) {
	handler := WebhooksResource{
		hooks:  store,
		rabbit: newRabbitFarm(rabbit),
	}

	ws := new(restful.WebService)
	ws.Path("/webhooks").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/").To(handler.List).
		Doc("get all webhooks").
		Operation("List").
		Returns(200, "OK", []Webhook{}))

	ws.Route(ws.GET("/{id}").To(handler.Get).
		Doc("get a webhook").
		Operation("Get").
		Param(ws.PathParameter("id", "identifier of the webhook").DataType("string")).
		Writes(Webhook{}))

	ws.Route(ws.PUT("/").To(handler.Create).
		Doc("create a webhook").
		Operation("Create").
		Reads(Webhook{}))

	ws.Route(ws.DELETE("/{id}").To(handler.Delete).
		Doc("delete a webhook").
		Operation("Delete").
		Param(ws.PathParameter("id", "identifier of the webhook").DataType("string")))

	c.Add(ws)
}

// POST /webhooks
func (h *WebhooksResource) Create(req *restful.Request, res *restful.Response) {
	hook := new(Webhook)
	err := req.ReadEntity(&hook)

	if failOnError(res, err) {
		return
	}

	err = h.hooks.Add(hook)

	if failOnError(res, err) {
		return
	}

	h.rabbit.SaveUrlQueue(hook.Id.Hex())

	uri := fmt.Sprintf("/webhooks/%s", hook.Id.Hex())
	res.AddHeader("Location", uri)
	res.WriteHeader(http.StatusCreated)
}

func (h *WebhooksResource) List(req *restful.Request, res *restful.Response) {
	hooks := h.hooks.All()
	res.WriteEntity(hooks)
}

func (h *WebhooksResource) Get(req *restful.Request, res *restful.Response) {
	id := req.PathParameter("id")
	hook := h.hooks.GetById(id)

	if len(hook.Id) == 0 {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	res.WriteEntity(hook)
}

func (h *WebhooksResource) Delete(req *restful.Request, res *restful.Response) {
	id := req.PathParameter("id")
	err := h.hooks.Delete(id)

	if failOnError(res, err) {
		return
	}

	h.rabbit.DeleteUrlQueue(id)

	res.WriteHeader(http.StatusNotFound)
}

func failOnError(response *restful.Response, err error) bool {
	if err == nil {
		return false
	}

	msg := fmt.Sprintf("An error occurred: %s", err.Error())
	response.WriteErrorString(http.StatusInternalServerError, msg)
	return true
}
