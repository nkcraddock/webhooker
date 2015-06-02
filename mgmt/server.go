package mgmt

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Handler interface {
	RegisterRoutes(r *mux.Router)
}

type MgmtServer struct {
	m        *mux.Router
	handlers []Handler
}

func NewMgmtServer(catchall Handler, handlers []Handler) (*MgmtServer, error) {
	m := mux.NewRouter()

	// Let all the handlers register their routes
	api := m.PathPrefix("/api").Subrouter()
	for _, h := range handlers {
		h.RegisterRoutes(api)
	}

	// wire up the catchall route to the main router
	catchall.RegisterRoutes(m)

	return &MgmtServer{m, handlers}, nil
}

func (s *MgmtServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.m.ServeHTTP(w, r)
}
