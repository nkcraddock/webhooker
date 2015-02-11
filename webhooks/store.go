package webhooks

import (
	"github.com/nkcraddock/meathooks/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Store interface {
	Add(*Webhook) error
	All() []Webhook
	GetById(string) Webhook
}

type mgoStore struct {
	conn *db.Connection
}

func NewMongoStore(conn *db.Connection) Store {
	store := &mgoStore{conn: conn}
	store.registerSchema()
	return store
}

func (s *mgoStore) do(c func(*mgo.Collection) error) error {
	return s.conn.WithCollection("webhooks", c)
}

func (s *mgoStore) registerSchema() error {
	// Set up the webhooks index
	return s.do(func(c *mgo.Collection) (err error) {
		err = c.EnsureIndex(mgo.Index{
			Key:        []string{"CallbackURL"},
			Unique:     true,
			DropDups:   true,
			Background: true,
		})
		return
	})
}

func (s *mgoStore) Add(hook *Webhook) error {
	hook.Id = bson.NewObjectId()
	err := s.do(func(c *mgo.Collection) (err error) {
		err = c.Insert(hook)
		return err
	})

	return err
}

func (s *mgoStore) All() (webhooks []Webhook) {
	err := s.do(func(c *mgo.Collection) (err error) {
		err = c.Find(nil).Sort("callback_url").All(&webhooks)
		return err
	})

	if err != nil {
		panic(err)
	}
	return webhooks
}

func (s *mgoStore) GetById(id string) (webhook Webhook) {
	s.do(func(c *mgo.Collection) (err error) {
		err = c.FindId(bson.ObjectIdHex(id)).One(&webhook)
		return err
	})

	return
}
