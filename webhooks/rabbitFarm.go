package webhooks

import "github.com/nkcraddock/rabbit-hole"

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

	r.conn.DeclareQueue("/", id, queue)
	r.conn.DeclareExchange("/", id, exch)
}
