package db

import (
	"errors"
	proto "btdxcx.com/micro/taxons-srv/proto/imp"
)

// DB is 数据库接口
type DB interface {
	Init() error
	Taxons
}

// Taxons is Taxons数据接口
type Taxons interface {
	Read(dbname string) (* proto.TaxonsMessage, error)
	Delete(dbname string, id string) error
	Create(dbname string, data *proto.TaxonsMessage) (string, error)
	Update(dbname string, data *proto.TaxonsMessage) error
}

var (
	db DB

	// ErrNotFound is 找不到数据错误
	ErrNotFound = errors.New("not found")
)

// Register db Imp
func Register(backend DB) {
	db = backend
}

// Init 数据库初始化
func Init() error {
	return db.Init()
}

// Read 读取数据
func Read(dbname string) (* proto.TaxonsMessage, error) {
	return db.Read(dbname)
}

// Delete 删除数据
func Delete(dbname string, id string) error {
	return db.Delete(dbname, id)
}

// Create 创造数据
func Create(dbname string, data *proto.TaxonsMessage) (string, error) {
	return db.Create(dbname, data)
}

// Update 更新数据
func Update(dbname string, data *proto.TaxonsMessage) error {
	return db.Update(dbname, data)
}
