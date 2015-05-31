package testhelp

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
)

type TestServer struct {
	Handler http.Handler
}

func (srv *TestServer) GET(path string, data interface{}) *httptest.ResponseRecorder {
	return srv.request("GET", path, nil, data)
}

func (srv *TestServer) POST(path string, data, result interface{}) *httptest.ResponseRecorder {
	body, _ := json.Marshal(data)
	return srv.request("POST", path, body, result)
}

func (srv *TestServer) request(verb, path string, data []byte, result interface{}) *httptest.ResponseRecorder {
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(data)

	r, _ := http.NewRequest(verb, path, buf)
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")

	w := httptest.NewRecorder()
	srv.Handler.ServeHTTP(w, r)

	if result != nil {
		body, err := ioutil.ReadAll(w.Body)
		if err != nil {
			log.Println("Error reading body", err)
			return w
		}

		if err := json.Unmarshal(body, result); err != nil {
			log.Println("Error unmarshalling data", err, "\nDATA:\n", string(body))
		}
	}

	return w
}
