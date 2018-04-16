package mongodb

import (
	"btdxcx.com/micro/product-srv/db"
	"gopkg.in/mgo.v2"
)

// Mongo DB
type Mongo struct {
	session *mgo.Session
}

var (
	// DBUrl mongodb URL
	DBUrl = "localhost:27017"
)

func init() {
	db.Register(new(Mongo))
}

// Init 数据库初始化
func (m *Mongo) Init() error {
	session, err := mgo.Dial(DBUrl)
	if err != nil {
		return err
	}
	session.SetMode(mgo.Monotonic, true)
	m.session = session
	return nil
}

// Deinit 资源释放
func (m *Mongo) Deinit() {
	if m.session != nil {
		m.session.Close()
	}
}
