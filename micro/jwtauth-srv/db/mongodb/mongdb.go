package mongodb

import (
	"github.com/micro/go-micro/errors"

	"btdxcx.com/micro/jwtauth-srv/db"
	proto "btdxcx.com/micro/jwtauth-srv/proto/auth"
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
	collectionName = "jwtauth"
)

// Jwtauth DB
type Jwtauth struct {
	id       bson.ObjectId `bson:"_id,omitempty"`
	ClientID string        `bson:"client_id"`
	JTI      string        `bson:"jti"`
	Access   string        `bson:"access_token"`
	Refresh  string        `bson:"refresh_token"`
	Cipher   string        `bson:"cipher"`
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

// Read Auth
func (m *Mongo) Read(dbName string, accessToken string) (*proto.Record, error) {

	c := m.session.DB(dbName).C(collectionName)
	result, err := m.read(c, accessToken)
	if err != nil {
		return nil, err
	}
	record := &proto.Record{
		Jti:          result.JTI,
		ClientId:     result.ClientID,
		AccessToken:  result.Access,
		RefreshToken: result.Refresh,
		Cipher:       result.Cipher,
	}

	return record, nil
}

func (m *Mongo) read(c *mgo.Collection, accessToken string) (*Jwtauth, error) {
	result := &Jwtauth{}
	if err := c.Find(bson.M{"access_token": accessToken}).One(result); err != nil {
		return nil, err
	}
	return result, nil
}

// Create Auth
func (m *Mongo) Create(dbName string, record *proto.Record) error {

	c := m.session.DB(dbName).C(collectionName)

	_, err := m.read(c, record.ClientId)
	if err == nil {
		return errors.Conflict("jwtauth.srv", "account conflict.")
	}

	result := &Jwtauth{
		id:       bson.NewObjectId(),
		JTI:      record.Jti,
		ClientID: record.ClientId,
		Access:   record.AccessToken,
		Refresh:  record.RefreshToken,
		Cipher:   record.Cipher,
	}

	return c.Insert(result)
}

// Update Auth
func (m *Mongo) Update(dbName string, record *proto.Record) error {

	c := m.session.DB(dbName).C(collectionName)
	selector := bson.M{"jti": record.Jti}

	updataData := bson.M{"$set": bson.M{
		"access_token":  record.AccessToken,
	}}

	return c.Update(selector, updataData)
}

// DeleteAccessToken Auth
func (m *Mongo) DeleteAccessToken(dbName string, token string) error {
	c := m.session.DB(dbName).C(collectionName)
	return c.Remove(bson.M{"access_token": token})
}

// DeleteRefreshToken Auth
func (m *Mongo) DeleteRefreshToken(dbName string, token string) error {
	c := m.session.DB(dbName).C(collectionName)
	return c.Remove(bson.M{"refresh_token": token})
}

// ReadFormRefreshToken read token
func (m *Mongo) ReadFormRefreshToken(dbName string, refreshToken string) (*proto.Record, error) {
	c := m.session.DB(dbName).C(collectionName)
	result, err := m.readFrom(c, refreshToken)
	if err != nil {
		return nil, err
	}
	record := &proto.Record{
		Jti:          result.JTI,
		ClientId:     result.ClientID,
		AccessToken:  result.Access,
		RefreshToken: result.Refresh,
		Cipher:       result.Cipher,
	}

	return record, nil
}

func (m *Mongo) readFrom(c *mgo.Collection, refreshToken string) (*Jwtauth, error) {
	result := &Jwtauth{}
	if err := c.Find(bson.M{"refresh_token": refreshToken}).One(result); err != nil {
		return nil, err
	}
	return result, nil
}