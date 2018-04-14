package db

import (
	proto "btdxcx.com/micro/product-srv/proto/product"
)

// DB is 数据库接口
type DB interface {
	Init() error
	Deinit()
	Attribute
	Option
}

// Attribute is Attribute数据接口
type Attribute interface {
	CreateAttribute(string, *proto.AttributesRecord) error
	ReadAttributes(string, int, int) (*[]*proto.AttributesRecord, error)
	ReadAttribute(string, string) (*proto.AttributesRecord, error)
	UpdateAttribute(string, string, *proto.AttributesRecord) error
	DeleteAttribute(string, string) error
}

// Option is Option
type Option interface {
	CreateOption(string, *proto.OptionRecord) error
	ReadOptions(string, int, int) (*[]*proto.OptionRecord, error)
	ReadOption(string, string) (*proto.OptionRecord, error)
	UpdateOption(string, string, *proto.OptionRecord) error
	DeleteOption(string, string) error
}


var (
	db DB
)

// Register db Imp
func Register(backend DB) {
	db = backend
}

// Init 数据库初始化
func Init() error {
	return db.Init()
}

// Deinit 析构
func Deinit() {
	db.Deinit()
}

// CreateAttribute create
func CreateAttribute(dbname string, record *proto.AttributesRecord) error {
	return db.CreateAttribute(dbname, record)
}

// ReadAttributes read list
func ReadAttributes(dbname string, offset int, limit int) (*[]*proto.AttributesRecord, error) {
	return db.ReadAttributes(dbname, offset, limit)
}

// ReadAttribute read
func ReadAttribute(dbname string, code string) (*proto.AttributesRecord, error) {
	return db.ReadAttribute(dbname, code)
}

// UpdateAttribute update
func UpdateAttribute(dbname string, code string, record *proto.AttributesRecord) error {
	return db.UpdateAttribute(dbname, code, record)
}

// DeleteAttribute delete
func DeleteAttribute(dbname string, code string) error {
	return db.DeleteAttribute(dbname, code)
}

//////////////////////////////

// CreateOption create
func CreateOption(dbname string, record *proto.OptionRecord) error {
	return db.CreateOption(dbname, record)
}

// ReadOptions read
func ReadOptions(dbname string, offset int, limit int) (*[]*proto.OptionRecord, error) {
	return db.ReadOptions(dbname, offset, limit)
}

// ReadOption read
func ReadOption(dbname string, code string) (*proto.OptionRecord, error) {
	return db.ReadOption(dbname, code)
}

// UpdateOption update
func UpdateOption(dbname string, code string, record *proto.OptionRecord) error {
	return db.UpdateOption(dbname, code, record)
}

// DeleteOption delete
func DeleteOption(dbname string, code string) error {
	return db.DeleteOption(dbname, code)
}
