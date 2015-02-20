package webhooks

type Store interface {
	Add(*Webhook) error
	All() ([]Webhook, error)
	GetById(string) Webhook
	Delete(string) error
}
