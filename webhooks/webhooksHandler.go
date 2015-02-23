package webhooks

import (
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful"
)

type webhooksHandler struct {
	hooks Store
}

func RegisterHooks(c *restful.Container, store Store) {
	handler := webhooksHandler{hooks: store}

	ws := new(restful.WebService)
	ws.Path("/api/webhooks").
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
func (h *webhooksHandler) Create(req *restful.Request, res *restful.Response) {
	hook := new(Webhook)
	err := req.ReadEntity(&hook)

	if failOnError(res, err) {
		return
	}

	hook.Id = getId()

	err = h.hooks.AddHook(hook)

	if failOnError(res, err) {
		return
	}

	uri := fmt.Sprintf("/webhooks/%s", hook.Id)
	res.AddHeader("Location", uri)
	res.WriteHeader(http.StatusCreated)
}

func (h *webhooksHandler) List(req *restful.Request, res *restful.Response) {
	hooks, err := h.hooks.AllHooksFor("12341234")

	if failOnError(res, err) {
		return
	}

	res.WriteEntity(hooks)
}

func (h *webhooksHandler) Get(req *restful.Request, res *restful.Response) {
	id := req.PathParameter("id")
	hook, _ := h.hooks.GetHookById(id)

	if len(hook.Id) == 0 {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	res.WriteEntity(hook)
}

func (h *webhooksHandler) Delete(req *restful.Request, res *restful.Response) {
	id := req.PathParameter("id")
	err := h.hooks.DeleteHook(id)

	if failOnError(res, err) {
		return
	}

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
