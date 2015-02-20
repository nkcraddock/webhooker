package webhooks

import (
	"crypto/sha1"
	"encoding/base64"
	"net/url"
	"strings"

	"github.com/michaelklishin/rabbit-hole"
)

const sourceExchange = "amq.topic"

type RabbitStore struct {
	rb        *rabbithole.Client
	vh        string
	endpoints map[string]*endpoint
}

type endpoint struct {
	url   string
	hooks []Webhook
}

func NewRabbitStore(r *rabbithole.Client, vh string) Store {
	setupVhost(r, vh)

	s := &RabbitStore{
		rb: r,
		vh: vh,
	}

	s.reloadEndpoints()

	return s
}

func (r *RabbitStore) All() (hooks []Webhook, err error) {
	hooks = make([]Webhook, 0)
	for _, v := range r.endpoints {
		hooks = append(hooks, v.hooks...)
	}
	return
}

func (r RabbitStore) GetById(id string) Webhook {
	return Webhook{}
}

func (r *RabbitStore) Delete(id string) error {
	return nil
}

func (r *RabbitStore) Add(wh *Webhook) (err error) {
	qn, err := r.setupUrlQueue(wh.CallbackURL)

	if err != nil {
		return err
	}

	props, err := r.bindExchange(sourceExchange, qn, wh.Filter)

	wh.Id = props

	r.reloadEndpoints()

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

func (r *RabbitStore) getQueueUrls() (urls map[string]*endpoint, err error) {
	urls = make(map[string]*endpoint)

	qs, err := r.rb.ListQueuesIn(r.vh)
	if err != nil {
		return
	}

	for _, q := range qs {
		if arg, ok := q.Arguments["url"]; ok {
			if url, ok := arg.(string); ok {
				urls[q.Name] = &endpoint{url: url, hooks: make([]Webhook, 0)}
			}
		}
	}

	return
}

func (r *RabbitStore) reloadEndpoints() {
	endpoints, err := r.getQueueUrls()

	if err != nil {
		return
	}

	bs, err := r.rb.ListBindingsIn(r.vh)

	if err != nil {
		return
	}

	for _, b := range bs {
		if b.DestinationType == "exchange" {
			if ep, ok := endpoints[b.Destination]; ok {
				h := NewWebhook(b.PropertiesKey, ep.url, b.RoutingKey)
				ep.hooks = append(ep.hooks, h)
			}
		}
	}

	r.endpoints = endpoints
}
