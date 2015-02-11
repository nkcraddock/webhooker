package db

import "gopkg.in/mgo.v2"

type Connection struct {
	database string
	url      string
	mongo    *mgo.Session
}

func Dial(url string, database string) (*Connection, error) {
	conn := &Connection{database: database, url: url}
	session, err := conn.GetSession()

	if err == nil {
		session.Close() // I promise it's temporary
	}

	return conn, err
}

func (r *Connection) GetSession() (*mgo.Session, error) {
	if r.mongo == nil {
		var err error
		r.mongo, err = mgo.Dial(r.url)
		if err != nil {
			return nil, err
		}
	}
	return r.mongo.Clone(), nil
}

func (r *Connection) WithCollection(collection string, s func(*mgo.Collection) error) error {
	session, err := r.GetSession()
	if err != nil {
		return err
	}
	defer session.Close()
	c := session.DB(r.database).C(collection)
	return s(c)
}
