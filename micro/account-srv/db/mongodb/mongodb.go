package mongodb

import (
	"github.com/micro/go-micro/errors"
	"time"

	"btdxcx.com/micro/account-srv/db"
	proto "btdxcx.com/micro/account-srv/proto/account"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
	collectionName = "account"
)

// Account DB
type Account struct {
	id           bson.ObjectId     `bson:"_id,omitempty"`
	Type         string            `bson:"type"`
	ClientID     string            `bson:"client_id"`
	ClientSecret string            `bson:"client_secret"`
	Metadata     map[string]string `bson:"metadata"`
	CreatedAt    time.Time         `bson:"created_at"`
	UpdatedAt    time.Time         `bson:"updated_at"`
}

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

// Read Account
func (m *Mongo) Read(dbName string, clientID string) (*proto.Record, error) {

	c := m.session.DB(dbName).C(collectionName)
	result, err := m.read(c, clientID)
	if err != nil {
		return nil, err
	}
	record := &proto.Record{
		Id:           result.id.Hex(),
		Type:         result.Type,
		ClientId:     result.ClientID,
		ClientSecret: result.ClientSecret,
		Metadata:     result.Metadata,
		Created:      result.CreatedAt.Unix(),
		Updated:      result.UpdatedAt.Unix(),
	}

	return record, nil
}

func (m *Mongo) read(c *mgo.Collection, clientID string) (*Account, error) {
	result := &Account{}
	if err := c.Find(bson.M{"client_id": clientID}).One(result); err != nil {
		return nil, err
	}
	return result, nil
}

// Create Account
func (m *Mongo) Create(dbName string, record *proto.Record) error {

	c := m.session.DB(dbName).C(collectionName)

	_, err := m.read(c, record.ClientId)
	if err == nil {
		return errors.Conflict("account.srv", "account conflict.")
	}

	index := mgo.Index{
		Key:        []string{"client_id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
	}

	if err := c.EnsureIndex(index); err != nil {
		return err
	}

	result := &Account{
		id:           bson.NewObjectId(),
		Type:         record.Type,
		ClientID:     record.ClientId,
		ClientSecret: record.ClientSecret,
		Metadata:     record.Metadata,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return c.Insert(result)
}

// Update Account
func (m *Mongo) Update(dbName string, record *proto.Record) error {

	c := m.session.DB(dbName).C(collectionName)
	selector := bson.M{"_id": bson.ObjectIdHex(record.Id)}

	updataData := bson.M{"$set": bson.M{
		"type":          record.Type,
		"client_secret": record.ClientSecret,
		"metadata":      record.Metadata,
		"updated_at":    time.Now()}}

	return c.Update(selector, updataData)
}

// Delete Account
func (m *Mongo) Delete(dbName string, id string) error {
	c := m.session.DB(dbName).C(collectionName)
	return c.RemoveId(bson.ObjectIdHex(id))
}

// Search Account
func (m *Mongo) Search(dbName string, request *proto.SearchRequest) (*[]*proto.Record, error) {

	c := m.session.DB(dbName).C(collectionName)

	if request.ClientId != "" {
		result := &[]Account{}
		if err := c.Find(bson.M{"client_id": request.ClientId}).All(result); err != nil {
			return nil, err
		}
		return mapper(result), nil
	}

	result := &[]Account{}
	if err := c.Find(bson.M{"type": request.Type}).Skip(int(request.Offset)).Limit(int(request.Limit)).All(result); err != nil {
		return nil, err
	}

	return mapper(result), nil
}

func mapper(accounts *[]Account) *[]*proto.Record {
	record := []*proto.Record{}
	for _, a := range *accounts {
		r := proto.Record{Id: a.id.Hex(),
			Type:         a.Type,
			ClientId:     a.ClientID,
			ClientSecret: a.ClientSecret,
			Metadata:     a.Metadata,
			Created:      a.CreatedAt.Unix(),
			Updated:      a.UpdatedAt.Unix(),
		}
		record = append(record, &r)
	}
	return &record
}
