package webhooks

import "github.com/nu7hatch/gouuid"

type Hook struct {
	Id            string `json:"id,omitifempty"`
	CallbackUrl   string `json:"url"`
	Secret        string `json:"secret"`
	RatePerMinute int    `json:"rate"`
}

func NewHook(url string, rate int) *Hook {
	return &Hook{
		Id:            getId(),
		CallbackUrl:   url,
		Secret:        getId(),
		RatePerMinute: rate,
	}
}

func (h *Hook) NewFilter(src, evt, key string) *Filter {
	return &Filter{
		Id:   getId(),
		Src:  src,
		Evt:  evt,
		Key:  key,
		Hook: h.Id,
	}
}

func getId() string {
	id, _ := uuid.NewV4()
	return id.String()
}
