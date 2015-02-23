package webhooks

import (
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful"
)

type webhookersHandler struct {
	hooks Store
}

type HookerRegistration struct {
	Callback string `json:"callback"`
}

func RegisterHookers(c *restful.Container, store Store) {
	handler := webhookersHandler{hooks: store}

	ws := new(restful.WebService)
	ws.Path("/api/webhookers").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/").To(handler.List).
		Doc("get all webhookers").
		Operation("List").
		Returns(200, "OK", []Webhooker{}))

	ws.Route(ws.GET("/{id}").To(handler.Get).
		Doc("get a webhooker").
		Operation("Get").
		Param(ws.PathParameter("id", "identifier of the webhooker").DataType("string")).
		Writes(Webhooker{}))

	ws.Route(ws.POST("/").To(handler.Create).
		Doc("create a webhooker").
		Operation("Create").
		Reads(HookerRegistration{}))

	ws.Route(ws.DELETE("/{id}").To(handler.Delete).
		Doc("delete a webhooker").
		Operation("Delete").
		Param(ws.PathParameter("id", "identifier of the webhook").DataType("string")))

	c.Add(ws)
}

func (h *webhookersHandler) Create(req *restful.Request, res *restful.Response) {
	reg := new(HookerRegistration)
	err := req.ReadEntity(&reg)

	if failOnError(res, err) {
		return
	}

	hooker := NewWebHooker(reg.Callback)
	err = h.hooks.AddHooker(hooker)

	if failOnError(res, err) {
		return
	}

	uri := fmt.Sprintf("/webhooker/%s", hooker.Id)
	res.AddHeader("Location", uri)
	res.WriteHeader(http.StatusCreated)
}

func (h *webhookersHandler) List(req *restful.Request, res *restful.Response) {
	hookers, err := h.hooks.AllHookers()

	if failOnError(res, err) {
		return
	}

	res.WriteEntity(hookers)
}

func (h *webhookersHandler) Get(req *restful.Request, res *restful.Response) {
	id := req.PathParameter("id")
	hooker, _ := h.hooks.GetHooker(id)

	if len(hooker.Id) == 0 {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	res.WriteEntity(hooker)
}

func (h *webhookersHandler) Delete(req *restful.Request, res *restful.Response) {
	id := req.PathParameter("id")
	err := h.hooks.DeleteHooker(id)

	if failOnError(res, err) {
		return
	}

	res.WriteHeader(http.StatusNotFound)
}
