package commons

import (
	"reflect"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

var (
	session     *mgo.Session
	mongo       *Mongo
	ErrNotFound = mgo.ErrNotFound
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

type ObjectIDs []ObjectID

func (o ObjectIDs) Hex() string {
	return o.Hex()
}

type ObjectID interface {
	Hex() string
}
type objectID bson.ObjectId

func (o *objectID) Hex() string {
	return o.Hex()
}

func (m *Mongo) ObjectIDHex(id string) ObjectID {
	if len(id) != 24 {
		return nil
	}
	return bson.ObjectIdHex(id)
}

func (m *Mongo) ObjectIDsHex(ids []string) []ObjectID {
	var oIDs []ObjectID
	for _, v := range ids {
		if id := m.ObjectIDHex(v); id != nil {
			oIDs = append(oIDs, id)
		}
	}
	return oIDs
}

func (m *Mongo) NewObjectID() ObjectID {
	return bson.NewObjectId()
}

func (m *Mongo) mc(collection string, f func(*mgo.Collection) (string, error)) (string, error) {
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

func (m *Mongo) mdc(dbName string, collection string, f func(*mgo.Collection) (string, error)) (string, error) {
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

func (m *Mongo) insert(dbName, collName string, i interface{}) (string, error) {
	v := reflect.ValueOf(i).Elem()
	id := v.FieldByName("ID")
	newID := m.NewObjectID()
	if id.IsNil() {
		id.Set(reflect.ValueOf(newID))
	}
	v.FieldByName("CTime").Set(reflect.ValueOf(time.Now()))
	return m.mdc(dbName, collName, func(coll *mgo.Collection) (string, error) {
		return newID.Hex(), coll.Insert(i)
	})
}
func (m *Mongo) InsertDC(dbName, collName string, i interface{}) (string, error) {
	return m.insert(dbName, collName, i)
}
func (m *Mongo) InsertC(collName string, i interface{}) (string, error) {
	return m.insert(m.info.Database, collName, i)
}

func (m *Mongo) upsert(dbName, collName string, q bson.M, i interface{}) (string, error) {
	v := reflect.ValueOf(i).Elem()
	id := v.FieldByName("ID")
	newID := m.NewObjectID()
	if id.IsNil() {
		id.Set(reflect.ValueOf(newID))
	}
	v.FieldByName("CTime").Set(reflect.ValueOf(time.Now()))
	return m.mdc(dbName, collName, func(coll *mgo.Collection) (string, error) {
		info, err := coll.Upsert(q, i)
		if err != nil {
			return "", err
		}
		if info != nil {
			return info.UpsertedId.(bson.ObjectId).Hex(), nil
		}
		return "", nil
	})
}
func (m *Mongo) UpdateDC(dbName, collName string, q, i bson.M) (string, error) {
	return m.update(dbName, collName, q, i)
}
func (m *Mongo) UpdateC(collName string, q, i bson.M) (string, error) {
	return m.update(m.info.Database, collName, q, i)
}

func (m *Mongo) update(dbName, collName string, q, i map[string]interface{}) (string, error) {
	return m.mdc(dbName, collName, func(coll *mgo.Collection) (string, error) {
		data := bson.M{"$set": i}
		err := coll.Update(q, data)
		if err != nil {
			return "", err
		}
		return "", nil
	})
}
func (m *Mongo) UpsertDC(dbName, collName string, q bson.M, i interface{}) (string, error) {
	return m.upsert(dbName, collName, q, i)
}
func (m *Mongo) UpsertC(collName string, q bson.M, i interface{}) (string, error) {
	return m.upsert(m.info.Database, collName, q, i)
}

func (m *Mongo) view(dbName, collName string, query, result interface{}, selectQ ...interface{}) (string, error) {
	return m.mdc(dbName, collName, func(coll *mgo.Collection) (string, error) {
		var q *mgo.Query
		if len(selectQ) > 0 {
			q = coll.Find(query).Select(selectQ[0])
		} else {
			q = coll.Find(query)
		}
		err := q.All(result)
		if err != nil {
			result = nil
			if err != ErrNotFound {
				return "", err
			}
		}
		return "", err
	})
}

func (m *Mongo) ViewAllDC(dbName, collName string, query, result interface{}, selectQ ...interface{}) (string, error) {
	return m.view(dbName, collName, query, result, selectQ...)
}
func (m *Mongo) ViewAllC(collName string, query, result interface{}, selectQ ...interface{}) (string, error) {
	return m.view(m.info.Database, collName, query, result, selectQ...)
}

func (m *Mongo) viewOne(dbName, collName string, query, result interface{}, selectQ ...interface{}) (string, error) {
	return m.mdc(dbName, collName, func(coll *mgo.Collection) (string, error) {
		var q *mgo.Query
		if len(selectQ) > 0 {
			q = coll.Find(query).Select(selectQ[0])
		} else {
			q = coll.Find(query)
		}
		err := q.One(result)
		if err != nil {
			result = nil
			if err != ErrNotFound {
				return "", err
			}
		}
		return "", err
	})
}

func (m *Mongo) ViewOneDC(dbName, collName string, query, result interface{}, selectQ ...interface{}) (string, error) {
	return m.viewOne(dbName, collName, query, result, selectQ...)
}
func (m *Mongo) ViewOneC(collName string, query, result interface{}, selectQ ...interface{}) (string, error) {
	return m.viewOne(m.info.Database, collName, query, result, selectQ...)
}

// type Iquery interface {
// 	//New(string, interface{}) bson.M
// 	In(string, interface{}) bson.M
// }

type Query bson.M

func (m *Mongo) NewQuery(field string, value interface{}) bson.M {
	return bson.M{field: value}
}

func (m *Mongo) In(field string, value interface{}) bson.M {
	return bson.M{field: bson.M{"$in": value}}
}
func (m *Mongo) Select(qs []string) bson.M {
	q := bson.M{}
	for _, m := range qs {
		q[m] = 1
	}
	return q
}
