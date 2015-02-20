package webhooks

type Webhook struct {
	Id          string `json:"id" bson:"_id"`
	CallbackURL string `json:"callback_url" bson:"callback_url"`
	Filter      string `json:"filter" bson:"filter"`
}

func NewWebhook(id string, callback string, filter string) Webhook {
	return Webhook{
		Id:          id,
		CallbackURL: callback,
		Filter:      filter,
	}
}
