package storage

import (
	c "github.com/korovaisdead/go-simple-membership/config"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func getSession() (Session, error) {
	config := c.Get()
	if config.Db.Database == "test" {
		return &MongoSession{}, nil
	}

	di := &mgo.DialInfo{
		Addrs: []string{config.Db.Host + config.Db.Port},
	}
	session, err := mgo.DialWithInfo(di)
	if err != nil {
		return nil, err
	}
	return MongoSession{session}, nil
}

type Session interface {
	Close()
	DB(name string) Database
}

type Database interface {
	C(name string) Collection
}

type Collection interface {
	Find(query interface{}) Query
	Insert(docs ...interface{}) error
}

type Query interface {
	One(result interface{}) (err error)
	All(result interface{}) error
}

type MongoDatabase struct {
	*mgo.Database
}

func (d MongoDatabase) C(name string) Collection {
	return &MongoCollection{Collection: d.Database.C(name)}
}

type MongoSession struct {
	*mgo.Session
}

func (s MongoSession) DB(name string) Database {
	return &MongoDatabase{Database: s.Session.DB(name)}
}

type MongoCollection struct {
	*mgo.Collection
}

func (c MongoCollection) Find(query interface{}) Query {
	return MongoQuery{Query: c.Collection.Find(query)}
}

type MongoQuery struct {
	*mgo.Query
}

//User represents the user model inside database
type User struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	Email     string        `bson:"email" json:"email"`
	Firstname string        `bson:"firstname" json:"firstname"`
	Lastname  string        `bson:"lastname" json:"lastname"`
	Password  string        `bson:"password" json:"password"`
	Phone     string        `bson:"phone" json:"phone"`
	Salt      string        `bson:"salt" json:"salt"`
}
