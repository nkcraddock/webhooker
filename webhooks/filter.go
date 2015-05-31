package webhooks

type Filter struct {
	Id   string `json:"id"`
	Src  string `json:"src"`
	Evt  string `json:"evt"`
	Key  string `json:"key"`
	Hook string `json:"hook"`
}
