package q

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/michaelklishin/rabbit-hole"
	"github.com/nkcraddock/webhooker/webhooks"
	"github.com/nu7hatch/gouuid"
)

const sourceExchange = "amq.topic"

type RabbitStore struct {
	src string
	rb  *rabbithole.Client
	vh  string
}

func NewRabbitStore(r *rabbithole.Client, vh string) *RabbitStore {
	setupVhost(r, vh)

	s := &RabbitStore{
		rb: r,
		vh: vh,
	}

	return s
}

func (r *RabbitStore) SaveHook(hook *webhooks.Hook) error {
	if hook.Id == "" {
		hook.Id = getId()
	}

	// Create a queue
	if err := r.createQueue(hook.Id, nil); err != nil {
		return err
	}

	// Create an exchange
	if err := r.createExchange(hook.Id); err != nil {
		return err
	}

	// BIND THEM!
	if err := r.bindQueue(hook.Id); err != nil {
		return err
	}

	return nil
}

func (r *RabbitStore) GetHook(id string) (*webhooks.Hook, error) {
	q, err := r.rb.GetQueue(r.vh, id)
	if err != nil {
		return nil, err
	}

	return &webhooks.Hook{Id: q.Name}, nil
}

func (r *RabbitStore) GetHooks(query string) ([]*webhooks.Hook, error) {

	qs, err := r.rb.ListQueuesIn(r.vh)
	if err != nil {
		return nil, err
	}

	hooks := make([]*webhooks.Hook, len(qs))
	for ix, q := range qs {
		hooks[ix] = &webhooks.Hook{Id: q.Name}
	}

	return hooks, nil
}

func (r *RabbitStore) SaveFilter(f *webhooks.Filter) error {
	if f.Hook == "" {
		return fmt.Errorf("Invalid hook")
	}

	hook, err := r.GetHook(f.Hook)
	if err != nil {
		return fmt.Errorf("Invalid hook %s", f.Hook)
	}

	if f.Id == "" {
		f.Id = getId()
	}

	topic := getTopicFilter(f)
	args := map[string]interface{}{
		"FilterId": f.Id,
	}

	if prop, err := r.bindExchange(sourceExchange, hook.Id, topic, args); err != nil {
		return err
	} else {
		// Store the properties key. We need it later to delete this thing.
		f.RmqProps = prop
	}

	return nil
}

func (r *RabbitStore) GetFilter(hook, id string) (*webhooks.Filter, error) {
	bs, err := r.rb.ListBindingsIn(r.vh)
	if err != nil {
		return nil, err
	}

	for _, b := range bs {
		if filterId, ok := b.Arguments["FilterId"]; ok {
			if fid, ok := filterId.(string); ok {
				return &webhooks.Filter{Id: fid}, nil
			}
		}
	}

	return nil, fmt.Errorf("Not found")
}

func (r *RabbitStore) GetFilters(hook string) ([]*webhooks.Filter, error) {
	bs, err := r.rb.ListQueueBindings(r.vh, hook)
	if err != nil {
		return nil, err
	}

	filters := make([]*webhooks.Filter, len(bs))
	for ix, b := range bs {
		filters[ix] = &webhooks.Filter{
			Id: b.PropertiesKey,
		}
	}

	return filters, nil
}

func setupVhost(r *rabbithole.Client, vh string) error {
	if _, err := r.PutVhost(vh, rabbithole.VhostSettings{Tracing: false}); err != nil {
		return err
	}

	permissions := rabbithole.Permissions{Configure: ".*", Write: ".*", Read: ".*"}
	_, err := r.UpdatePermissionsIn(vh, r.Username, permissions)

	return err
}

func (r *RabbitStore) createQueue(name string, args map[string]interface{}) error {
	_, err := r.rb.DeclareQueue(r.vh, name, rabbithole.QueueSettings{
		Durable:    false,
		AutoDelete: false,
		Arguments:  args,
	})

	return err
}

func (r *RabbitStore) createExchange(name string) error {
	_, err := r.rb.DeclareExchange(r.vh, name, rabbithole.ExchangeSettings{
		Type: "topic",
	})

	return err
}

func (r *RabbitStore) bindQueue(qn string) error {
	_, err := r.rb.DeclareBinding(r.vh, rabbithole.BindingInfo{
		Source:          qn,
		Destination:     qn,
		DestinationType: "queue",
		RoutingKey:      "#",
	})

	return err
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

func getId() string {
	id, _ := uuid.NewV4()
	return id.String()
}

func getTopicFilter(f *webhooks.Filter) string {
	return fmt.Sprintf("%s.%s.%s", f.Src, f.Evt, f.Key)
}
