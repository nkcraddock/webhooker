package webhooks

type memoryStore struct {
	hookers map[string]Webhooker
	hooks   map[string]Webhook
}

func newMemoryStore() *memoryStore {
	return &memoryStore{
		hookers: make(map[string]Webhooker),
		hooks:   make(map[string]Webhook),
	}
}

func (ms *memoryStore) getAll(hooker string) (hooks []Webhook) {
	hooks = make([]Webhook, 0)
	for _, h := range ms.hooks {
		hooks = append(hooks, h)
	}
	return hooks
}
