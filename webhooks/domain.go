package webhooks

import "gopkg.in/mgo.v2/bson"

type Webhook struct {
	Id          bson.ObjectId `json:"id" bson:"_id"`
	CallbackURL string        `json:"callback_url" bson:"callback_url"`
	Filter      string        `json:"filter" bson:"filter"`
}
