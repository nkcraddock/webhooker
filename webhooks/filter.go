package webhooks

type Filter struct {
	Id       string `json:"id"`
	Src      string `json:"src"`
	Evt      string `json:"evt"`
	Key      string `json:"key"`
	Hook     string `json:"hook"`
	RmqProps string `json:"rmqprops,omitempty"`
}

func (f *Filter) Sanitize() interface{} {
	return sanitizedFilter{Filter: f}
}

type sanitizedFilter struct {
	*Filter
	RmqProps omit `json:"rmqprops,omitempty"`
}
