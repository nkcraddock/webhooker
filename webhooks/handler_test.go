package webhooks

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestGet(t *testing.T) {
	router := mux.NewRouter()
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/webhooks", nil)

	store := &fakeStore{
		all: func() []Webhook {
			return []Webhook{
				NewWebhook("54dac2b3c7f7324b40000001", "localhost/callback", "*"),
			}
		},
	}

	RegisterHandler(router, store)

	router.ServeHTTP(w, r)

	response := parseResponse(w.Body)

	if len(response) != 1 {
		t.Errorf("Got the wrong number of results. Got %d, expected 1", len(response))
	}
}

func parseResponse(r io.Reader) []Webhook {
	var hooks []Webhook
	b, _ := ioutil.ReadAll(r)
	json.Unmarshal(b, &hooks)
	return hooks
}

type fakeStore struct {
	add     func(*Webhook) error
	all     func() []Webhook
	getById func(string) Webhook
}

func (f *fakeStore) Add(wh *Webhook) error {
	return f.add(wh)
}

func (f *fakeStore) All() []Webhook {
	return f.all()
}

func (f *fakeStore) GetById(id string) Webhook {
	return f.getById(id)
}
