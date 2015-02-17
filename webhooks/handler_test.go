package webhooks

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/emicklei/go-restful"
)

func TestList(t *testing.T) {
	container := restful.NewContainer()

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/webhooks", nil)

	store := &fakeStore{
		all: func() []Webhook {
			return []Webhook{
				NewWebhook("54dac2b3c7f7324b40000001", "localhost/callback", "*"),
			}
		},
	}

	Register(container, store, nil)
	container.ServeHTTP(w, r)

	response := parseResponseSet(w.Body)

	if len(response) != 1 {
		t.Errorf("Got the wrong number of results. Got %d, expected 1", len(response))
	}
}

func TestGetById(t *testing.T) {
	container := restful.NewContainer()
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/webhooks/54dac2b3c7f7324b40000001", nil)

	hook := NewWebhook("54dac2b3c7f7324b40000001", "localhost/callback", "*")

	store := &fakeStore{
		getById: func(id string) Webhook {
			return hook
		},
	}

	Register(container, store, nil)
	container.ServeHTTP(w, r)

	response := parseResponse(w.Body)

	if !reflect.DeepEqual(response, hook) {
		t.Errorf("Got the wrong response %e, expected %e", response, hook)
	}

	if w.Code != 200 {
		t.Errorf("Got the wrong response code %d expected 200", w.Code)
	}

}

func parseResponse(r io.Reader) Webhook {
	var hook Webhook
	b, _ := ioutil.ReadAll(r)
	json.Unmarshal(b, &hook)
	return hook
}

func parseResponseSet(r io.Reader) []Webhook {
	var hooks []Webhook
	b, _ := ioutil.ReadAll(r)
	json.Unmarshal(b, &hooks)
	return hooks
}

type fakeStore struct {
	add     func(*Webhook) error
	all     func() []Webhook
	getById func(string) Webhook
	del     func(string) error
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

func (f *fakeStore) Delete(id string) error {
	return f.del(id)
}
