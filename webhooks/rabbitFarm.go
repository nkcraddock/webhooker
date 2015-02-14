package webhooks

import (
	"crypto/sha1"
	"encoding/base64"

	"github.com/nkcraddock/rabbit-hole"
)

const vhost = "/"

type rabbitFarm struct {
	conn *rabbithole.Client
}

func newRabbitFarm(r *rabbithole.Client) *rabbitFarm {
	return &rabbitFarm{conn: r}
}

func (r *rabbitFarm) SaveUrlQueue(id string) {
	queue := rabbithole.QueueSettings{
		Durable:    false,
		AutoDelete: true,
	}

	exch := rabbithole.ExchangeSettings{
		Type:       "topic",
		Durable:    false,
		AutoDelete: true,
	}

	bind := rabbithole.BindingInfo{
		Source:          id,
		Vhost:           vhost,
		Destination:     id,
		DestinationType: "q",
		RoutingKey:      "#",
		PropertiesKey:   PropertiesKey(id, "#"),
	}

	r.conn.DeclareQueue(vhost, id, queue)
	r.conn.DeclareExchange(vhost, id, exch)
	r.conn.DeclareBinding(vhost, bind.PropertiesKey, bind)
}

func (r *rabbitFarm) DeleteUrlQueue(id string) {
	bind := rabbithole.BindingInfo{
		Source:          id,
		Vhost:           vhost,
		Destination:     id,
		DestinationType: "q",
		RoutingKey:      "#",
		PropertiesKey:   PropertiesKey(id, "#"),
	}

	r.conn.DeleteBinding(vhost, bind.PropertiesKey, bind)
	r.conn.DeleteQueue(vhost, id)
	r.conn.DeleteExchange(vhost, id)
}

func PropertiesKey(id, filter string) string {
	sha := sha1.New()
	sha.Write([]byte(filter))
	return base64.URLEncoding.EncodeToString(sha.Sum(nil))
}
