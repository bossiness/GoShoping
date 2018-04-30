package mongodb

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	proto "btdxcx.com/micro/product-srv/proto/product"
)

const (
	attributesCollectionName = "attributes"
)

// Attribute DB
type Attribute struct {
	ID            bson.ObjectId `bson:"_id,omitempty"`
	Code          string        `bson:"code"`
	Name          string        `bson:"name"`
	Type          string        `bson:"type"`
	Configuration []string      `bson:"configuration,omitempty"`
	UpdatedAt     int64         `bson:"updated_at,omitempty"`
	CreatedAt     int64         `bson:"created_at,omitempty"`
}

// CreateAttribute Insert
func (m *Mongo) CreateAttribute(dbname string, record *proto.AttributesRecord) error {
	c := m.session.DB(dbname).C(attributesCollectionName)

	index := mgo.Index{
		Key:        []string{"code"},
		Unique:     true,
		DropDups:   true,
		Background: true,
	}

	if err := c.EnsureIndex(index); err != nil {
		return err
	}

	doc := &Attribute{
		ID:            bson.NewObjectId(),
		Code:          record.Code,
		Name:          record.Name,
		Type:          record.Type,
		Configuration: record.Configuration,
		UpdatedAt:     time.Now().Unix(),
		CreatedAt:     time.Now().Unix(),
	}

	return c.Insert(doc)
}

// ReadAttributes Find
func (m *Mongo) ReadAttributes(dbname string, offset int, limit int) (*[]*proto.AttributesRecord, error) {
	c := m.session.DB(dbname).C(attributesCollectionName)

	results := []Attribute{}
	if err := c.Find(nil).Skip(offset).Limit(limit).All(&results); err != nil {
		return nil, err
	}

	records := []*proto.AttributesRecord{}
	for _, item := range results {
		record := &proto.AttributesRecord{
			Code:          item.Code,
			Type:          item.Type,
			Name:          item.Name,
			Configuration: item.Configuration,
			UpdatedAt:     item.UpdatedAt,
			CreatedAt:     item.CreatedAt,
		}
		records = append(records, record)
	}
	return &records, nil
}

// ReadAttribute Find
func (m *Mongo) ReadAttribute(dbname string, code string) (*proto.AttributesRecord, error) {
	c := m.session.DB(dbname).C(attributesCollectionName)

	result := Attribute{}
	if err := c.Find(bson.M{"code": code}).One(&result); err != nil {
		return nil, err
	}

	record := proto.AttributesRecord{
		Code:          result.Code,
		Type:          result.Type,
		Name:          result.Name,
		Configuration: result.Configuration,
		UpdatedAt:     result.UpdatedAt,
		CreatedAt:     result.CreatedAt,
	}

	return &record, nil
}

// UpdateAttribute update
func (m *Mongo) UpdateAttribute(dbname string, code string, record *proto.AttributesRecord) error {
	c := m.session.DB(dbname).C(attributesCollectionName)

	selector := bson.M{"code": code}
	updataData := bson.M{"$set": bson.M{
		// "type":          record.Type,
		"configuration": record.Configuration,
		"updated_at":    time.Now()}}

	return c.Update(selector, updataData)
}

// DeleteAttribute Remove
func (m *Mongo) DeleteAttribute(dbname string, code string) error {
	c := m.session.DB(dbname).C(attributesCollectionName)

	selector := bson.M{"code": code}
	return c.Remove(selector)
}
