package client

import (
	"log"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

type Handler struct {
	resources ResourceLocator
}

func NewHandler(res ResourceLocator) *Handler {
	return &Handler{
		resources: res,
	}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.Methods("GET").HandlerFunc(h.ServeHTTP)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	data, err := h.resources.Get(path)

	if err == nil {
		w.Header().Set("Content-Type", getContentType(path))
		if _, err := w.Write(data); err != nil {
			log.Println(err)
		}
		return
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(http.StatusText(http.StatusNotFound)))
}

func getContentType(path string) string {
	if mt := mime.TypeByExtension(filepath.Ext(path)); mt != "" {
		return mt
	}

	return "text/html"
}
