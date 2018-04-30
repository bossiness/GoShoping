package mongodb

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	proto "btdxcx.com/micro/product-srv/proto/product"
)

const (
	optionsCollectionName = "options"
)

// Option DB
type Option struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Code      string        `bson:"code,omitempty"`
	Name      string        `bson:"name,omitempty"`
	Options   []OptionValue `bson:"options,omitempty"`
	UpdatedAt int64         `bson:"updated_at,omitempty"`
	CreatedAt int64         `bson:"created_at,omitempty"`
}

// OptionValue inc
type OptionValue struct {
	Value       string `bson:"value,omitempty"`
	Description string `bson:"description,omitempty"`
}

// CreateOption Insert
func (m *Mongo) CreateOption(dbname string, record *proto.OptionRecord) error {
	c := m.session.DB(dbname).C(optionsCollectionName)

	index := mgo.Index{
		Key:        []string{"code"},
		Unique:     true,
		DropDups:   true,
		Background: true,
	}

	if err := c.EnsureIndex(index); err != nil {
		return err
	}

	incs := []OptionValue{}
	for _, op := range record.Options {
		inc := OptionValue{
			Value:       op.Value,
			Description: op.Description,
		}
		incs = append(incs, inc)
	}

	doc := Option{
		ID:        bson.NewObjectId(),
		Code:      record.Code,
		Name:      record.Name,
		Options:   incs,
		UpdatedAt: time.Now().Unix(),
		CreatedAt: time.Now().Unix(),
	}

	return c.Insert(&doc)

}

// ReadOptions Find
func (m *Mongo) ReadOptions(dbname string, offset int, limit int) (*[]*proto.OptionRecord, error) {
	c := m.session.DB(dbname).C(optionsCollectionName)

	results := []Option{}
	if err := c.Find(nil).Skip(offset).Limit(limit).All(&results); err != nil {
		return nil, err
	}

	records := []*proto.OptionRecord{}
	for _, item := range results {
		options := []*proto.OptionRecord_OptionValue{}
		for _, op := range item.Options {
			options = append(options,
				&proto.OptionRecord_OptionValue{
					Value:       op.Value,
					Description: op.Description,
				})
		}
		record := &proto.OptionRecord{
			Code:      item.Code,
			Name:      item.Name,
			Options:   options,
			UpdatedAt: item.UpdatedAt,
			CreatedAt: item.CreatedAt,
		}
		records = append(records, record)
	}
	return &records, nil
}

// ReadOption Find
func (m *Mongo) ReadOption(dbname string, code string) (*proto.OptionRecord, error) {
	c := m.session.DB(dbname).C(optionsCollectionName)

	result := Option{}
	selector := bson.M{"code": code}
	if err := c.Find(selector).One(&result); err != nil {
		return nil, err
	}

	options := []*proto.OptionRecord_OptionValue{}
	for _, op := range result.Options {
		options = append(options,
			&proto.OptionRecord_OptionValue{
				Value:       op.Value,
				Description: op.Description,
			})
	}
	record := &proto.OptionRecord{
		Code:      result.Code,
		Name:      result.Name,
		Options:   options,
		UpdatedAt: result.UpdatedAt,
		CreatedAt: result.CreatedAt,
	}

	return record, nil
}

// UpdateOption Update
func (m *Mongo) UpdateOption(dbname string, code string, record *proto.OptionRecord) error {
	c := m.session.DB(dbname).C(optionsCollectionName)

	incs := []OptionValue{}
	for _, op := range record.Options {
		inc := OptionValue{
			Value:       op.Value,
			Description: op.Description,
		}
		incs = append(incs, inc)
	}

	selector := bson.M{"code": code}
	updataData := bson.M{"$set": bson.M{
		"options":    incs,
		"updated_at": time.Now()}}

	return c.Update(selector, updataData)
}

// DeleteOption Remove
func (m *Mongo) DeleteOption(dbname string, code string) error {
	c := m.session.DB(dbname).C(optionsCollectionName)

	selector := bson.M{"code": code}
	return c.Remove(selector)
}
