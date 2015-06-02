package webhooks

// Used for evil.
type omit *struct{}

type Hook struct {
	Id            string `json:"id,omitempty"`
	CallbackUrl   string `json:"url"`
	Secret        string `json:"secret,omitempty"`
	RatePerMinute int    `json:"rate"`
}

func (h *Hook) Sanitize() interface{} {
	return sanitizedHook{Hook: h}
}

type sanitizedHook struct {
	*Hook
	Secret omit `json:"secret,omitempty"`
}

func NewHook(url string, rate int) *Hook {
	return &Hook{
		CallbackUrl:   url,
		RatePerMinute: rate,
	}
}

func (h *Hook) NewFilter(src, evt, key string) *Filter {
	return &Filter{
		Src:  src,
		Evt:  evt,
		Key:  key,
		Hook: h.Id,
	}
}
