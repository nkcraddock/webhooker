package testhelp

import (
	"net/http"
	"net/http/httptest"
)

type TestServer struct {
	Handler http.Handler
}

func (srv *TestServer) GET(path string) *httptest.ResponseRecorder {
	return srv.request("GET", path)
}

func (srv *TestServer) request(verb, path string) *httptest.ResponseRecorder {
	r, _ := http.NewRequest(verb, path, nil)
	w := httptest.NewRecorder()
	srv.Handler.ServeHTTP(w, r)
	return w
}
