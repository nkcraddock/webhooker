package domain

type Store interface {
	SaveHook(h *Hook) error
	GetHook(id string) (*Hook, error)
	GetHooks(query string) ([]*Hook, error)
	SaveFilter(f *Filter) error
	GetFilters(hook string) ([]*Filter, error)
	DeleteFilter(hook, id string) error
}
