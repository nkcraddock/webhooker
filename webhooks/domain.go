package webhooks

import (
	"github.com/nu7hatch/gouuid"
)

type Webhooker struct {
	Id       string `json:"id" bson:"id"`
	Secret   string `json:"-" bson:"secret"`
	Callback string `json:"callback" bson:"callback"`
}

type Webhook struct {
	Id     string `json:"id" bson:"id"`
	Src    string `json:"src" bson:"src"`
	Evt    string `json:"evt" bson:"evt"`
	Key    string `json:"key" bson:"key"`
	Props  string `json:"props" bson:"props"`
	Hooker string `json:"hooker" bson:"hooker"`
}

func NewWebHooker(callback string) *Webhooker {
	return &Webhooker{
		Id:       getId(),
		Secret:   getId(),
		Callback: callback,
	}
}

func NewWebhook(src, evt, key, hooker string) Webhook {
	return Webhook{
		Id:     getId(),
		Src:    src,
		Evt:    evt,
		Key:    key,
		Hooker: hooker,
	}
}

func getId() string {
	id, _ := uuid.NewV4()
	return id.String()
}
