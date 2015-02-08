package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Repo struct {
	database string
	url      string
	mongo    *mgo.Session
}

func ConnectRepo(url string, database string) *Repo {
	return &Repo{database: database, url: url}
}

func (r *Repo) AddWebhook(wh *Webhook) error {
	wh.Id = bson.NewObjectId()
	err := r.withCollection("webhooks", func(c *mgo.Collection) (err error) {
		err = c.Insert(wh)
		return err
	})

	return err
}

func (r *Repo) GetWebhooks() (webhooks []Webhook) {
	err := r.withCollection("webhooks", func(c *mgo.Collection) (err error) {
		err = c.Find(nil).Sort("callback_url").All(&webhooks)
		return err
	})

	if err != nil {
		panic(err)
	}
	return webhooks
}

func (r *Repo) GetWebhook(id string) (webhook Webhook) {

	err := r.withCollection("webhooks", func(c *mgo.Collection) (err error) {
		err = c.FindId(bson.ObjectIdHex(id)).One(&webhook)
		return err
	})

	if err != nil {
		panic(err)
	}

	return webhook
}

func (r *Repo) getSession() *mgo.Session {
	if r.mongo == nil {
		var err error
		r.mongo, err = mgo.Dial(r.url)
		if err != nil {
			panic(err)
		}
	}
	return r.mongo.Clone()
}

func (r *Repo) withCollection(collection string, s func(*mgo.Collection) error) error {
	session := r.getSession()
	defer session.Close()
	c := session.DB(r.database).C(collection)
	return s(c)
}
