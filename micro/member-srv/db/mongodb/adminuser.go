package mongodb

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	proto "btdxcx.com/micro/member-srv/proto/member"
)

const (
	adminusersCollectionName = "adminusers"
)

// AdminUser DB
type AdminUser struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	Username     string        `bson:"username,omitempty"`
	FirstName    string        `bson:"firstName,omitempty"`
	LastName     string        `bson:"lastName,omitempty"`
	Phone        string        `bson:"phone,omitempty"`
	Email        string        `bson:"email,omitempty"`
	Portrait     string        `bson:"portrait,omitempty"`
	RegisteredAt int64         `bson:"registered_at,omitempty"`
	AccessAt     int64         `bson:"access_at,omitempty"`
	Enable       bool          `bson:"enable,omitempty"`
	Role         []string      `bson:"role,omitempty"`
}

// CreateAdminUser create AdminUser
func (m *Mongo) CreateAdminUser(dbname string, record *proto.AdminUserRecord) error {
	c := m.session.DB(dbname).C(adminusersCollectionName)

	index := mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		DropDups:   true,
		Background: true,
	}

	if err := c.EnsureIndex(index); err != nil {
		return err
	}

	doc := &AdminUser{
		ID:                  bson.NewObjectId(),
		Username:            record.Username,
		FirstName:           record.FirstName,
		LastName:            record.LastName,
		Phone:               record.Phone,
		Email:               record.Email,
		Portrait:            record.Portrait,
		RegisteredAt:        time.Now().Unix(),
		AccessAt:            time.Now().Unix(),
		Role:                record.Role,
	}

	return c.Insert(doc)
}

// ReadAdminUsers read AdminUsers
func (m *Mongo) ReadAdminUsers(dbname string, offset int, limit int) (*[]*proto.AdminUserRecord, error) {
	c := m.session.DB(dbname).C(adminusersCollectionName)

	results := []AdminUser{}
	if err := c.Find(nil).Skip(offset).Limit(limit).All(&results); err != nil {
		return nil, err
	}

	records := []*proto.AdminUserRecord{}
	for _, it := range results {
		record := &proto.AdminUserRecord{
			Id:                  it.ID.Hex(),
			Username:            it.Username,
			FirstName:           it.FirstName,
			LastName:            it.LastName,
			Phone:               it.Phone,
			Email:               it.Email,
			Portrait:            it.Portrait,
			RegisteredAt:        it.RegisteredAt,
			AccessAt:            it.AccessAt,
			Role:                it.Role,
		}
		records = append(records, record)
	}
	return &records, nil
}

// ReadAdminUser read a AdminUser
func (m *Mongo) ReadAdminUser(dbname string, id string) (*proto.AdminUserRecord, error) {
	c := m.session.DB(dbname).C(adminusersCollectionName)

	result := AdminUser{}
	if err := c.FindId(bson.ObjectIdHex(id)).One(&result); err != nil {
		return nil, err
	}

	record := proto.AdminUserRecord{
		Id:                  result.ID.Hex(),
		Username:            result.Username,
		FirstName:           result.FirstName,
		LastName:            result.LastName,
		Phone:               result.Phone,
		Email:               result.Email,
		Portrait:            result.Portrait,
		RegisteredAt:        result.RegisteredAt,
		AccessAt:            result.AccessAt,
		Role:                result.Role,
	}

	return &record, nil
}

// ReadAdminUserFromName update a AdminUser
func (m *Mongo) ReadAdminUserFromName(dbname string, username string) (*proto.AdminUserRecord, error) {
	c := m.session.DB(dbname).C(adminusersCollectionName)

	result := AdminUser{}
	if err := c.Find(bson.M{"username": username}).One(&result); err != nil {
		return nil, err
	}

	record := proto.AdminUserRecord{
		Id:                  result.ID.Hex(),
		Username:            result.Username,
		FirstName:           result.FirstName,
		LastName:            result.LastName,
		Phone:               result.Phone,
		Email:               result.Email,
		Portrait:            result.Portrait,
		RegisteredAt:        result.RegisteredAt,
		AccessAt:            result.AccessAt,
		Role:                result.Role,
	}

	return &record, nil
}

// UpdateAdminUser update a AdminUser
func (m *Mongo) UpdateAdminUser(dbname string, id string, record *proto.AdminUserRecord) error {
	c := m.session.DB(dbname).C(adminusersCollectionName)

	updataData := bson.M{"$set": bson.M{
		"firstName": record.FirstName,
		"lastName":  record.LastName,
		"phone":     record.Phone,
		"email":     record.Email,
		"portrait":  record.Portrait,
	}}

	return c.UpdateId(bson.ObjectIdHex(id), updataData)
}

// DeleteAdminUser delete a AdminUser
func (m *Mongo) DeleteAdminUser(dbname string, id string) error {
	c := m.session.DB(dbname).C(adminusersCollectionName)
	return c.RemoveId(bson.ObjectIdHex(id))
}
