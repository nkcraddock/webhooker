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

func (r *RabbitStore) AllHooksFor(hooker string) (hooks []Webhook, err error) {
	hooks = r.store.getAll(hooker)
	return
}

func (r RabbitStore) GetHookById(id string) (hook *Webhook, err error) {
	return &Webhook{}, nil
}

func (r *RabbitStore) DeleteHook(id string) (err error) {
	return nil
}

func (r *RabbitStore) AddHook(wh *Webhook) (err error) {
	hooker := r.store.hookers[wh.Hooker]
	qn, err := r.setupUrlQueue(hooker.Callback)

	if err != nil {
		return err
	}

	filter := fmt.Sprintf("%s.%s.%s", wh.Evt, wh.Src, wh.Key)
	props, err := r.bindExchange(sourceExchange, qn, filter)

	wh.Id = props

	return
}

func setupVhost(r *rabbithole.Client, vh string) (err error) {
	permissions := rabbithole.Permissions{Configure: "*", Write: "*", Read: "*"}
	_, err = r.PutVhost(vh, rabbithole.VhostSettings{Tracing: false})

	if err != nil {
		return
	}

	_, err = r.UpdatePermissionsIn(vh, r.Username, permissions)

	return
}

func (r *RabbitStore) setupUrlQueue(url string) (qn string, err error) {
	qn, err = r.createQueue(url)

	if err != nil {
		return
	}

	err = r.createExchange(qn)

	if err != nil {
		return
	}

	err = r.bindQueue(qn)

	return
}

func generateEndpointQueueName(url string) string {
	h := sha1.New()
	h.Write([]byte(url))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func (r *RabbitStore) createQueue(url string) (name string, err error) {
	name = generateEndpointQueueName(url)
	_, err = r.rb.DeclareQueue(r.vh, name, rabbithole.QueueSettings{
		Durable:    false,
		AutoDelete: false,
		Arguments:  map[string]interface{}{"url": url},
	})

	return
}

func (r *RabbitStore) createExchange(name string) (err error) {
	_, err = r.rb.DeclareExchange(r.vh, name, rabbithole.ExchangeSettings{
		Type: "direct",
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

func (r *RabbitStore) bindExchange(src, dst, filter string) (props string, err error) {
	res, err := r.rb.DeclareBinding(r.vh, rabbithole.BindingInfo{
		Source:          src,
		Destination:     dst,
		DestinationType: "exchange",
		RoutingKey:      filter,
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

	return
}

// func (r *RabbitStore) getQueueUrls() (urls map[string]*endpoint, err error) {
// 	urls = make(map[string]*endpoint)

// 	qs, err := r.rb.ListQueuesIn(r.vh)
// 	if err != nil {
// 		return
// 	}

// 	for _, q := range qs {
// 		if arg, ok := q.Arguments["url"]; ok {
// 			if url, ok := arg.(string); ok {
// 				urls[q.Name] = &endpoint{url: url, hooks: make([]Webhook, 0)}
// 			}
// 		}
// 	}

// 	return
// }

// func (r *RabbitStore) reloadEndpoints() {
// 	endpoints, err := r.getQueueUrls()

// 	if err != nil {
// 		return
// 	}

// 	bs, err := r.rb.ListBindingsIn(r.vh)

// 	if err != nil {
// 		return
// 	}

// 	for _, b := range bs {
// 		if b.DestinationType == "exchange" {
// 			if ep, ok := endpoints[b.Destination]; ok {
// 				h := NewWebhook(b.PropertiesKey, ep.url, b.RoutingKey)
// 				ep.hooks = append(ep.hooks, h)
// 			}
// 		}
// 	}

// 	r.endpoints = endpoints
// }
