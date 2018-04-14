package mongodb

import (
	"time"

	"btdxcx.com/micro/product-srv/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	proto "btdxcx.com/micro/product-srv/proto/product"
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
	attributesCollectionName = "attributes"
	optionsCollectionName    = "options"
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

// Attribute DB
type Attribute struct {
	ID            bson.ObjectId `bson:"_id,omitempty"`
	Code          string        `bson:"code"`
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
		"type":          record.Type,
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

// Option DB
type Option struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Code      string        `bson:"code,omitempty"`
	Options   []OptionValue `bson:"options,omitempty"`
	UpdatedAt int64         `json:"updated_at,omitempty"`
	CreatedAt int64         `json:"created_at,omitempty"`
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
