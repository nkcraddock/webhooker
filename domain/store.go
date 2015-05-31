package domain

type Store interface {
	SaveHook(h *Hook) error
	GetHook(id string) (*Hook, error)
	GetHooks(query string) ([]*Hook, error)
}
