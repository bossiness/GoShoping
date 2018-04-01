package mongodb

import (
	"gopkg.in/mgo.v2/bson"
	"btdxcx.com/micro/shop-srv/db"
	"gopkg.in/mgo.v2"

	proto "btdxcx.com/micro/shop-srv/proto/shop"
)

// Mongo DB
type Mongo struct {
	session *mgo.Session
}

var (
	// DBUrl mongodb URL
	DBUrl = "localhost:27017"
)

const (
	databaseName  = "center"
	keyCollection = "shop-key"
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

// ShopKey DB
type ShopKey struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	UUID    string        `bson:"uuid"`
	BackKey string        `bson:"back_key"`
	MiniKey string        `bson:"mini_key"`
}

// ReadKey form uuid
func (m *Mongo) ReadKey(uuid string) (*proto.ShopKeyID, error) {
	c := m.session.DB(databaseName).C(keyCollection)

	result := &ShopKey{}
	if err := c.Find(bson.M{ "uuid": uuid}).One(result); err != nil {
		return nil, err
	}

	shopKey := &proto.ShopKeyID{ 
		BackKey: result.BackKey, 
		MiniKey: result.MiniKey }

	return shopKey, nil
}

// CreateKey form uuid
func (m *Mongo) CreateKey(uuid string, proto *proto.ShopKeyID) error {
	c := m.session.DB(databaseName).C(keyCollection)

	shopKey := &ShopKey {
		ID: bson.NewObjectId(),
		UUID: uuid,
		BackKey: proto.BackKey,
		MiniKey: proto.MiniKey,
	}
	return c.Insert(shopKey)
}

// DeleteKey form uuid
func (m *Mongo) DeleteKey(uuid string) error {
	c := m.session.DB(databaseName).C(keyCollection)
	return c.Remove(bson.M{ "uuid": uuid })
}
