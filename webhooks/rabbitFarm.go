package webhooks

import "github.com/michaelklishin/rabbit-hole"

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
	}
	r.conn.DeclareQueue(vhost, id, queue)
	r.conn.DeclareExchange(vhost, id, exch)
	r.conn.DeclareBinding(vhost, bind)
}

func (r *rabbitFarm) DeleteUrlQueue(id string) {
	r.conn.DeleteQueue(vhost, id)
	r.conn.DeleteExchange(vhost, id)
}
