package mongodb

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	proto "btdxcx.com/micro/member-srv/proto/member"
)

const (
	customersCollectionName = "customers"
)

// Customer DB
type Customer struct {
	ID                  bson.ObjectId `bson:"_id,omitempty"`
	Username            string        `bson:"username,omitempty"`
	FirstName           string        `bson:"firstName,omitempty"`
	LastName            string        `bson:"lastName,omitempty"`
	Phone               string        `bson:"phone,omitempty"`
	Email               string        `bson:"email,omitempty"`
	Portrait            string        `bson:"portrait,omitempty"`
	Gender              string        `bson:"gender,omitempty"`
	Birthday            int64         `bson:"birthday,omitempty"`
	Groups              []string      `bson:"groups,omitempty"`
	RegisteredAt        int64         `bson:"registered_at,omitempty"`
	AccessAt            int64         `bson:"access_at,omitempty"`
	Integral            int64         `bson:"integral,omitempty"`
	BuyNumber           int64         `bson:"buy_number,omitempty"`
	TotalPurchaseAmount int64         `bson:"totalPurchaseAmount,omitempty"`
	Role                []string      `bson:"role,omitempty"`
	Superior            string        `bson:"superior,omitempty"`
}

// CreateCustomer create Customer
func (m *Mongo) CreateCustomer(dbname string, record *proto.CustomerRecord) error {
	c := m.session.DB(dbname).C(customersCollectionName)

	index := mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		DropDups:   true,
		Background: true,
	}

	if err := c.EnsureIndex(index); err != nil {
		return err
	}

	doc := &Customer{
		ID:                  bson.NewObjectId(),
		Username:            record.Username,
		FirstName:           record.FirstName,
		LastName:            record.LastName,
		Phone:               record.Phone,
		Email:               record.Email,
		Portrait:            record.Portrait,
		Gender:              record.Gender,
		Groups:              record.Groups,
		RegisteredAt:        time.Now().Unix(),
		AccessAt:            time.Now().Unix(),
		Integral:            record.Integral,
		Birthday:            record.Birthday,
		TotalPurchaseAmount: record.TotalPurchaseAmount,
		Role:                record.Role,
		Superior:            record.Superior,
	}

	return c.Insert(doc)
}

// ReadCustomers read Customers
func (m *Mongo) ReadCustomers(dbname string, offset int, limit int) (*[]*proto.CustomerRecord, error) {
	c := m.session.DB(dbname).C(customersCollectionName)

	results := []Customer{}
	if err := c.Find(nil).Skip(offset).Limit(limit).All(&results); err != nil {
		return nil, err
	}

	records := []*proto.CustomerRecord{}
	for _, it := range results {
		record := &proto.CustomerRecord{
			Id:                  it.ID.Hex(),
			Username:            it.Username,
			FirstName:           it.FirstName,
			LastName:            it.LastName,
			Phone:               it.Phone,
			Email:               it.Email,
			Portrait:            it.Portrait,
			Gender:              it.Gender,
			Groups:              it.Groups,
			RegisteredAt:        it.RegisteredAt,
			AccessAt:            it.AccessAt,
			Integral:            it.Integral,
			Birthday:            it.Birthday,
			TotalPurchaseAmount: it.TotalPurchaseAmount,
			Role:                it.Role,
			Superior:            it.Superior,
		}
		records = append(records, record)
	}
	return &records, nil
}

// ReadCustomer read a Customer
func (m *Mongo) ReadCustomer(dbname string, id string) (*proto.CustomerRecord, error) {
	c := m.session.DB(dbname).C(customersCollectionName)

	result := Customer{}
	if err := c.FindId(bson.ObjectIdHex(id)).One(&result); err != nil {
		return nil, err
	}

	record := proto.CustomerRecord{
		Id:                  result.ID.Hex(),
		Username:            result.Username,
		FirstName:           result.FirstName,
		LastName:            result.LastName,
		Phone:               result.Phone,
		Email:               result.Email,
		Portrait:            result.Portrait,
		Gender:              result.Gender,
		Groups:              result.Groups,
		RegisteredAt:        result.RegisteredAt,
		AccessAt:            result.AccessAt,
		Integral:            result.Integral,
		Birthday:            result.Birthday,
		TotalPurchaseAmount: result.TotalPurchaseAmount,
		Role:                result.Role,
		Superior:            result.Superior,
	}

	return &record, nil
}

// ReadCustomerFromName update a Customer
func (m *Mongo) ReadCustomerFromName(dbname string, username string) (*proto.CustomerRecord, error) {
	c := m.session.DB(dbname).C(customersCollectionName)

	result := Customer{}
	if err := c.Find(bson.M{"username": username}).One(&result); err != nil {
		return nil, err
	}

	record := proto.CustomerRecord{
		Id:                  result.ID.Hex(),
		Username:            result.Username,
		FirstName:           result.FirstName,
		LastName:            result.LastName,
		Phone:               result.Phone,
		Email:               result.Email,
		Portrait:            result.Portrait,
		Gender:              result.Gender,
		Groups:              result.Groups,
		RegisteredAt:        result.RegisteredAt,
		AccessAt:            result.AccessAt,
		Integral:            result.Integral,
		Birthday:            result.Birthday,
		TotalPurchaseAmount: result.TotalPurchaseAmount,
		Role:                result.Role,
		Superior:            result.Superior,
	}

	return &record, nil
}

// UpdateCustomer update a Customer
func (m *Mongo) UpdateCustomer(dbname string, id string, record *proto.CustomerRecord) error {
	c := m.session.DB(dbname).C(customersCollectionName)

	updataData := bson.M{"$set": bson.M{
		"firstName": record.FirstName,
		"lastName":  record.LastName,
		"phone":     record.Phone,
		"email":     record.Email,
		"portrait":  record.Portrait,
		"gender":    record.Gender,
		"birthday":  record.Birthday,
		"groups":    record.Groups,
	}}

	return c.UpdateId(bson.ObjectIdHex(id), updataData)
}

// DeleteCustomer delete a Customer
func (m *Mongo) DeleteCustomer(dbname string, id string) error {
	c := m.session.DB(dbname).C(customersCollectionName)
	return c.RemoveId(bson.ObjectIdHex(id))
}
