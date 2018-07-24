package commons

import (
	"reflect"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

var (
	session *mgo.Session
	mongo   *Mongo
)

type Mongo struct {
	info *mgo.DialInfo
}

func NewMongo(ip, port, dataBase, user, pwd string, poolLimit int) *Mongo {
	if mongo == nil {
		m.Lock()
		defer m.Unlock()
		if mongo == nil {
			mongo = &Mongo{
				info: &mgo.DialInfo{
					Addrs:     []string{ip},
					Direct:    false,
					Timeout:   time.Second * 1,
					Username:  user,
					Password:  pwd,
					Database:  dataBase,
					PoolLimit: poolLimit,
				},
			}
		}
	}
	return mongo
}
func (m *Mongo) session() *mgo.Session {
	if session == nil {
		var err error
		session, err = mgo.DialWithInfo(m.info)
		if err != nil {
			Console().Panic(err)
		}
	}
	return session.Clone()
}

func (m *Mongo) mc(collection string, f func(*mgo.Collection) error) error {
	session := m.session()
	defer func() {
		session.Close()
		if err := recover(); err != nil {
			Console().Panic(err)
		}
	}()
	c := session.DB(m.info.Database).C(collection)
	return f(c)
}

func (m *Mongo) mdc(dbName string, collection string, f func(*mgo.Collection) error) error {
	session := m.session()
	defer func() {
		session.Close()
		if err := recover(); err != nil {
			Console().Panic(err)
		}
	}()
	c := session.DB(dbName).C(collection)
	return f(c)
}

func (m *Mongo) insert(dbName, collName string, i interface{}) error {
	v := reflect.ValueOf(i).Elem()
	v.FieldByName("ID").Set(reflect.ValueOf(bson.NewObjectId().String()))
	v.FieldByName("CTime").Set(reflect.ValueOf(time.Now()))
	return m.mdc(dbName, collName, func(coll *mgo.Collection) error {
		err := coll.Insert(i)
		return err
	})
}
func (m *Mongo) InsertDC(dbName, collName string, i interface{}) error {
	return m.insert(dbName, collName, i)
}
func (m *Mongo) InsertC(collName string, i interface{}) error {
	return m.insert(m.info.Database, collName, i)
}

func (m *Mongo) view(dbName, collName string, query, result interface{}) error {
	return m.mdc(dbName, collName, func(coll *mgo.Collection) error {
		err := coll.Find(query).All(&result)
		if err != nil {
			return err
		}
		return err
	})
}

func (m *Mongo) ViewAllDC(dbName, collName string, query, result interface{}) error {
	return m.view(dbName, collName, query, result)
}
func (m *Mongo) ViewAllC(collName string, query, result interface{}) error {
	return m.view(m.info.Database, collName, query, result)
}

func (m *Mongo) viewOne(dbName, collName string, query, result interface{}) error {
	return m.mdc(dbName, collName, func(coll *mgo.Collection) error {
		err := coll.Find(query).One(&result)
		if err != nil {
			return err
		}
		return err
	})
}

func (m *Mongo) ViewOneDC(dbName, collName string, query, result interface{}) error {
	return m.viewOne(dbName, collName, query, result)
}
func (m *Mongo) ViewOneC(collName string, query, result interface{}) error {
	return m.viewOne(m.info.Database, collName, query, result)
}
