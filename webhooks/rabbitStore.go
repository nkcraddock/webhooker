package webhooks

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"github.com/michaelklishin/rabbit-hole"
)

const sourceExchange = "amq.topic"

type RabbitStore struct {
	rb    *rabbithole.Client
	vh    string
	store *memoryStore
}

func NewRabbitStore(r *rabbithole.Client, vh string) Store {
	setupVhost(r, vh)

	s := &RabbitStore{
		rb: r,
		vh: vh,
	}

	// Load the webhooks from rabbit into a memory store
	ms, _ := s.reloadMemoryStore()
	s.store = ms

	return s
}

// HOOKERS

func (r *RabbitStore) AddHooker(hooker *Webhooker) (err error) {
	if _, ok := r.store.hookers[hooker.Id]; ok {
		return fmt.Errorf("This Webhooker is already registered '%s'", hooker.Callback)
	}

	r.store.hookers[hooker.Id] = *hooker
	r.putHooker(hooker)

	return
}

func (r *RabbitStore) AllHookers() (hookers []Webhooker, err error) {
	hookers = make([]Webhooker, 0)

	for _, h := range r.store.hookers {
		hookers = append(hookers, h)
	}

	return
}

func (r *RabbitStore) DeleteHooker(id string) (err error) {
	_, ok := r.store.hookers[id]

	if !ok {
		return fmt.Errorf("No such hooker '%s'", id)
	}

	_, err = r.rb.DeleteQueue(r.vh, id)

	if err != nil {
		return
	}

	_, err = r.rb.DeleteExchange(r.vh, id)

	if err != nil {
		return
	}

	delete(r.store.hookers, id)

	// DELETE ALL THE HOOKS TOO

	return
}

func (r *RabbitStore) GetHooker(id string) (hooker *Webhooker, err error) {
	h, ok := r.store.hookers[id]

	if !ok {
		return nil, fmt.Errorf("No such hooker '%s'", id)
	}

	return &h, nil
}

// HOOKS

func (r *RabbitStore) AllHooksFor(hooker string) (hooks []Webhook, err error) {
	hooks = r.store.getAll(hooker)
	return
}

func (r RabbitStore) GetHookById(id string) (hook *Webhook, err error) {
	return &Webhook{}, nil
}

func (r *RabbitStore) DeleteHook(id string) (err error) {
	hook, ok := r.store.hooks[id]
	if !ok {
		return fmt.Errorf("No such hook '%s'", id)
	}

	_, err = r.rb.DeleteBinding(r.vh, rabbithole.BindingInfo{
		Source:          sourceExchange,
		Destination:     hook.Hooker,
		DestinationType: "exchange",
		RoutingKey:      getHookFilter(&hook),
		PropertiesKey:   hook.Props,
	})

	if err != nil {
		return
	}

	delete(r.store.hooks, id)

	return
}

func (r *RabbitStore) AddHook(wh *Webhook) (err error) {
	// make sure theres a hooker
	_, ok := r.store.hookers[wh.Hooker]

	if !ok {
		return fmt.Errorf("Couldn't find the hooker '%s'", wh.Hooker)
	}

	filter := getHookFilter(wh)

	args := map[string]interface{}{
		"id":  wh.Id,
		"src": wh.Src,
		"evt": wh.Evt,
		"key": wh.Key,
	}
	props, err := r.bindExchange(sourceExchange, wh.Hooker, filter, args)

	wh.Props = props

	r.store.hooks[wh.Id] = *wh

	return
}

func setupVhost(r *rabbithole.Client, vh string) (err error) {
	_, err = r.PutVhost(vh, rabbithole.VhostSettings{Tracing: false})

	if err != nil {
		return
	}

	permissions := rabbithole.Permissions{Configure: ".*", Write: ".*", Read: ".*"}
	_, err = r.UpdatePermissionsIn(vh, r.Username, permissions)

	return
}

func generateEndpointQueueName(url string) string {
	h := sha1.New()
	h.Write([]byte(url))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func (r *RabbitStore) putHooker(h *Webhooker) (err error) {
	args := map[string]interface{}{
		"callback": h.Callback,
		"secret":   h.Secret,
		"email": h.Email,
		"name": h.Name,
	}

	err = r.createQueue(h.Id, args)

	if err != nil {
		return
	}

	err = r.createExchange(h.Id)

	if err != nil {
		return
	}

	err = r.bindQueue(h.Id)

	return
}

func (r *RabbitStore) createQueue(name string, args map[string]interface{}) (err error) {
	_, err = r.rb.DeclareQueue(r.vh, name, rabbithole.QueueSettings{
		Durable:    false,
		AutoDelete: false,
		Arguments:  args,
	})

	return
}

func (r *RabbitStore) createExchange(name string) (err error) {
	_, err = r.rb.DeclareExchange(r.vh, name, rabbithole.ExchangeSettings{
		Type: "topic",
	})

	return
}

func (r *RabbitStore) bindQueue(qn string) (err error) {
	_, err = r.rb.DeclareBinding(r.vh, rabbithole.BindingInfo{
		Source:          qn,
		Destination:     qn,
		DestinationType: "queue",
		RoutingKey:      "#",
	})

	return
}

func (r *RabbitStore) bindExchange(src, dst, filter string, args map[string]interface{}) (props string, err error) {
	res, err := r.rb.DeclareBinding(r.vh, rabbithole.BindingInfo{
		Source:          src,
		Destination:     dst,
		DestinationType: "exchange",
		RoutingKey:      filter,
		Arguments:       args,
	})

	if err != nil {
		return
	}

	props, err = url.QueryUnescape(strings.Split(res.Header.Get("Location"), "/")[1])

	return
}

func (r *RabbitStore) reloadMemoryStore() (s *memoryStore, err error) {
	s = newMemoryStore()

	qs, err := r.rb.ListQueuesIn(r.vh)
	if err != nil {
		return
	}

	for _, q := range qs {
		secret, _ := q.Arguments["secret"].(string)
		callback, _ := q.Arguments["callback"].(string)
		s.hookers[q.Name] = Webhooker{
			Id:       q.Name,
			Secret:   secret,
			Callback: callback,
		}
	}

	bs, err := r.rb.ListBindingsIn(r.vh)
	if err != nil {
		return
	}

	for _, b := range bs {
		// We're only interested in exchange->exchange bindings
		if b.DestinationType != "exchange" {
			continue
		}

		id, _ := b.Arguments["id"].(string)
		src, _ := b.Arguments["src"].(string)
		evt, _ := b.Arguments["evt"].(string)
		key, _ := b.Arguments["key"].(string)

		s.hooks[id] = Webhook{
			Id:     id,
			Src:    src,
			Evt:    evt,
			Key:    key,
			Hooker: b.Destination,
			Props:  b.PropertiesKey,
		}
	}

	return
}

func getHookFilter(wh *Webhook) string {
	return fmt.Sprintf("%s.%s.%s", wh.Src, wh.Evt, wh.Key)
}
