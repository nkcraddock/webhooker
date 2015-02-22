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
		allHooksFor: func(hooker string) ([]Webhook, error) {
			return []Webhook{
				NewWebhook("src", "evt", "key", "hooker"),
			}, nil
		},
	}

	RegisterHooks(container, store)
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

	hook := NewWebhook("src", "evt", "key", "hooker")

	store := &fakeStore{
		getHookById: func(id string) (*Webhook, error) {
			return &hook, nil
		},
	}

	RegisterHooks(container, store)
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
	addHook       func(*Webhook) error
	allHooksFor   func(string) ([]Webhook, error)
	getHookById   func(string) (*Webhook, error)
	deleteHook    func(string) error
	addHooker     func(*Webhooker) error
	allHookers    func() ([]Webhooker, error)
	getHookerById func(string) (*Webhooker, error)
	deleteHooker  func(string) error
}

func (f *fakeStore) AddHooker(h *Webhooker) error {
	return f.addHooker(h)
}

func (f *fakeStore) AllHookers() ([]Webhooker, error) {
	return f.allHookers()
}

func (f *fakeStore) GetHookerById(id string) (*Webhooker, error) {
	return f.getHookerById(id)
}

func (f *fakeStore) DeleteHooker(id string) error {
	return f.deleteHooker(id)
}

func (f *fakeStore) AddHook(wh *Webhook) error {
	return f.addHook(wh)
}

func (f *fakeStore) AllHooksFor(hooker string) ([]Webhook, error) {
	return f.allHooksFor(hooker)
}

func (f *fakeStore) GetHookById(id string) (*Webhook, error) {
	return f.getHookById(id)
}

func (f *fakeStore) DeleteHook(id string) error {
	return f.deleteHook(id)
}
