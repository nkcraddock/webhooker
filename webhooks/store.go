package webhooks

type Store interface {
	AddHook(*Webhook) error
	AllHooksFor(string) ([]Webhook, error)
	GetHookById(string) (*Webhook, error)
	DeleteHook(string) error
	AddHooker(*Webhooker) error
	AllHookers() ([]Webhooker, error)
	//	DeleteHooker(id string) error
}
