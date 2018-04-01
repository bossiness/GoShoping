package mongodb

import (
	"btdxcx.com/micro/shop-srv/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

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
	ID     bson.ObjectId     `bson:"_id,omitempty"`
	UUID   string            `bson:"uuid"`
	TagKey map[string]string `bson:"tag_key"`
}

// ReadKey form uuid
func (m *Mongo) ReadKey(uuid string) (*proto.ShopTagKeys, error) {
	c := m.session.DB(databaseName).C(keyCollection)

	result := &ShopKey{}
	if err := c.Find(bson.M{"uuid": uuid}).One(result); err != nil {
		return nil, err
	}

	shopKey := &proto.ShopTagKeys{
		Tagkeys: result.TagKey,
	}

	return shopKey, nil
}

// CreateKey form uuid
func (m *Mongo) CreateKey(uuid string, proto *proto.ShopTagKeys) error {
	c := m.session.DB(databaseName).C(keyCollection)

	index := mgo.Index {
		Key: []string{"uuid"},
		Unique: true,
		DropDups: true,
		Background: true,
	}

	if err := c.EnsureIndex(index); err != nil {
		return err
	}

	shopKey := &ShopKey{
		ID:     bson.NewObjectId(),
		UUID:   uuid,
		TagKey: proto.Tagkeys,
	}
	return c.Insert(shopKey)
}

// DeleteKey form uuid
func (m *Mongo) DeleteKey(uuid string) error {
	c := m.session.DB(databaseName).C(keyCollection)
	return c.Remove(bson.M{"uuid": uuid})
}
