package db

import (
	proto "btdxcx.com/micro/account-srv/proto/account"
)

// DB is 数据库接口
type DB interface {
	Init() error
	Deinit()
	Account
}

// Account is Account数据接口
type Account interface {
	Read(string, string) (*proto.Record, error)
	Create(string, *proto.Record) error
	Update(string, *proto.Record) error
	Delete(string, string) error
	Search(string, *proto.SearchRequest) (*[]*proto.Record, error)
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
func Deinit()  {
	db.Deinit()
}

// Read Account
func Read(dbName string, clientID string) (*proto.Record, error) {
	return db.Read(dbName, clientID)
}

// Create Account
func Create(dbName string, record *proto.Record) error {
	return db.Create(dbName, record)
}

// Update Account
func Update(dbName string, record *proto.Record) error {
	return db.Update(dbName, record)
}

// Delete Account
func Delete(dbName string, id string) error {
	return db.Delete(dbName, id)
}

// Search Account
func Search(dbName string, request *proto.SearchRequest) (*[]*proto.Record, error) {
	return db.Search(dbName, request)
}
